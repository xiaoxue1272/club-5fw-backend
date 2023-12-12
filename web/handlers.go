package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func requestLogger(context *gin.Context) {
	request := context.Request
	if logger.IsLevelEnabled(logger.DebugLevel) {
		logger.Debugf(
			"Request income\nMethod: %s\nPath: %s\nContent Length: %d\nHeaders: %s\nCookies: %s\nRemoteAddr: %s\n",
			request.Method,
			request.RequestURI,
			request.ContentLength,
			request.Header,
			request.Cookies(),
			request.RemoteAddr,
		)
	} else if logger.IsLevelEnabled(logger.InfoLevel) {
		logger.Info(
			"Request income\nMethod: %s\nPath: %s\nClientIP: %s\n",
			request.Method,
			request.RequestURI,
			request.ContentLength,
			request.Header,
			request.Cookies(),
			context.ClientIP(),
		)
	}
}

func recovery(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recovery a panic : %v\n", err)
			context.String(http.StatusInternalServerError, "Server Internal Error")
			context.Abort()
		}
	}()
	context.Next()
}

func jwtAuthorization(context *gin.Context) {
	authHeader := context.Request.Header.Get("Authorization")
	authHeaderSplit := strings.Split(authHeader, " ")
	if len(authHeaderSplit) != 2 || strings.ToLower(authHeaderSplit[0]) != "bearer" {
		context.String(http.StatusBadRequest, "invalid authorization header format")
		context.Abort()
		return
	}
	authToken := authHeaderSplit[1]
	jwt, err := resolveJwt(authToken)
	if err != nil {
		if jwt != nil {
			context.String(http.StatusForbidden, err.Error())
		} else {
			context.String(http.StatusUnauthorized, err.Error())
		}
		context.Abort()
		return
	}
	context.Set("jwt", jwt)
}

func corsAccess() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowFiles = true
	return cors.New(config)
}

func errorHandler(ctx *gin.Context) {
	errors := ctx.Errors
	if len(errors.Errors()) > 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.JSON())
	}
}
