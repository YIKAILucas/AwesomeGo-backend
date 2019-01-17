package mongoDao

import (
	"awesomeProject/src/middleware/mongo"
	"fmt"
)

func GetDAU() {
	ms, connect := mongo.Connect(mongo.DB_NAME, mongo.DAU)
	defer ms.Close()

	result := make([]mongo.DauJson, 0, 10)

	_ = connect.Find(nil).All(&result)

	fmt.Println(result)
}
