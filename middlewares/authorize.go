package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SecretKey = "XXXXXXXXXXX123456"
)

//权限验证中间件
func ValidDataTokenMiddleWare(c *gin.Context) {
	token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

	if err != nil {
		//log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "Unauthorized access",
		})
		c.Abort()
		return
	} else {
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)

			c.Set("user_name", claims["UserName"])
			c.Set("user_id", claims["UserId"])

			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "token is not valid",
			})
			c.Abort()
			return
		}
	}
}
