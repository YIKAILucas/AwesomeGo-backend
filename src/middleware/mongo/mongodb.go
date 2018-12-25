package mongo

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

var globalSession *mgo.Session

type MongodbInfo struct {
	DbHost    string
	AuthDb    string
	AuthUser  string
	AuthPass  string
	PoolLimit int
	Timeout   time.Duration
}

func InfoInit(info *MongodbInfo) {
	info.DbHost = "mongodb://root@root:127.0.0.1:27017"
	info.AuthDb = "admin"
	info.AuthUser = "root"
	info.AuthPass = "root"
	info.PoolLimit = 100
	info.Timeout = 60 * time.Second
}

func init() {
	s, err := mgo.Dial("106.12.130.179:27017")
	if err != nil {
		log.Fatalf("Create Session: %s\n", err)
	}
	globalSession = s

}

/*
权限型MongoDB
*/
//func init() {
//	info := &MongodbInfo{}
//	InfoInit(info)
//	dialInfo := &mgo.DialInfo{
//		Addrs:     []string{dbhost}, // 数据库地址 Dbhost: mongodb://user@123456:106.12.130.179:27017
//		Source:    info.AuthDb,      // 设置权限的数据库 authdb: admin
//		Username:  info.AuthUser,    // 设置的用户名 authuser: user
//		Password:  info.AuthPass,    // 设置的密码 authpass: 123456
//		PoolLimit: info.PoolLimit,   // 连接池的数量 poollimit: 100
//		Timeout:   info.Timeout,     // 连接超时时间 timeout: 60 * time.Second
//	}
//
//	s, err := mgo.DialWithInfo(dialInfo)
//	if err != nil {
//		log.Fatalf("Create Session: %s\n", err)
//	}
//	globalSession = s
//}

func connect(db, collection string) (*mgo.Session, *mgo.Collection) {
	ms := globalSession.Copy()
	c := ms.DB(db).C(collection)
	ms.SetMode(mgo.Monotonic, true)
	return ms, c
}

func Insert(db, collection string, doc interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Insert(doc)
}
func FindOne(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).One(result)
}

func FindAll(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).All(result)
}
func Update(db, collection string, selector, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Update(selector, update)
}

//更新，如果不存在就插入一个新的数据 `upsert:true`
func Upsert(db, collection string, selector, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	_, err := c.Upsert(selector, update)
	return err
}

// `multi:true`
func UpdateAll(db, collection string, selector, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	_, err := c.UpdateAll(selector, update)
	return err
}
func Remove(db, collection string, selector interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Remove(selector)
}

func RemoveAll(db, collection string, selector interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	_, err := c.RemoveAll(selector)
	return err
}
func FindPage(db, collection string, page, limit int, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).Skip(page * limit).Limit(limit).All(result)
}
func IsEmpty(db, collection string) bool {
	ms, c := connect(db, collection)
	defer ms.Close()
	count, err := c.Count()
	if err != nil {
		log.Fatal(err)
	}
	return count == 0
}

func Count(db, collection string, query interface{}) (int, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Count()
}

func Test2() {
	err := Insert("db", "123", map[string]interface{}{"id": 7, "name": "tongjh", "age": 25})
	if err != nil {

	}
}

type Data struct {
	Title   string
	Des     string
	Content string
	Img     string
	Date    time.Time
}

func Test() {
	data := &Data{
		//Id:      bson.NewObjectId().Hex(),
		Title:   "标题",
		Des:     "博客描述信息",
		Content: "博客的内容信息",
		Img:     "https://upload-images.jianshu.io/upload_images/8679037-67456031925afca6.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/700",
		Date:    time.Now(),
	}

	_ = Insert("acke", "test", data)
}
