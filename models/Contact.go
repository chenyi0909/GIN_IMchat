package models

import (
	"gorm.io/gorm"
)

// Contact 人员关系结构体
type Contact struct {
	gorm.Model
	OwnerID  uint   //谁的关系信息
	TargetID uint   //对应的谁
	Type     int    //对应的类型0  1   2
	Desc     string //预留的描述信息
}

func (table *Contact) TableName() string {
	return "contact"
}
