package tests

import (
	"fmt"
	"net"
	"testing"

	"github.com/dengjiawen8955/go_utils/erru"
	"github.com/dengjiawen8955/go_utils/logger"
	"github.com/dengjiawen8955/go_utils/netu"
	"github.com/dengjiawen8955/gocache/src/gocache"
)

var addr = ":8889"

//功能测试
func Test_netCtx_server(t *testing.T) {
	gc := gocache.NewGoCache()
	gc.Run(addr)
}
func Test_netCtx_client_set(t *testing.T) {
	conn, _ := net.Dial("tcp", addr)
	nc := netu.NewConnCtx(conn)
	nc.WriteData([]byte("set k1 v1"))
	b, err := nc.ReadData()
	if err != nil {
		logger.Error(erru.NewError("nc.ReadData", "read data from conn err"))
	}
	fmt.Printf("string(b): %v\n", string(b))
}

func Test_netCtx_client_get(t *testing.T) {
	conn, _ := net.Dial("tcp", addr)
	nc := netu.NewConnCtx(conn)
	nc.WriteData([]byte("get k1"))
	b, err := nc.ReadData()
	if err != nil {
		logger.Error(erru.NewError("nc.ReadData", "read data from conn err"))
	}
	fmt.Printf("string(b): %v\n", string(b))
}
