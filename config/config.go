package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitConfig() {
	viper.SetConfigName("config")   // 配置文件的文件名
	viper.SetConfigType("yaml")     // 配置文件的后缀
	viper.AddConfigPath("./config") // 获取到配置文件的路径
	if err := viper.ReadInConfig(); err != nil {
		panic("读取配置失败：" + err.Error())
	}
}

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
}

func NewDB() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到标准输出
		logger.Config{
			SlowThreshold: time.Second, // 慢查询阈值，超过此时间的查询将被记录
			LogLevel:      logger.Info, // 记录信息级别（Info、Warn、Error）
			Colorful:      true,        // 输出带颜色
		},
	)

	username := viper.GetString("mysql.username")
	password := viper.GetString("mysql.password")
	host := viper.GetString("mysql.host")
	port := viper.GetInt("mysql.port")
	database := viper.GetString("mysql.database")

	address := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)

	db, err := gorm.Open(mysql.Open(address), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	// 获取数据库连接实例并设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(100) // 设置最大打开连接数
	sqlDB.SetMaxIdleConns(10)  // 设置最大空闲连接数

	return db, nil
}
