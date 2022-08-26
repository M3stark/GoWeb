package libs

import (
	"GoWebProject/dao"
	"errors"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

/*
	IsValid 验证密码是否有效，必须包含大小写字母、数字、特殊字符.或*
	----------------+-------------------
	     errCode    |      err
	----------------+-------------------
			-1		    密码不能为空
	----------------+-------------------
			-2 			缺少大写字母
	----------------+-------------------
			-3 			缺少小写字母
	----------------+-------------------
			-4 			缺少数字
	----------------+-------------------
			-5			缺少特殊字符(./* ）
	----------------+-------------------
*/
func IsValid(password string) (errCode int, err error) {

	if len(password) == 0 {
		return -1, errors.New("密码不能为空")
	}

	// 统计大写字母、小写字母、数字、特殊字符个数
	capitalLetterCount, LowerLetterCount, numberCount, SpecialCount := 0, 0, 0, 0

	for _, v := range password {
		switch {
		case 'A' <= v && v <= 'Z':
			capitalLetterCount++
		case 'a' <= v && v <= 'z':
			LowerLetterCount++
		case '0' <= v && v <= '9':
			numberCount++
		case v == '.' || v == '*':
			SpecialCount++
		default:
		}
	}

	if capitalLetterCount == 0 {
		return -2, errors.New("缺少大写字母")
	} else if LowerLetterCount == 0 {
		return -3, errors.New("缺少小写字母")
	} else if numberCount == 0 {
		return -4, errors.New("缺少数字")
	} else if SpecialCount == 0 {
		return -5, errors.New("缺少特殊字符.或*")
	}

	return 0, nil
}

func CheckCaptcha(context *gin.Context) (errCode int) {
	/*
		验证码处理逻辑：
		1、读取客户端POST请求的验证码；
		2、在redis中获取CaptchaId；
		3、验证验证码是否输入正确；
	*/
	reqCaptcha := context.PostForm("captcha")
	getCaptcha := dao.Redis.Get(context, "CaptchaId")
	reqCaptchaId, reqErr := getCaptcha.Result()

	if reqErr == redis.Nil {
		return -1
	} else if reqErr != nil {
		return -2
	}

	if !captcha.VerifyString(reqCaptchaId, reqCaptcha) {
		return -3
	}
	return 0
}
