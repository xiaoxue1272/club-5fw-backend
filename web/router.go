package web

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
)

func initRouters(engine *gin.Engine) {
	initJwt()
	basicEndPoint(&engine.RouterGroup)
	authGroup := engine.Group("/auth", JwtAuthorization())
	authGroupEndpoint(authGroup)
}

func basicEndPoint(group *gin.RouterGroup) {
	group.GET("/login", func(context *gin.Context) {
		authToken, err := generateJwt(&gin.H{"user": "test"})
		if err != nil {
			logger.Warnf("generate Jwt Token Failed.\nError: %v\n", err)
			context.String(http.StatusInternalServerError, err.Error())
			context.Abort()
			return
		}
		context.JSON(http.StatusOK, &gin.H{"token": authToken})
	})
}

func authGroupEndpoint(group *gin.RouterGroup) {
	group.GET("/hello", func(context *gin.Context) {
		value, _ := context.Get("data")
		context.PureJSON(http.StatusOK, value)
	})
}
