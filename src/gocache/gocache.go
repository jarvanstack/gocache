package gocache

import (
	"fmt"
	"net"
	"strings"

	"github.com/dengjiawen8955/go_utils/erru"
	"github.com/dengjiawen8955/go_utils/netu"
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
		go g.handler(c)
	}

}
func (g *GoCache) handler(conn net.Conn) {
	defer func() {
		e := recover()
		if e != nil {
			erru.Trace(e)
			fmt.Printf("err=%#v\n", e)
		}
	}()
	cc := netu.NewConnCtx(conn)
	//conn可以循环利用
	for {
		data, _ := cc.ReadData()
		//2.切割(先用这种笨办法)
		splits := strings.Split(string(data), " ")
		if len(splits) >= 2 {
			//不区分大小写比较
			if strings.EqualFold("get", splits[0]) {
				//GET
				d, ok := g.Map[splits[1]]
				if !ok {
					d = []byte("NOT EXISTS " + splits[1])
				}
				cc.WriteData(d)
			} else if strings.EqualFold("set", splits[0]) {
				//SET
				var d []byte
				if len(splits) >= 3 {
					g.Map[splits[1]] = []byte(splits[2])
					d = []byte("OK")
				} else {
					d = []byte("NOT ENOUGH PARAMITERS")
				}
				cc.WriteData(d)
			}
		} else {
			panic(erru.NewError("strings.Split", "split err while split request cmd"))
		}
	}

}
