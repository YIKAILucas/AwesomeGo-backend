package hbase

import (
	"fmt"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"github.com/tsuna/gohbase/pb"
	"golang.org/x/net/context"
	"testing"
)

func Test(t *testing.T) {
	option := gohbase.EffectiveUser("main")
	client := gohbase.NewClient("106.12.130.179", option)

	v := map[string]map[string][]byte{
		"image": map[string][]byte{
			"rgb": []byte("test"),
			"ir":  []byte("test"),
		},
	}
	putRequest, err := hrpc.NewPutStr(context.Background(), "collection_img", "12", v)
	if err != nil {
		fmt.Println("err")
	}
	rsp, err := client.Put(putRequest)
	if err != nil {
		fmt.Println("err")

	}

	for _, cell := range rsp.Cells {
		fmt.Println(string((*pb.Cell)(cell).GetFamily()))
		fmt.Println(string((*pb.Cell)(cell).GetQualifier()))
		fmt.Println(string((*pb.Cell)(cell).GetValue()))
		fmt.Println((*pb.Cell)(cell).GetTimestamp())
	}
}
