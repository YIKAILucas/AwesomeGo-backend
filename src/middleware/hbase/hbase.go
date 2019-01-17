package hbase

import (
	"fmt"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"github.com/tsuna/gohbase/pb"
	"golang.org/x/net/context"
)

func main() {

	client := gohbase.NewClient("10.18.218.17,10.18.218.5,10.18.218.9")

	getRequest, _ := hrpc.NewGetStr(context.Background(), "emp", "1")
	getRsp, _ := client.Get(getRequest)

	for _, cell := range getRsp.Cells {
		fmt.Println(string((*pb.Cell)(cell).GetFamily()))
		fmt.Println(string((*pb.Cell)(cell).GetQualifier()))
		fmt.Println(string((*pb.Cell)(cell).GetValue()))
		fmt.Println((*pb.Cell)(cell).GetTimestamp())
	}
}
