package service

import (
	"github.com/gin-gonic/gin"
)

// FindUserByName
// @Tags 首页
// @Success 200 {string} json{"code","message"}
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "welcome!!!",
	})
}
