package tsdb

import (
	"fmt"
	"time"
)

func main() {
	timestamp := time.Now().Unix()
	fmt.Println(time.Unix(timestamp, 0))

}
