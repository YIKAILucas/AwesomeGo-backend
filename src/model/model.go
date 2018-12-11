package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

//定义结构体(xorm支持双向映射)
type User struct {
	User_id    int64  `xorm:"pk autoincr"` //指定主键并自增
	Name       string `xorm:"unique"`      //唯一的
	Balance    string
	Time       int64 `xorm:"updated"` //修改后自动更新时间
	Creat_time int64 `xorm:"created"` //创建时间
	//Version    int     `xorm:"version"` //乐观锁
}

//定义orm引擎
var X *xorm.Engine

//创建orm引擎
func init() {
	var err error
	X, err = xorm.NewEngine("mysql", "root:123@tcp(127.0.0.1:3306)/acke?charset=utf8")
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	if err := X.Sync(new(User)); err != nil {
		log.Fatal("数据表同步失败:", err)
	}
}

//增
func Insert(name string, balance string) (int64, bool) {
	user := new(User)
	user.Name = name
	user.Balance = balance
	affected, err := X.Insert(user)
	if err != nil {
		return affected, false
	}
	return affected, true
}

//删
func Del(id int64) {
	user := new(User)
	X.Id(id).Delete(user)
}

//改
func update(id int64, user *User) bool {
	affected, err := X.ID(id).Update(user)
	if err != nil {
		log.Fatal("错误:", err)
	}
	if affected == 0 {
		return false
	}
	return true
}

//查
func getinfo(id int64) *User {
	user := &User{User_id: id}
	is, _ := X.Get(user)
	if !is {
		log.Fatal("搜索结果不存在!")
	}
	return user
}
