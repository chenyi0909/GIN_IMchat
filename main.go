package main

import (
	"GIN_IMchat/router"
	"GIN_IMchat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()

	r := router.Router()
	r.Run()
}
