package web

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaoxue1272/club-5fw-backend/config"
	"strconv"
)

var engine *gin.Engine

func init() {
	engine = gin.New()
	engine.Use(requestLogger(), recovery())
	initRouters(engine)
}

func StartWebServer(webConfig *config.WebConfiguration) {
	err := engine.Run(webConfig.Host + ":" + strconv.Itoa(webConfig.Port))
	if err != nil {
		panic(err)
	}
}
