package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// Session 中间件，处理session
func Session(sessionName string, sessionKey string) gin.HandlerFunc {
	/*
		Session信息是存放在server端，但session id是存放在client cookie的.
		使用Redis保存session，调用session.Save()，可以将sessionID保存至Redis.
	*/
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte(sessionKey))

	// 使用Cookie作为Session的存储方式
	//store := cookie.NewStore([]byte(sessionKey))
	// TODO 为什么用Cookie无法获取sessionID？

	store.Options(sessions.Options{
		MaxAge: 300, // seconds
		Path:   "/",
	})
	return sessions.Sessions(sessionName, store)
}
