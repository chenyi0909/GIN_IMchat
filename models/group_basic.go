package models

import (
	"gorm.io/gorm"
)

// GroupBasic 群信息结构体
type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerID uint
	Icon    string
	Type    int
	Desc    string
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
