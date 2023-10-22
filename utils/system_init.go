package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB         *gorm.DB
	Red        *redis.Client
	TimeFormat = "2006-01-02 15:03:04"
)

func InitConfig() {
	viper.SetConfigFile("/home/chenyi/workspace/golangProject/GIN_IMchat/config/app.yml") //必须写绝对路径
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("InitConfig:", err)
	}
	fmt.Println("config app inited")
}

func InitMySQL() {
	//自定义日志模板  打印SQL语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢SQL的阈值
			LogLevel:      logger.Info, // 日志级别
			Colorful:      true,        //彩色
		},
	)

	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")),
		&gorm.Config{Logger: newLogger})
	/*	if err != nil {
		panic("failed to connect database")
	}*/
	fmt.Println("config mysql:", viper.GetString("mysql.dns"))
	//user := models.UserBasic{}
	//DB.Find(&user)
	//fmt.Println(user)

}

func InitRedis() {

	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	pong, err := Red.Ping(context.TODO()).Result()
	if err != nil {
		fmt.Println("init redis err:", err)
		os.Exit(-1)
	}
	fmt.Println("init redis success,pong:", pong)
}

const (
	PulishKey = "websocket"
)

// Publish 发布消息到Redis
func Publish(ctx context.Context, channel string, msg string) error {
	fmt.Println("publish......")
	err := Red.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 订阅Redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {

	sub := Red.Subscribe(ctx, channel)
	//defer sub.Close()
	//for msg := range sub.Channel(){
	//	fmt.Printf("channel=%s  message=%s\n",msg.Channel, msg.Payload)
	//}
	fmt.Println("Subscribe......", ctx)
	msg, err := sub.ReceiveMessage(ctx)
	fmt.Println("after Subscribe......")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("channel=%s  message=%s\n", msg.Channel, msg.Payload)
	return msg.Payload, err
}
