package client

import (
	"bytes"
	"fmt"
	"net"

	"github.com/dengjiawen8955/go_utils/erru"
	"github.com/dengjiawen8955/go_utils/netu"
	"github.com/dengjiawen8955/gocache/src/consistenthash"
)

//GoCache封装调用
type Client struct {
	//1. 键是地址 value 是封装的连接
	Map map[string]*netu.ConnCtx
	//2. 一致性hash论
	Circle *consistenthash.NodeCircle
}

func NewClient(addrs ...string) *Client {
	m := make(map[string]*netu.ConnCtx, 0)
	circle := consistenthash.NewNodeCircle(3, nil)
	for _, addr := range addrs {
		c, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Printf("err=%#v\n", erru.NewError("net.Dial", "GoCache NewClient() Dial Error "+addr))
			continue
		}
		fmt.Printf("%s\n", "dial success "+addr)
		circle.AddNodes(addr)
		m[addr] = netu.NewConnCtx(c)
	}
	return &Client{
		Map:    m,
		Circle: circle,
	}
}

//Set 方法
// 1.找到哪个链接
// 2.执行Set返回值
// 3.读取set返回值
func (c *Client) Set(key string, value []byte) ([]byte, error) {
	nodeKey, err := c.Circle.GetNode(key)
	fmt.Printf("node=%#v\n", nodeKey)
	if err != nil {
		return nil, erru.NewError("Circle.GetNode", "hash circle GetNode() not exist")
	}
	cc, ok := c.Map[nodeKey]
	if !ok {
		return nil, erru.NewError("Map[nodeKey]", "Map[nodeKey] not ok")
	}
	msg := fmt.Sprintf("set %s ", key)
	if err != nil {
		return nil, erru.NewError("baseu.IntToBytes", "")
	}
	var b bytes.Buffer
	b.WriteString(msg)
	b.Write(value)
	cc.WriteData(b.Bytes())
	return cc.ReadData()
}

//Get 方法
// 1.找到哪个链接
// 2.执行Get返回值
func (c *Client) Get(key string) ([]byte, error) {
	nodeKey, err := c.Circle.GetNode(key)
	if err != nil {
		return nil, erru.NewError("Circle.GetNode", "hash circle GetNode() not exist")
	}
	cc, ok := c.Map[nodeKey]
	if !ok {
		return nil, erru.NewError("Map[nodeKey]", "Map[nodeKey] not ok")
	}
	msg := fmt.Sprintf("get %s", key)
	cc.WriteData([]byte(msg))
	return cc.ReadData()
}
