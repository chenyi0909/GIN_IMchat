package main

import (
	"GIN_IMchat/models"
	"GIN_IMchat/router"
	"GIN_IMchat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	if utils.DB != nil {
		utils.DB.AutoMigrate(&models.UserBasic{})
		utils.DB.AutoMigrate(&models.Message{})
		utils.DB.AutoMigrate(&models.GroupBasic{})
	}
	utils.InitRedis()

	r := router.Router()
	r.Run()
}
