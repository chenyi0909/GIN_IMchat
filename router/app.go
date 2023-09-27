package router

import (
	"GIN_IMchat/docs"
	"GIN_IMchat/service"

	//docs "github.com/GIN_IMcaht/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/index", service.GetIndex)

	//传统API--获取用户列表
	//r.GET("/user/GetUserList", service.GetUserList)
	//RESTful API--获取用户列表
	r.GET("/user", service.GetUserList)

	//传统API--添加用户
	//r.POST("/user/CreateUser", service.CreateUser)
	//RESTful API--添加用户
	r.POST("/user", service.CreateUser)

	//传统API--删除用户
	//r.GET("/user/DeleteUser", service.DeleteUser)
	//RESTful API--删除全部用户用户
	r.DELETE("/user", service.DeleteUser)

	//传统API--修改某个用户
	//r.POST("/user/UpdateUser", service.UpdateUser)
	//RESTful API--修改某个用户
	r.PUT("/user", service.UpdateUser)

	//传统API--获取单个用户信息
	//r.GET("/user/FindUser", service.FindUserByName)
	//RESTful API--获取单个用户信息
	r.GET("/oneuser", service.FindUserByName)

	r.POST("/user/login", service.FindUserByNameAndPwd)

	//发送消息
	r.GET("user/sendMsg", service.SendMsg)
	r.GET("user/sendUserMsg", service.SendUserMsg)
	return r
}
