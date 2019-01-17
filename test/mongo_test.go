package main

import (
	"awesomeProject/src/middleware/mongo"
	"fmt"
	"reflect"
)

func main() {
	err := mongo.Insert("acke", "test", map[string]interface{}{"id": 7, "name": "tongjh", "age": 25})
	mongo.Test()
	if err != nil {
		fmt.Println("错误")
	}
	b := []string{}
	fmt.Println(reflect.TypeOf(b), reflect.ValueOf(b).Kind())

}
