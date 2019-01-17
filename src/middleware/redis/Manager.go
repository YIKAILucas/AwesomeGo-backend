package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", "106.12.130.179:6379")
	if err != nil {
		fmt.Println("conn redis failed, err:", err)
		return
	}
	defer c.Close()
}
