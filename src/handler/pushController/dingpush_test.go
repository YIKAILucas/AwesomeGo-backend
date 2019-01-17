package pushController

import (
	"fmt"
	"testing"
)

func TestDingPush(test *testing.T) {
	//x:=""
	//fmt.Println(getToken(&x))
	//fmt.Println("token"+x)
	err := PubMessage("golang")
	if err != nil {
		fmt.Println(err)
	}
}
