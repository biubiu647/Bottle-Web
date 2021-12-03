package websocket

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"io"
	"net/http"
	"os"
	"software_bottle/dao"
	"software_bottle/model"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"log"
)

type Hub struct {
	register   chan *Client
	unregister chan *Client
	msgChan    chan []byte
	clienlist  map[int]*Client
}

type Client struct {
	id   int
	conn *websocket.Conn
}

type MsgContent struct {
	Sender   int    `json:"sender"`
	Receiver int    `json:"receiver"`
	Message  string `json:"message"`
}
type ReplyContent struct {
	UserId    int    `json:"user_id"`
	ImgBase64 string `json:"img_base64"`
	Message   string `json:"message"`
}

var hub = &Hub{
	register:   make(chan *Client),
	unregister: make(chan *Client),
	msgChan:    make(chan []byte),
	clienlist:  make(map[int]*Client),
}

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var Filed = []string{"user_id", "img_base64", "message"}

func (h *Hub) Run() {
	for {
		//依次为注册，注销，转发
		select {
		case c := <-h.register:
			log.Printf("客户端 %d 加入\n", c.id)
			h.clienlist[c.id] = c
		case c := <-h.unregister:
			log.Printf("客户端 %d 退出\n", c.id)
			delete(h.clienlist, c.id)
		case v := <-h.msgChan:
			msg := &MsgContent{}
			_ = json.Unmarshal(v, msg)
			if client, ok := h.clienlist[msg.Receiver]; ok {
				log.Printf("转发自客户端 %d 向客户端 %d 发送的消息: %s\n",
					msg.Sender, msg.Receiver, msg.Message)
				client.WriteMsg(v)
			} else {
				//接收方离线的处理，存在redis中
				SaveOfflineMsg(v)
				log.Println("用户已离线")
			}
		}
	}
}

//获取离线时未接收的回复
func SaveOfflineMsg(v []byte) {
	conn := dao.RedisCache.Conn()
	defer conn.Close()
	var reply ReplyContent
	var msg MsgContent
	_ = json.Unmarshal(v, &msg)
	err := GetReply(v, &reply)
	if err != nil {
		log.Println("write OfflineMsg error:", err)
		return
	}
	key := strconv.Itoa(msg.Receiver)
	replyBytes, _ := json.Marshal(reply)
	conn.Send("MULTI")
	conn.Send("SADD", key, replyBytes)
	_, err = conn.Do("EXEC")
	if err != nil {
		log.Println("save OfflineMsg in redis error:", err)
		return
	}
}
func (c *Client) WriteMsg(v []byte) {

	var reply ReplyContent
	err := GetReply(v, &reply)
	if err != nil {
		log.Println("write Msg error:", err)
		return
	}
	replyBytes, _ := json.Marshal(reply)
	err = c.conn.WriteMessage(websocket.TextMessage, replyBytes)
	if err != nil {
		log.Println("write Msg error:", err)
		return
	}
	log.Printf("发送到客户端的消息: %s\n", v)
}

//将sender发送的消息进行处理
func GetReply(v []byte, reply *ReplyContent) error {
	var msg MsgContent
	_ = json.Unmarshal(v, &msg)
	var user model.User
	err := model.GetAUserById(&user, msg.Sender)
	if err != nil {
		//log.Println("获取用户信息失败", err.Error())
		return err
	}

	imagePath := user.ImagePath
	f, err := os.Open(imagePath)
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, f); err != nil {
		//log.Println("获取头像失败", err.Error())
		return err
	}
	slice := strings.Split(user.ImagePath, ".")
	ext := slice[len(slice)-1]
	reply.UserId = user.UserId
	reply.Message = msg.Message
	reply.ImgBase64 = "data:image/" + ext + ";base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	return nil
}

func (c *Client) ReadMsg() {
	defer func() {
		hub.unregister <- c
		_ = c.conn.Close()
	}()
	//客户端不断进行读取
	for {
		_, v, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("%d read Msg error:%s\n", c.id, err.Error())
			break
		}
		//为消息添加上发送者的id
		msg := &MsgContent{}
		_ = json.Unmarshal(v, msg)
		msg.Sender = c.id

		v, _ = json.Marshal(msg)
		//转发消息
		hub.msgChan <- v
	}
}

func WsHandle(c *gin.Context) {
	go hub.Run()
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	v, _ := c.Get("user_id")
	uid := int(v.(float64))
	client := &Client{
		id:   uid,
		conn: conn,
	}
	hub.register <- client
	//为每个客户端建立一个读取协程
	go client.ReadMsg()
}

func SendOfflineMsg(c *gin.Context) {
	conn := dao.RedisCache.Conn()
	defer conn.Close()
	v, _ := c.Get("user_id")
	userId := int(v.(float64))
	key := strconv.Itoa(userId)
	replyBytes, err := redis.Values(conn.Do("SMEMBERS", key))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusOK,
			"msg":  err.Error(),
		})
		return
	}
	_, _ = conn.Do("DEL", key)
	for _, val := range replyBytes {
		//fmt.Println(index)
		var reply ReplyContent
		replyByte := val.([]byte)
		err = json.Unmarshal(replyByte, &reply)
		if err != nil {
			//c.JSON(http.StatusInternalServerError, gin.H{
			//	"code": http.StatusInternalServerError,
			//	"msg":  err.Error(),
			//})
		} else {
			c.JSON(http.StatusOK, gin.H{
				//"code":  http.StatusOK,
				"reply": reply,
			})
		}
	}
}
