package model

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"time"
)

type User struct {
	gorm.Model
	Birthday time.Time
	Age      int
	Name     string `gorm:"size:255"`       // string默认长度为255, 使用这种tag重设。
	Num      int    `gorm:"AUTO_INCREMENT"` // 自增

	IgnoreMe int `gorm:"-"` // 忽略这个字段
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func openDB(username, password, addr, name string) *gorm.DB {
	config := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		true,
		//"Asia/Shanghai"),
		"Local")

	db, err := gorm.Open("mysql", config)
	if err != nil {
		log.Errorf(err, "Database connection failed. Database name: %s", name)
		defer db.Close()
	}

	// set for db connection
	setupDB(db)

	return db
}
//"root:root@(106.12.130.179:3307)/acke_test?
func setupDB(db *gorm.DB) {
	db.LogMode(viper.GetBool("db.log"))
	db.DB().SetMaxOpenConns(viper.GetInt("db.max_open_connection")) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(viper.GetInt("db.max_idle_connection")) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
}

func InitDockerDB() *gorm.DB {
	return openDB(viper.GetString("docker_db.username"),
		viper.GetString("docker_db.password"),
		viper.GetString("docker_db.addr"),
		viper.GetString("docker_db.name"))
}

func DBInit() {
	db := InitDockerDB()
	setupDB(db)
	// 自动迁移模式
	db.AutoMigrate(&User{}, &Product{})
	// 启用日志记录器
	db.LogMode(true)
	user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
	//
	db.Create(&user)
	if db.NewRecord(user) {
		fmt.Println("插入失败,主键为空")
	}

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
