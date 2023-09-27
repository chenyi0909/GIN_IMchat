package service

import (
	"GIN_IMchat/models"
	"GIN_IMchat/utils"
	"github.com/gorilla/websocket"
	"time"

	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// GetUserList
// @Summary 查询全部用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()

	c.JSON(200, gin.H{
		"message": data,
	})
}

// FindUserByNameAndPwd
// @Summary 登录用户
// @Tags 用户模块
// @param username formData string false "用户名"
// @param password formData string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/login [post]
func FindUserByNameAndPwd(c *gin.Context) {
	user := models.UserBasic{}
	username := c.PostForm("username")
	passwd := c.PostForm("password")

	user = models.FindUserByName(username)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "用户未注册",
		})
		return
	}

	flag := utils.ValidPassword(passwd, user.Salt, user.Passwd)
	if !flag {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "用户名或密码错误",
		})
		return
	}

	user = models.FindUserByNameAndPwd(username, user.Passwd)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"user":    user,
	})
}

// FindUserByName
// @Summary 查询用户
// @Tags 用户模块
// @param name query string false "用户名"
// @Success 200 {string} json{"code","message"}
// @Router /oneuser [get]
func FindUserByName(c *gin.Context) {
	username := c.Query("name")
	data := models.FindUserByName(username)
	if data.ID != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": data,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无查询结果！！！",
		})
	}
}

// CreateUser
// @Summary 新增用户(用户注册)
// @Tags 用户模块
// @param username formData string false "用户名"
// @param password formData string false "密码"
// @param repassword formData string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	//user.Name =  c.Query("name")
	//user.Name = "ruan"
	//

	/*	username := c.Query("username")
		passwd := c.Query("password")
		repasswd := c.Query("repassword")*/
	username := c.PostForm("username")
	passwd := c.PostForm("password")
	repasswd := c.PostForm("repassword")
	if passwd != repasswd {
		c.JSON(-1, gin.H{
			"message":    "两次密码不一致",
			"username":   username,
			"password":   passwd,
			"repassword": repasswd,
		})
		return
	}
	//重复注册校验
	data := models.FindUserByName(username)
	if data.ID != 0 {
		c.JSON(-1, gin.H{
			"message":    "用户名以存在",
			"username":   username,
			"password":   passwd,
			"repassword": repasswd,
		})
		return
	}
	user.Name = username
	//user.Passwd = passwd
	//初始化随机数种子，以确保每次随机数都不同
	rand.Seed(time.Now().Unix())
	//获取随机数，用于生成密码
	salt := fmt.Sprintf("%06d", rand.Int31()) //获取该随机数的方法在包 “math/rand” 中
	user.Passwd = utils.MakePassword(passwd, salt)
	user.Salt = salt

	models.CreateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"message":    "新增用户成功！",
		"username":   username,
		"password":   passwd,
		"repassword": repasswd,
	})
}

// DeleteUser
// @Summary 注销用户
// @Tags 用户模块
// @param id query string false "user-id"
// @Success 200 {string} json{"code","message"}
// @Router /user [delete]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	//user.Name = c.Query("用户名")
	//user.Name = "ruan"
	//id,_ := strconv.Atoi(c.PostForm("id"))
	id, _ := strconv.Atoi(c.Query("id"))
	fmt.Println(id)
	//删除前应该校验一下参数，此处未处理
	models.DeleteUser(user, uint(id))
	c.JSON(http.StatusOK, gin.H{
		"message": "删除用户成功！",
		"id":      uint(id),
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param ori_username formData string false "原用户名"
// @param username formData string false "用户名"
// @param passwd formData string false "密码"
// @param Email formData string false "新邮箱"
// @Success 200 {string} json{"code","message"}
// @Router /user [put]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	utils.DB.First(&user, "name", c.PostForm("ori_username")) //根据主键获取匹配上的第一条数据（主键唯一）
	user.Name = c.PostForm("username")
	//user.Name = "陈忆"
	user.Passwd = c.PostForm("passwd")
	user.Email = c.PostForm("Email")
	//newValue := user.Passwd
	//newValue := c.Query("password")

	//newEmail := c.Query("email")
	//进行参数校验，是否合法
	ok, err := govalidator.ValidateStruct(user) //校验参数是否合法
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数不合法",
			"error":   err,
		})
		return
	}
	ret := govalidator.IsEmail(user.Email) //校验邮箱地址
	if !ret {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "-1",
			"message": "邮件地址错误",
			"data":    user.Email,
		})
		return
	}
	//此处应该将密码重新加密添加到结构体  最好线重新生成一个随机数

	//utils.DB.First(&user, "name", "chenyi")//根据主键获取匹配上的第一条数据（主键唯一）
	//models.UpdateUser2(user,newValue)
	models.UpdateUser2(user)
	c.JSON(http.StatusOK, gin.H{
		"message": "修改用户成功！",
		"user":    user,
	})
}

// 防止跨域站点伪请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	//打开一个websocket用于发送信息
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		fmt.Println("closed websocket......")
		err := ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	//订阅消息
	msg, err := utils.Subscribe(c, utils.PulishKey)
	if err != nil {
		fmt.Println(err)
	}
	tm := time.Now().Format("2006-01-02 15:04:05")
	m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
	fmt.Println(m)
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println(err)
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
