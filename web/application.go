package web

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaoxue1272/club-5fw-backend/config"
	"strconv"
)

var engine *gin.Engine

var runAddr string

func Init(webConfig *config.WebConfiguration) {
	runAddr = webConfig.Host + ":" + strconv.Itoa(webConfig.Port)
	engine = gin.New()
	initJwt(webConfig.Jwt)
	initJsonCipher(webConfig.CipherJson)
	engine.Use(recovery, corsAccess(), requestLogger)
	initRouters(engine)
	engine.Use(errorHandler)
}

func StartWebServer() {
	err := engine.Run(runAddr)
	if err != nil {
		panic(err)
	}
}
