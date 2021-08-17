package main

// gocache 主函数
import (
	"github.com/dengjiawen8955/gocache/src/gocache"
)

var addr1 = ":8888"
var addr2 = ":8889"

func main() {
	go func() {
		gc := gocache.NewGoCache()
		gc.Run(addr1)
	}()
	gc := gocache.NewGoCache()
	gc.Run(addr2)
}
