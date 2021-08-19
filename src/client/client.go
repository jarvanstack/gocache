package client

import (
	"fmt"
	"net"

	"github.com/dengjiawen8955/go_utils/erru"
	"github.com/dengjiawen8955/go_utils/netu"
	"github.com/dengjiawen8955/gocache/pb/pbgc"
	"github.com/dengjiawen8955/gocache/src/consistenthash"
	"google.golang.org/protobuf/proto"
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
func (c *Client) Set(key string, value []byte) (*pbgc.GoCacheResponse, error) {
	nodeKey, err := c.Circle.GetNode(key)
	// fmt.Printf("node=%#v\n", nodeKey)
	if err != nil {
		return nil, erru.NewError("Circle.GetNode", "hash circle GetNode() not exist")
	}
	cc, ok := c.Map[nodeKey]
	if !ok {
		return nil, erru.NewError("Map[nodeKey]", "Map[nodeKey] not ok")
	}
	return pbReq(cc, "set", key, value)
}

//Get 方法
// 1.找到哪个链接
// 2.执行Get返回值
func (c *Client) Get(key string) (*pbgc.GoCacheResponse, error) {
	nodeKey, err := c.Circle.GetNode(key)
	if err != nil {
		return nil, erru.NewError("Circle.GetNode", "hash circle GetNode() not exist")
	}
	cc, ok := c.Map[nodeKey]
	if !ok {
		return nil, erru.NewError("Map[nodeKey]", "Map[nodeKey] not ok")
	}
	return pbReq(cc, "get", key, nil)
}

//使用 proto buf 通信, 发送 gocache 请求. 并返回响应.
func pbReq(cc *netu.ConnCtx, cmd, key string, value []byte) (*pbgc.GoCacheResponse, error) {
	var err error
	resp := &pbgc.GoCacheResponse{}
	req := &pbgc.GoCacheRequest{
		Cmd:   cmd,
		Key:   key,
		Value: value,
	}
	bs, err := proto.Marshal(req)
	if err != nil {
		return nil, erru.NewError("proto.Marshal", "proto.Marshal(req)")
	}
	cc.WriteData(bs)
	respB, err := cc.ReadData()
	if err != nil {
		return nil, erru.NewError("ConnCtx.ReadData", "cc.ReadData while get key")
	}
	err = proto.Unmarshal(respB, resp)
	if err != nil {
		return nil, erru.NewError("proto.Unmarshal", "roto.Unmarshal get key resp")
	}
	return resp, nil
}
