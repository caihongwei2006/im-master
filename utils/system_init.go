package utils

import (
	"fmt"

	"log"
	"os"
	"time"

	"context"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string
	Email         string
	ClientIP      string
	ClientPost    string
	LoginTime     uint64
	HeartBeatTime uint64
	LogOutTime    uint64
	IsLogout      bool
	DeviceInfo    string
}

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("config file error!!!!!", err)
	}
	fmt.Println("config file success!!!!!!!")
	fmt.Println("config ", viper.Get("mysql"))
	fmt.Println("config mysql", viper.Get("mysql"))
}
func InitMysql() {
	//自定义日志部分，打印sql语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, //慢sql阈值
			LogLevel:      logger.Info, //日志级别
			Colorful:      true,        //是否开启彩色打印
		},
	)

	var err error
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println("\n\n\n, something is wrong with connecting to database", err)
	}
	user := UserBasic{}
	DB.Find(&user)
	fmt.Println("utils systeminit.go working")
}

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.db"),
		PoolSize:     viper.GetInt("redis.poolsize"),
		MinIdleConns: viper.GetInt("redis.minidleconns"),
	})
	pong, err := Redis.Ping(Redis.Context()).Result()
	if err != nil {
		fmt.Println("redis connect failed", err)
	} else {
		fmt.Println("redis connect success", pong)
	}
}

const (
	PublishKey = "websocket"
)

// Publish 将发布消息到redis
func Publish(ctx context.Context, channel string, message string) error {
	var err error
	err = Redis.Publish(ctx, channel, message).Err()
	if err != nil {
		fmt.Println("publish error", err)
	}
	return err
}
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Redis.Subscribe(ctx, channel)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println("Subscribe error", err)
	}
	fmt.Println("Subscribe...", sub)
	return msg.Payload, err
}
