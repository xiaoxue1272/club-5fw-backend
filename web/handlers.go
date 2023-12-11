package web

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func requestLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		request := context.Request
		logger.Debugf(
			"Request income\nMethod: %s\nPath: %s\nContent Length: %d\nHeaders: %s\nCookies: %s\nRemote Address: %s\n",
			request.Method,
			request.RequestURI,
			request.ContentLength,
			request.Header,
			request.Cookies(),
			request.RemoteAddr,
		)
	}
}

func recovery() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Recovery a panic : %v\n", err)
				context.String(http.StatusInternalServerError, "Server Internal Error")
			}
		}()
		context.Next()
	}
}

func JwtAuthorization() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.Request.Header.Get("Authorization")
		authHeaderSplit := strings.Split(authHeader, " ")
		if len(authHeaderSplit) != 2 || strings.ToLower(authHeaderSplit[0]) != "bearer" {
			context.String(http.StatusBadRequest, "invalid authorization header format")
			context.Abort()
			return
		}
		authToken := authHeaderSplit[1]
		data, err := resolveJwt(authToken)
		if err != nil {
			if data != nil {
				context.String(http.StatusForbidden, err.Error())
			} else {
				context.String(http.StatusUnauthorized, err.Error())
			}
			context.Abort()
			return
		}
		context.Set("data", data)
	}
}
