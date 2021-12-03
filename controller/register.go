package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"software_bottle/model"
)

//注册
func Register(c *gin.Context) {
	var user model.User
	err := c.BindJSON(&user)
	if len(user.UserName) > 24 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "username can't more than eight Chinese characters",
		})
		return
	}
	if len(user.Password) < 8 || len(user.Password) > 15 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "password's length should between 8 to 15",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}
	//用户名已存在
	if err = model.UserIsExists(user.UserName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}
	user.ImagePath = model.IMAG_PATH + model.DefaultImg
	if err = model.CreateAUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
	} else {

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"msg":     "Success",
			"user_id": user.UserId,
			//"data": user,
		})
	}
}
