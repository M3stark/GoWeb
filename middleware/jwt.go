package middleware

import (
	"GoWebProject/libs"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// JWTAuth jwt认证中间件
func JWTAuth() func(context *gin.Context) {
	return func(context *gin.Context) {
		// 从session中获取token
		session := sessions.Default(context)
		if session.Get("token") == nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token已失效",
			})
			context.Abort()
			return
		}

		authHeader := session.Get("token").(string)
		if authHeader == "" {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "请求的auth为空",
			})
			context.Abort()
			return
		}

		reqToken, err := libs.ParserToken(authHeader)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token验证失败",
			})
			context.Abort()
			return
		}

		context.Set("username", reqToken.Username)
		// 后面通过 context.Get("username")获取当前请求的用户信息

		context.Next()
	}
}
