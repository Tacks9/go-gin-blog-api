package util

import (
	"go-gin-blog-api/pkg/setting"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(setting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	// jwt 所需要的标准字段
	jwt.StandardClaims
}

// 生成 token
func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "go-gin-blog-api",
		},
	}

	// 创建内部加密方式
	// HMAC SHA-256  采用 SHA-256 哈希函数和一个密钥来生成签名
	// 具体过程是将JWT的头部和有效载荷进行Base64编码，然后使用一个密钥对编码后的字符串进行签名
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成签名
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// 解析 token
func ParseToken(token string) (*Claims, error) {
	// 解析鉴权
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 检验超时时间
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
