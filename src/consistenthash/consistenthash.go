package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"

	"github.com/dengjiawen8955/go_utils/erru"
)

//一致性hash
//hash环 + 副本节点

//自定义hash函数
type Hash func(data []byte) uint32

//节点环实现
type NodeCircle struct {
	//hash值计算函数
	hashFunc Hash
	//节点的副本数量
	copiesNum int
	//排序后的节点的hash值, 数量为 节点数量 * 副本数量
	nodesHash []int
	//key:节点的hash值,value:节点名称
	nodesHashMap map[int]string
}

// 初始化 NodeCircle 节点环
//  copiesNum 为副本的数量, copiesNum > 1
//  hashFunc 用户自定义 hash 函数, 可以为 nil
func NewNodeCircle(copiesNum int, hashFunc Hash) *NodeCircle {
	n := &NodeCircle{
		hashFunc:     hashFunc,
		copiesNum:    copiesNum,
		nodesHashMap: make(map[int]string),
	}
	if n.hashFunc == nil {
		//默认 hash 函数
		n.hashFunc = crc32.ChecksumIEEE
	}
	return n
}

// 添加节点
//  keys 可以为节点的名称, 通过 keys 计算 hash 值
func (n *NodeCircle) AddNodes(nodes ...string) {
	for _, node := range nodes {
		//复制副本数量
		for i := 0; i < n.copiesNum; i++ {
			//计算每个副本的 hash 值
			h := int(n.hashFunc([]byte(strconv.Itoa(i) + node)))
			//将副本的 hash 值添加到 keys 中
			n.nodesHash = append(n.nodesHash, h)
			//将副本的 hash 存入map, value 为原体
			n.nodesHashMap[h] = node
		}
	}
	//调用快速排序库,将 keys 排序
	sort.Ints(n.nodesHash)
}

// 通过 key 获得 Node 节点
func (n *NodeCircle) GetNode(key string) (string, error) {
	var err error
	if len(n.nodesHash) == 0 {
		//如果没有节点,就报错
		err = erru.NewError("GetNode.NoNode", "len(n.nodesHash) no node in node circle")
		return "", err
	}
	//计算 key 的 hash 值
	h := int(n.hashFunc([]byte(key)))
	//使用二分法计算出 h 应该属于的节点的下标
	idx := sort.Search(len(n.nodesHash), func(i int) bool {
		return n.nodesHash[i] >= h
	})
	//取模是因为最后一节环的 hash 节点为第一个节点 idx = 0
	//TODO: 试一试入如果为 9 然后 -1 如何
	// idx := idx%len(n.nodesHash)
	if idx == len(n.nodesHash) {
		idx = 0
	}
	return n.nodesHashMap[n.nodesHash[idx]], nil
}
