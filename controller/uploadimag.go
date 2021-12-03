package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"software_bottle/model"
	"strconv"
	"strings"
	"time"
)

//
func UploadImag(c *gin.Context) {
	//解析图片，数据保存到本地，路径持久化到数据库
	imgFile, imgHandler, err := c.Request.FormFile("img")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	defer func() {
		_ = imgFile.Close()
	}()
	current := time.Now()
	//时间戳命名图片，防止重复
	ext := strings.ToLower(path.Ext(imgHandler.Filename))
	if !(ext == ".jpg" || ext == ".png" || ext == ".bmp" || ext == ".jpeg") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "图片格式只支持.jpg,.png,.bmp,.jpeg",
		})
		return
	}
	imgPath := model.IMAG_PATH + strconv.FormatInt(current.Unix(), 10) + ext
	f, err := os.Create(imgPath)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func() {
		_ = f.Close()
	}()
	//将用户上传的图片保存到服务器本地
	_, err = io.Copy(f, imgFile)
	if err != nil {
		uploadError(err, c)
		return
	}
	username, ok := c.Get("user_name")
	if !ok {
		fmt.Println("Not get username")
		return
	}

	var user model.User
	err = model.GetAUser(&user, username.(string))
	if err != nil {
		uploadError(err, c)
		return
	}
	//删除服务器本地原本用户头像
	if user.ImagePath != "" && user.ImagePath != (model.IMAG_PATH+"default.jpg") {
		_ = os.Remove(user.ImagePath)
	}
	user.ImagePath = imgPath
	err = model.UpdateUser(&user)
	if err != nil {
		uploadError(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "img post success",
	})

}

func uploadError(err error, c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code": http.StatusInternalServerError,
		"msg":  "uploadImg error:" + err.Error(),
	})
}
