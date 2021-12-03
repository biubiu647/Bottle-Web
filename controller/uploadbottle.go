package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"software_bottle/dao"
	"software_bottle/model"
	"strconv"
	"time"
)

/*
redis Hash(存放漂流瓶信息,头像)
key:user_id-时间戳
field0:msg
value:message

redis sorted set(存放漂流瓶key)
key:bottle_owner
member:user_id-时间戳

随机取，如果漂流瓶得到回复，就从redis中删除,回复的话就建立一对一聊天了

一对一聊天
如果对方离线，则把消息方入redis中，等对方上线再从redis中取出发送，然后删除redis
*/

type Bottle struct {
	Msg string `json:"message"`
}

const (
	EXPIRE = int64(1000 * 3600 * 24 * 30)
)

func UploadBottle(c *gin.Context) {
	var bottle Bottle
	err := c.BindJSON(&bottle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "can not get message:" + err.Error(),
		})
		return
	}
	fmt.Println(bottle)
	//int
	v, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "user_id isn't exists",
		})
		return
	}
	userId := int(v.(float64))
	cur := time.Now().Unix()
	key := strconv.Itoa(userId) + "-" + strconv.FormatInt(cur, 10)
	var user model.User
	if err = model.GetAUserById(&user, userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}
	conn := dao.RedisCache.Conn()
	defer conn.Close()
	//事务
	conn.Send("MULTI")
	conn.Send("HSET", key, dao.Field[0], bottle.Msg)
	conn.Send("ZADD", dao.Field[2], cur, key)
	conn.Send("EXPIRE", key, EXPIRE)
	conn.Send("PERSIST", dao.Field[2])
	_, err = conn.Do("EXEC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "redis save error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "send bottle success",
	})
}
