package main

// gocache 主函数
import (
	"fmt"
	"testing"

	"github.com/dengjiawen8955/gocache/src/client"
)

//1. 链接所有的server用地址作为 map key, connCtx 作为 value
//2. 维护一个hash轮,添加 map 的 key
//3. 将 get 和 set 封装一下结合 1 + 2 组合为一个结构体 client
func Test_gocache_client(t *testing.T) {
	c := client.NewClient(addr1, addr2)
	sb, _ := c.Set("k1", []byte("v1"))
	fmt.Printf("%s\n", string(sb))
	gb, _ := c.Get("k1")
	fmt.Printf("%s\n", string(gb))
	gb, _ = c.Get("k2")
	fmt.Printf("k2=%s\n", string(gb))
	gb, _ = c.Get("k1")
	fmt.Printf("k1=%s\n", string(gb))
	sb, _ = c.Set("k2", []byte("v3"))
	fmt.Printf("%s\n", string(sb))
	gb, _ = c.Get("k2")
	fmt.Printf("k2=%s\n", string(gb))
	gb, _ = c.Get("k1")
	fmt.Printf("k1=%s\n", string(gb))
	// OK
	// v1
	// k2=v2
	// k1=v1
	// OK
	// k2=v3
	// k1=v1
}

//上面的测试全部命中 8888
//现在试试命中 8889
//ok
func Test_set_8889(t *testing.T) {
	c := client.NewClient(addr1, addr2)
	c.Set("k1", []byte("v1"))
	c.Set("ff", []byte("v1"))
	c.Set("fasd", []byte("v1"))
	c.Set("fsda", []byte("v1"))
	// node=":8888"
	// node=":8889"
	// node=":8889"
	// node=":8888"
	gb, _ := c.Get("ff")
	fmt.Printf("k1=%s\n", string(gb))

}
