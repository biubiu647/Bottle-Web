package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	_ "software_bottle/config"
	"software_bottle/dao"
	"software_bottle/model"
	"software_bottle/routers"
)

func main() {
	gin.DefaultWriter = colorable.NewColorableStdout()
	dao.InitRedis()
	err := dao.InitMySQL()
	if err != nil {
		panic(err)
	}
	defer dao.DB.Close()
	//绑定映射模型与数据库
	dao.DB.AutoMigrate(&model.User{})
	//设置路由
	r := routers.Router()
	r.Run(":9999")
}
