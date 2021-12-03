package controller

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"software_bottle/model"
	"time"
)

const (
	SecretKey = "XXXXXXXXXXX123456"
)

type User struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

//登录
func Login(c *gin.Context) {
	var user User
	err := c.BindJSON(&user)
	//获取json失败
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code": http.StatusForbidden,
			"msg":  "error in decode request",
		})
		return
	}
	var uuser model.User
	//查不到记录
	if err = model.GetAUser(&uuser, user.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	} else {
		//密码不对
		if uuser.Password != user.Password {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  errors.New("password error").Error(),
			})
			return
		}
	}
	token := jwt.New(jwt.SigningMethodHS256)
	var claims = make(jwt.MapClaims)
	claims["UserName"] = uuser.UserName
	claims["UserId"] = uuser.UserId
	//token过期时间
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(24)).Unix()
	//token创建时间
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	//加密生成token
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Sign token error",
		})
		return
	}

	//返回token
	f, err := os.Open(uuser.ImagePath)
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, f); err != nil {
		fmt.Println("failed to copy picture")
	}
	c.JSON(http.StatusOK, gin.H{
		"code":       http.StatusOK,
		"msg":        "Success",
		"token":      tokenString,
		"img_base64": base64.StdEncoding.EncodeToString(buf.Bytes()),
	})
}

func Safe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "authorized success",
	})
}
