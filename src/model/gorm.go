package model

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	Birthday time.Time
	Age      int
	Name     string `gorm:"size:255"`       // string默认长度为255, 使用这种tag重设。
	Num      int    `gorm:"AUTO_INCREMENT"` // 自增

	IgnoreMe int `gorm:"-"` // 忽略这个字段
}

func DBInit() {
	config := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		viper.GetString("docker_db.username"),
		viper.GetString("docker_db.password"),
		viper.GetString("docker_db.addr"),
		viper.GetString("docker_db.name"),
		true,
		//"Asia/Shanghai")
		"Local")
	db, err := gorm.Open("mysql", config)
	DB = db
	if err != nil {
		log.Errorf(err, "Database connection failed. Database name: %s", viper.GetString("docker_db.name"))
		defer DB.Close()
	}
	DB.LogMode(viper.GetBool("db.log"))
	DB.DB().SetMaxOpenConns(viper.GetInt("db.max_open_connection")) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	DB.DB().SetMaxIdleConns(viper.GetInt("db.max_idle_connection")) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	// 自动迁移模式
	DB.AutoMigrate(&User{}, &Device{}, &DevicesLifeCycle{})
	// 启用日志记录器
	DB.LogMode(true)

	//tx := db.Begin()
	//// 注意，一旦你在一个事务中，使用tx作为数据库句柄
	//
	//if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
	//	tx.Rollback()
	//}
	//
	//if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
	//	tx.Rollback()
	//}
	//
	//tx.Commit()

}
