package controller

import (
	"GoWebProject/dao"
	"GoWebProject/libs"
	"GoWebProject/models"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// IndexHandler 加载首页
func IndexHandler(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", nil)
	return
}

// Register 用户注册页面
func Register(context *gin.Context) {
	context.HTML(http.StatusOK, "register.html", gin.H{})
	return
}

// RegisterHandler 用户注册
func RegisterHandler(context *gin.Context) {
	user := models.User{}
	err := context.Bind(&user)
	if err != nil {
		return
	}

	username := user.Username
	password := user.Password
	email := user.Email
	mobile := user.Mobile

	// 处理逻辑
	// 如果用户名为空
	if len(username) == 0 {
		context.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "用户名不能为空！",
		})
		return
	}
	// 1、用户名要验证唯一
	dao.DB.Where("username = ?", username).First(&user)
	if user.ID != 0 {
		context.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "当前用户名已被注册，请换一个吧！",
		})
		return
	}

	//2、密码需要复杂度验证
	if errCode, err := libs.IsValid(password); err != nil {
		switch errCode {
		case -1:
			context.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "密码不能为空！",
			})
			return
		case -2:
			context.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "密码缺少大写字母！",
			})
			return
		case -3:
			context.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "密码缺少小写字母！",
			})
			return
		case -4:
			context.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "密码缺少数字！",
			})
			return
		case -5:
			context.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "密码缺少特殊字符.或* ",
			})
			return
		}
	}

	// 3、密码要加密存储 - MD5算法
	// 注册的时候用md5算法加密并保存，登陆的时候将请求信息的密码md5加密，判断是否一致
	password = libs.GetMD5String(password)

	// 验证手机号是否为11位数
	if len(mobile) != 11 {
		context.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "手机号不是11位数",
		})
		return
	}

	insertUser := &models.User{
		Username: username,
		Password: password,
		Email:    email,
		Mobile:   mobile,
	}

	// 插入
	if err := dao.DB.Create(insertUser).Error; err != nil {
		fmt.Println("插入失败！", err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功！",
	})
	//context.HTML(http.StatusOK, "login.html", nil)
	return
}

// Login 用户登陆页面
func Login(context *gin.Context) {
	// Ref: https://segmentfault.com/q/1010000000308417
	// 将登陆错误次数记录在 Redis 中，每次发起GET请求进入登陆页面，
	// 登陆错误次数记为0
	loginFailedCount := 0
	err := dao.Redis.Set(context, "loginFailedCount", loginFailedCount, time.Minute*5).Err()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "登陆错误次数载入Redis失败",
		})
		return
	} else {
		fmt.Println("登陆错误次数已成功载入Redis")
	}

	// 每次手动刷新登陆页面，错误登陆次数将重置为0
	context.HTML(http.StatusOK, "login.html", gin.H{})
}

// LoginHandler 用户登陆
func LoginHandler(context *gin.Context) {
	// 从 Redis中取出登陆失败次数
	loginFailedCount, _ := dao.Redis.Get(context, "loginFailedCount").Result()
	//if getRedisErr != nil {
	//	panic(getRedisErr)
	//}
	failedCount, _ := strconv.Atoi(loginFailedCount)

	LoginUser := models.User{}
	err := context.Bind(&LoginUser)
	if err != nil {
		return
	}

	username := LoginUser.Username
	password := LoginUser.Password

	if err := context.ShouldBind(&LoginUser); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误！",
		})
		return
	}

	user := models.User{}

	// 用户名不存在
	dao.DB.Where("username = ?", username).First(&user)

	if user.ID == 0 {
		failedCount++ // 用户名不存在， 登陆错误+1
		dao.Redis.Set(context, "loginFailedCount", failedCount, time.Minute*5)
		if failedCount == 5 {
			// 如果输错超过4次，重新加载页面，刷新验证码
			context.HTML(http.StatusForbidden, "relogin.html", gin.H{
				"code": 403,
				"msg":  "已经输错5次，请重新登陆吧！",
			})
		} else {
			context.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "用户名不存在！",
			})
		}
		return
	}

	// 密码输入错误
	password = libs.GetMD5String(password)
	pErr := dao.DB.Where("password = ?", password).First(&user).Error
	if errors.Is(pErr, gorm.ErrRecordNotFound) {
		failedCount++ // 密码输入错误， 登陆错误+1
		dao.Redis.Set(context, "loginFailedCount", failedCount, time.Minute*5)
		if failedCount == 5 {
			// 如果输错超过4次，重新加载页面，刷新验证码
			context.HTML(http.StatusForbidden, "relogin.html", gin.H{
				"code": 403,
				"msg":  "已经输错5次，请重新登陆吧！",
			})
		} else {
			context.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "密码错误！",
			})
		}
		return
	}

	// 如果输错次数超过4次，验证码还输错了
	if failedCount >= 5 {
		checkCode := libs.CheckCaptcha(context)
		switch checkCode {
		case -1:
			context.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "未找到CaptchaId",
			})
			return
		case -2:
			context.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "未知错误",
			})
			return
		case -3:
			context.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "验证码输入错误",
			})
			return
		}
	}

	// 生成 JWT token
	token, tokenErr := libs.GenToken(username)
	if tokenErr != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "错误的token",
		})
		return
	}

	// 将token保存到session中
	session := sessions.Default(context)
	session.Set("token", token)
	session.Save()

	// 进入成功页面
	context.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "success",
		"token": token, // token string
	})
	return
}

// HomeHandler 用户登陆
func HomeHandler(context *gin.Context) {
	username, _ := context.Get("username")
	context.HTML(http.StatusOK, "homepage.html", gin.H{
		"usr": username,
	})
	return
}
