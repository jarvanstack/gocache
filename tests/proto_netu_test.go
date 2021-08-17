package tests

import (
	"fmt"
	"net"
	"testing"

	"github.com/dengjiawen8955/go_utils/netu"
	"github.com/dengjiawen8955/gocache/pb/kv"
	"google.golang.org/protobuf/proto"
)

//proto buf + netu 通信
//proto.Marshal()
//proto.UnMarshal()
//上面的2个方法封装
//proto 太好用了!!!
var addr3 = ":8888"

func Test_pb_service(t *testing.T) {
	l, _ := net.Listen("tcp", addr3)
	for {
		c, _ := l.Accept()
		cc := netu.NewConnCtx(c)
		b, _ := cc.ReadData()
		req := &kv.GetUserRequest{}
		proto.Unmarshal(b, req)
		fmt.Printf("req.Key: %v\n", req.Key)
		resp := &kv.GetUserResponse{Value: req.Key + " back"}
		wb, _ := proto.Marshal(resp)
		cc.WriteData(wb)
	}
}
func Test_pb_client(t *testing.T) {
	c, _ := net.Dial("tcp", addr3)
	cc := netu.NewConnCtx(c)
	req := &kv.GetUserRequest{Key: "req"}
	wb, _ := proto.Marshal(req)
	cc.WriteData(wb)
	rb, _ := cc.ReadData()
	resp := &kv.GetUserResponse{}
	proto.Unmarshal(rb, resp)
	fmt.Printf("resp.Value: %v\n", resp.Value)
}
