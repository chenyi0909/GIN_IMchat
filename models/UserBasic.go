package models

import (
	"GIN_IMchat/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Passwd        string
	Salt          string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIP      string
	ClientPort    string
	LoginTime     time.Time
	HeartbeatTime time.Time
	LogoutTime    time.Time
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

func FindUserByNameAndPwd(name, passwd string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and passwd = ?", name, passwd).First(&user)
	//设置登录Tocken
	str := fmt.Sprintf("%d", time.Now().Unix())
	tocken := utils.MD5Encode(str)
	utils.DB.Model(&user).Update("identity", tocken)
	return user
}
func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic, id interface{}) *gorm.DB {
	return utils.DB.Delete(&user, id) //根据主键值删除
}

func UpdateUser(user UserBasic, newValue string) *gorm.DB {
	return utils.DB.Model(&user).Update("Passwd", newValue)
}

func UpdateUser2(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Passwd: user.Passwd, Name: user.Name, Email: user.Email})
}
