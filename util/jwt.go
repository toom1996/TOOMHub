// @Description
// @Author    2020/8/20 10:01
package util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//var jwtSecret = []byte("toomhub")

type Claims struct {
	MiniId    int64
	CreatedAt int64
	jwt.StandardClaims
}

var identity *Claims

func GenerateToken(id int64) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(60 * time.Minute * 24 * 7)

	claims := Claims{
		id,
		nowTime.Unix(),
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "toomhub",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(GetConfig().Jwt.Secret))
	fmt.Println(err)
	return token, err
}

func ParseToken(token string, c *gin.Context) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetConfig().Jwt.Secret), nil
	})

	fmt.Println(tokenClaims)
	if tokenClaims != nil {
		if tokenClaims.Valid {
			if claims, ok := tokenClaims.Claims.(*Claims); ok {
				identity = &Claims{
					MiniId:    claims.MiniId,
					CreatedAt: claims.CreatedAt,
				}
				return claims, nil
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// token 过期了
				c.JSON(http.StatusOK, map[string]interface{}{
					"code": 401,
					"msg":  "token is expired",
				})
				c.Abort()
				return nil, nil
			}
		}
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 401,
		"msg":  "bad token",
	})
	c.Abort()
	return nil, nil
}

// @title 获取用户信息
func GetIdentity() *Claims {
	fmt.Println(identity)
	return identity
}
