package main

import (
	"GIN_IMchat/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/GIN_IMchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	//db.AutoMigrate(&models.UserBasic{}) //如果没有表，则会创建该表
	//db.AutoMigrate(&models.Contact{},&models.GroupBasic{})
	db.AutoMigrate(&models.Message{})
	// Create
	// user := &models.UserBasic{}
	// user.Name = "陈忆"
	// //db.Create(user)

	// // Read
	// fmt.Println(db.First(user, 9)) // 根据整型主键查找,并终端打印
	// //db.First(user, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// // Update - 将 product 的 price 更新为 200
	// db.Model(user).Update("Passwd", "1234")
	// // Update - 更新多个字段
	// //db.Model(user).Updates(UserBasic{Price: 200, Code: "F42"}) // 仅更新非零值字段
	// //db.Model(user).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// // Delete - 删除 product
	// //db.Delete(user, 8)
}
