package gocache

import (
	"fmt"
	"net"
	"strings"

	"github.com/dengjiawen8955/go_utils/erru"
	"github.com/dengjiawen8955/go_utils/netu"
	"github.com/dengjiawen8955/gocache/pb/pbgc"
	"google.golang.org/protobuf/proto"
)

//服务端
type GoCache struct {
	//全局map
	Map map[string][]byte
}

func NewGoCache() *GoCache {
	return &GoCache{
		Map: make(map[string][]byte),
	}
}

//启动
func (g *GoCache) Run(addr string) {
	var err error
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(erru.NewError("net.Listen", addr))
	}
	fmt.Printf("%s\n", "go cache listen and service at "+addr)
	for {
		c, err := l.Accept()
		if err != nil {
			panic(erru.NewError("net.Accept", ""))
		}
		cc := netu.NewConnCtx(c)
		go func(cc *netu.ConnCtx) {
			for {
				g.handler(cc)
			}
		}(cc)
	}

}
func (g *GoCache) handler(cc *netu.ConnCtx) {
	req := &pbgc.GoCacheRequest{}
	resp := &pbgc.GoCacheResponse{}
	defer func() {
		//拿到错误, 返回错误.
		e := recover()
		if e != nil {
			erru.Trace(e)
			e, _ := e.(*erru.Err)
			fmt.Printf("err=%#v\n", e)
			resp.Code = e.Code
			resp.Value = []byte(e.Msg)
			b, err := proto.Marshal(resp)
			if err != nil {
				fmt.Printf("err=%#v\n", err)
			}
			cc.WriteData(b)
		}
	}()
	data, err := cc.ReadData()
	if err != nil {
		panic(erru.NewError("cc.ReadData", ""))
	}
	err = proto.Unmarshal(data, req)
	if err != nil {
		panic(erru.NewError("proto.Unmarshal", ""))
	}
	if strings.EqualFold(req.Cmd, "get") {
		//get
		b, ok := g.Map[req.Key]
		if !ok {
			panic(erru.NewError("KeyNotExists", "req key not exists "+req.Key))
		}
		resp.Code = "ok"
		resp.Value = b
		b2, err := proto.Marshal(resp)
		if err != nil {
			panic(erru.NewError("proto.Marshal", "proto.Marshal(resp)"))
		}
		cc.WriteData(b2)
	} else if strings.EqualFold(req.Cmd, "set") {
		//set
		g.Map[req.Key] = req.Value
		resp.Code = "ok"
		resp.Value = nil
		b2, err := proto.Marshal(resp)
		if err != nil {
			panic(erru.NewError("proto.Marshal", "proto.Marshal(resp)"))
		}
		cc.WriteData(b2)
	}
}
