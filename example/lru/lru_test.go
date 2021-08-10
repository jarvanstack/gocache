package lru

import (
	"fmt"
	"testing"

	"github.com/dengjiawen8955/gocache/src/lru"
)

//LRU cache 基本使用.
func Test_function(t *testing.T) {
	k1 := "1"
	v1 := "10"
	k2 := "2"
	v2 := "20"
	//创建 lru 对象,并设置储存个数
	cache := lru.NewLRUCache(1)
	fmt.Println(cache.Size)
	//Set 储存
	cache.Set(k1, v1)
	//Set 储存覆盖
	cache.Set(k2, v2)
	//Get 值
	value1, _ := cache.Get(k1)
	fmt.Printf("value1=%#v\n", value1)
	value2, _ := cache.Get(k2)
	fmt.Printf("value2=%#v\n", value2)
}
