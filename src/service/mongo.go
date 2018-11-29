package service

import (
	"fmt"
	_ "log"

	"gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"
)

type Operater struct {
	mogSession *mgo.Session
	dbname     string
	document   string
}

//集合的结构 其他中AGE NAME HEIGHT 名字的首字母必须大写，如果不无法访问到
type person struct {
	AGE    int
	NAME   string
	HEIGHT int
}

//连接数据库
func (operater *Operater) connect() error {
	url := "127.0.0.1"
	mogsession, err := mgo.Dial(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	operater.mogSession = mogsession
	return nil
}

func main() {
	connect()
}
