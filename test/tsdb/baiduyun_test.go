package tsdb

import (
	"fmt"
	"testing"
	"time"
)

func Test(test *testing.T) {
	timestamp := time.Now()

	x := timestamp.Format("2006-01-02T15:04:05Z")
	fmt.Println(x)

}
