package controller

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"io"
	"math/rand"
	"net/http"
	"os"
	"software_bottle/dao"
	"software_bottle/model"
	"strconv"
	"strings"
	"time"
)

func SendImg(c *gin.Context) {

	userId, ok := c.Get("user_id")
	if !ok {

	}
	var user model.User
	err := model.GetAUserById(&user, int(userId.(float64)))
	if err != nil {

	}
	//拷贝图片
	f, err := os.Open(user.ImagePath)
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, f); err != nil {
		fmt.Println(err.Error())
	}
	//获取图片格式
	slice := strings.Split(user.ImagePath, ".")
	ext := slice[len(slice)-1]
	c.JSON(http.StatusOK, gin.H{
		"code":       http.StatusOK,
		"msg":        "success get Img",
		"img_base64": "data:image/" + ext + ";base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
		"img_type":   ext,
	})
}

func SendBottle(c *gin.Context) {
	conn := dao.RedisCache.Conn()
	defer conn.Close()
	id, _ := c.Get("user_id")
	userId := strconv.Itoa(int(id.(float64)))
	//删除过期漂流瓶
	DelExpireMembers(conn)
	rand.Seed(time.Now().Unix())
	//随机获取一个漂流瓶
	num, _ := redis.Int64(conn.Do("ZCARD", dao.Field[2]))
	//当存在漂流瓶时
	if num > 0 {
		//尝试100次，直到拿到其他人发送的漂流瓶则直接返回,否则在100次没取到之后就返回无漂流瓶可取.
		for i := 0; i < 100; i++ {
			index := rand.Int63n(num)
			val, err := redis.Values(conn.Do("ZRANGE", dao.Field[2], index, index, "WITHSCORES"))
			key := string(val[0].([]byte))
			//漂流瓶owner的ID
			s := strings.Split(key, "-")[0]
			//防止获取到自己的漂流瓶
			if s == userId {
				continue
			}
			senderId, _ := strconv.Atoi(s)

			var sender model.User
			_ = model.GetAUserById(&sender, senderId)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": http.StatusInternalServerError,
					"msg":  "获取漂流瓶失败",
				})
				return
			}
			value, err := redis.Values(conn.Do("HMGET", key, dao.Field[0]))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": http.StatusInternalServerError,
					"msg":  "获取漂流瓶内容失败",
				})
				return
			}

			imagePath := sender.ImagePath

			msg := string(value[0].([]byte))
			f, err := os.Open(imagePath)
			buf := bytes.NewBuffer(nil)
			if _, err = io.Copy(buf, f); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": http.StatusInternalServerError,
					"msg":  "获取漂流瓶用户头像失败",
				})
				return
			}
			slice := strings.Split(sender.ImagePath, ".")
			ext := slice[len(slice)-1]
			c.JSON(http.StatusOK, gin.H{
				"code":       http.StatusOK,
				"message":    msg,
				"user_name":  sender.UserName,
				"user_id":    sender.UserId,
				"img_base64": "data:image/" + ext + ";base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
				"img_type":   ext,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusInternalServerError,
		"msg":  "目前没有漂流瓶可以捞，请稍后重试",
	})
}

func DelExpireMembers(conn redis.Conn) {
	stop := time.Now().Unix() - EXPIRE
	//fmt.Println(stop)
	//删除一个月以前的漂流瓶
	_, _ = conn.Do("ZREMRANGEBYSCORE", dao.Field[2], 0, stop)
}
