package tests

import (
	"fmt"
	"testing"

	"github.com/dengjiawen8955/go_utils/erru"
)

//err interface 类型转换为 自定义 err 类型.
func Test_err_trans(t *testing.T) {
	var err error
	err = erru.NewError("Circle.GetNode", "hash circle GetNode() not exist")
	e := err.(*erru.Err)
	fmt.Printf("e.Code: %v\n", e.Code)
	//e.Code: Circle.GetNode
}
