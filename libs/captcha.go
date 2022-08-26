package libs

import (
	"GoWebProject/dao"
	"bytes"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// GenCaptcha 响应图形验证码
func GenCaptcha(context *gin.Context, length, width, height int) {
	captchaId := captcha.NewLen(length)

	// 将"CaptchaId"保存到Redis，有效期5分钟
	if err := dao.Redis.Set(context, "CaptchaId", captchaId, time.Minute*5).Err(); err != nil {
		fmt.Println(err)
	}

	if captchaId == "" {
		fmt.Println("CaptchaId不存在！")
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "CaptchaId不存在！",
		})
		return
	}

	var w bytes.Buffer
	if err := captcha.WriteImage(&w, captchaId, width, height); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "Captcha写入失败！",
		})
		return
	}

	context.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	context.Writer.Header().Set("Pragma", "no-cache")
	context.Writer.Header().Set("Expires", "0")
	context.Writer.Header().Set("Content-Type", "image/png")
	http.ServeContent(context.Writer, context.Request, captchaId, time.Time{}, bytes.NewReader(w.Bytes()))
	return
}
