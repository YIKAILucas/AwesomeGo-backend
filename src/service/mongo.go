package service

import (
	"fmt"
	"log"
	_ "log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	url := "106.12.130.179:27017"
	mogsession, err := mgo.Dial(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	operater.mogSession = mogsession
	return nil
}

var url = "106.12.130.179:27017"
var db_name = "日志"
var table = "Warning"

func main() {
	session, err := mgo.Dial(url) //连接服务器
	if err != nil {
		panic(err)
	}

	c := session.DB(db_name).C(table)

	query := c.Find(bson.M{"id": 7})
	ps := []person{}
	query.All(&ps)

	for _, v := range ps {
		log.Println(v)

	}

	_ = c.Insert(map[string]interface{}{"id": 7, "name": "tongjh", "age": 25}) //增

	objid := bson.ObjectIdHex("55b97a2e16bc6197ad9cad59")

	_ = c.RemoveId(7) //删除

	_ = c.UpdateId(objid, map[string]interface{}{"id": 8, "name": "aaaaa", "age": 30}) //改
	var one map[string]interface{}
	_ = c.FindId(objid).One(&one) //查询符合条件的一行数据
	fmt.Println(one)

	var result []map[string]interface{}
	_ = c.Find(nil).All(&result) //查询全部
	fmt.Println(result)
}
