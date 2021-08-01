package BLC


// 网络服务常量管理
// 协议
const PROTOCAL = "tcp"
// 命令长度
const CMMAND_LENGTH = 12

// 命令分类
const(
	// 验证当前节点末端区块是否是最新
	CMD_VERSION = "version"
	// 从最长链上获取区块
	CMD_GETBLOCKS = "getblocks"
	// 向其他节点展示当前节点有哪些区块
	CMD_INV = "inv"
	// 请求指定区块
	CMD_GETDATA = "getdata"
	//接受到新区快后，进行处理
	CMD_BLOCK = "block"
)
