package libs

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

/*
	JWT TOKEN 包括:
		- Header: 存储所使用的加密算法和Token类型；
		- Payload: 负载，也是一个JSON对象，规定了7个官方字段(签发，生效时间, etc..)
		- SIGNATURE: 对上述部分对签名，防止数据被串改。
*/

// Myclaims 自定义结构体，可以保存额外的字段 - username
type Myclaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Secret 指定一个密钥Secret，只能服务器知道
var Secret = []byte("August")

// TokenDuration JWT的过期时间
const TokenDuration = time.Hour * 2

// GenToken 生成token
func GenToken(username string) (string, error) {
	// 创建一个自己对声明
	claim := Myclaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenDuration).Unix(), // 过期时间
			Issuer:    "test-jwt",                           // 签发人
		},
	}

	// Access Token - 有效期较短的token
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// 使用指定的Secret签名并获得token字符串
	return token.SignedString(Secret)
}

// ParserToken 解析JWT
func ParserToken(tokenStr string) (*Myclaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Myclaims{}, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	// 校验token
	if claim, ok := token.Claims.(*Myclaims); ok && token.Valid {
		return claim, nil
	}
	return nil, errors.New("invalid token")
}

//TODO 刷新 token
