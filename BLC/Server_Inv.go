package BLC

type Inv struct {
	AddrFrom string // 当前节点
	Hashes [][]byte // 当前节点所拥有区块的Hash列表
}
