package routers

import (
	"github.com/gin-gonic/gin"
	"software_bottle/controller"
	"software_bottle/middlewares"
	"software_bottle/websocket"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	r.GET("/source", middlewares.ValidDataTokenMiddleWare, controller.Safe)
	r.POST("/uploadImg", middlewares.ValidDataTokenMiddleWare, controller.UploadImag)
	r.GET("/getImg", middlewares.ValidDataTokenMiddleWare, controller.SendImg)
	r.POST("/sendBottle", middlewares.ValidDataTokenMiddleWare, controller.UploadBottle)
	r.GET("/getBottle", middlewares.ValidDataTokenMiddleWare, controller.SendBottle)
	r.GET("/ws", middlewares.ValidDataTokenMiddleWare, websocket.WsHandle)
	r.GET("/getOfflineMsg", middlewares.ValidDataTokenMiddleWare, websocket.SendOfflineMsg)
	return r
}
