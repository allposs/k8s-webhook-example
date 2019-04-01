package encrypt

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var tokenCustomizeKey = "deppon.com"

// Claims jwt claims
type Claims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

// GenerateToken 生成token,uid用户id，expireSec过期秒数
func GenerateToken(user string, expireSec int) (tokenStr string, err error) {
	sec := time.Duration(expireSec)
	claims := Claims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * sec).Unix(), //自定义有效期，过期需要重新登录获取token,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用自定义字符串加密 and get the complete encoded token as a string
	tokenStr, err = token.SignedString([]byte(tokenCustomizeKey))
	return tokenStr, err
}

// ParseToken 解析Token, token是token值，string类型；key是自定义字条串，string类型
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenCustomizeKey), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
