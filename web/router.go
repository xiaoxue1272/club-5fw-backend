package web

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaoxue1272/club-5fw-backend/config"
	"github.com/xiaoxue1272/club-5fw-backend/model"
	"github.com/xiaoxue1272/club-5fw-backend/service"
	"net/http"
)

func initRouters(engine *gin.Engine) {
	basicEndPoint(&engine.RouterGroup)
	authGroupEndpoint(engine.Group("/auth", jwtAuthorization))
}

func basicEndPoint(group *gin.RouterGroup) {
	group.GET("/publicKey", func(context *gin.Context) {
		context.String(http.StatusOK, config.GetConfiguration().Web.CipherJson.Rsa.PublicKey)
	})

	group.POST("/sign", func(context *gin.Context) {
		userSign := &model.UserSign{}
		err := context.ShouldBindBodyWith(userSign, CipherJson)
		if err != nil {
			_ = context.Error(err)
			return
		}
		jwt, err := service.Sign(userSign)
		if err != nil {
			_ = context.Error(err)
			return
		}
		authToken, err := generateJwt(jwt)
		if err != nil {
			_ = context.Error(err)
			return
		}
		context.JSON(http.StatusOK, &gin.H{"token": authToken})
	})
}

func authGroupEndpoint(group *gin.RouterGroup) {
	group.GET("/user/icon", func(context *gin.Context) {
		// todo 获取用户头像
	})
	group.GET("/hello", func(context *gin.Context) {
		value, _ := context.Get("jwt")
		context.PureJSON(http.StatusOK, value)
	})
}
