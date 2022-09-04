package routers

import (
	"GoWebProject/controller"
	"GoWebProject/libs"
	"GoWebProject/middleware"
	"github.com/gin-gonic/gin"
)

/*
	SetUpRouter 注册路由
*/
func SetUpRouter() *gin.Engine {
	engine := gin.Default()

	// 告诉 gin 框架，模板文件引用的静态文件去哪找
	//router.Static("static", "static") // static 保存 css 和 js 静态文件

	engine.LoadHTMLGlob("templates/**/*")

	// 使用Session中间件，以保存token
	engine.Use(middleware.Session("cookie-style-session", "secret"))

	// 首页
	engine.GET("/", controller.IndexHandler)

	// 生成验证码
	engine.GET("captcha", func(context *gin.Context) {
		// 验证码位数，验证码图片宽度和高度
		libs.GenCaptcha(context, 4, 80, 30)
	})

	// 注册
	engine.GET("/register", controller.Register)
	engine.POST("/register", controller.RegisterHandler)

	// 登陆
	engine.GET("/login", controller.Login)
	engine.POST("/login", controller.LoginHandler)

	// home page
	engine.GET("/home", middleware.JWTAuth(), controller.HomeHandler)

	return engine
}
