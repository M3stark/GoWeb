package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Session 中间件，处理session
func Session(sessionName string, sessionKey string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(sessionKey))
	// 使用Cookie作为Session的存储方式
	store.Options(sessions.Options{
		MaxAge: 300, // seconds
		Path:   "/",
	})
	return sessions.Sessions(sessionName, store)
}
