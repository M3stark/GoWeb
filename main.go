package main

import (
	"GoWebProject/dao"
	"GoWebProject/models"
	"GoWebProject/routers"
)

func main() {
	// 创建数据库
	// sql: CREATE DATABASE myWeb;

	// 连接数据库
	err := dao.InitMySQL()
	if err != nil {
		panic(err)
	}

	defer dao.Close() // 程序退出关闭连接

	// 连接Redis
	dao.Redis = dao.InitRedis()
	defer dao.RedisClose()

	// 模型绑定
	err = dao.DB.AutoMigrate(&models.User{})
	if err != nil {
		return
	}

	router := routers.SetUpRouter()

	router.Run(":8001")
}
