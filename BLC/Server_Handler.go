package BLC

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

//请求处理文件

// version
func handleVersion(request []byte,bc *BlockChain){
	fmt.Println("the request of version handle...")
	var buff bytes.Buffer
	var data Version
	//1.解析请求
	dataBytes := request[12:]
	//2.生成version结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)

	if err := decoder.Decode(&data);nil!=err{
		log.Panicf("deccod the version struct failed! %v\n",err)
	}
	//3.获取请求放的区块高度
	versionHeight := data.Height
	//4.获取自身节点的区块高度
	height := bc.GetHeight()
	fmt.Println("高度：",height,"高度：",versionHeight)
	//如果当前节点的区块高度大于versionHeight
	//将当前节点版本信息发送给请求节点
	if height > int64(versionHeight){
		sendVersion(data.AddrFrom,bc)
	}else if height<int64(versionHeight){
		//如果当前节点的区块高度大于versionHeight
		//像发送方发送数据同步请求
		sendGetBlocks(data.AddrFrom)
	}

}

// getblocks
// 数据同步请求处理
func handleGetBlocks(request []byte,bc *BlockChain){
	fmt.Println("the request of get blocks handle...")
	var buff bytes.Buffer
	var data GetBlocks
	//1.解析请求
	dataBytes := request[12:]
	//2.生成GetBlocks结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)

	if err := decoder.Decode(&data);nil!=err{
		log.Panicf("deccod the getBlocks struct failed! %v\n",err)
	}
	//3. 获取区块链所有的区块哈希
	hashes := bc.GetBlockHashes()
	fmt.Println(hashes)
	sendInv(data.AddrFrom,hashes)
}

// getData
// 处理获取指定区块的请求
func handleGetData(request []byte,bc *BlockChain){
	fmt.Println("the request of GetData handle...")
	var buff bytes.Buffer
	var data GetData
	//1.解析请求
	dataBytes := request[12:]
	//2.生成GetData结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)

	if err := decoder.Decode(&data);nil!=err{
		log.Panicf("deccod the GetData struct failed! %v\n",err)
	}
	//3.通过传过来的区块hash,获取本地节点的区块
	blockByte := bc.GetBlock(data.ID)
	sendBlock(data.AddrFrom,blockByte)
}
//Block
//接收到新区块后处理
func handleBlock(request []byte,bc *BlockChain){
	fmt.Println("the request of block handle...")
	var buff bytes.Buffer
	var data BlockData
	//1.解析请求
	dataBytes := request[12:]
	//2.生成blockdata结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)

	if err := decoder.Decode(&data);nil!=err{
		log.Panicf("deccod the blockdata struct failed! %v\n",err)
	}
	//3.将接收到的区块添加到区块链中
	blockByte := data.Block
	block := DeSerialize(blockByte)
	bc.AddBlock(block)
	//4.更新utxo table
	uTXOSet:=UTXOSet{bc}
	uTXOSet.Update()

}
// Inv
func handleInv(request []byte,bc *BlockChain){
	fmt.Println("the request of inv handle...")
	var buff bytes.Buffer
	var data Inv
	//1.解析请求
	dataBytes := request[12:]
	//2.生成Inv结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)

	if err := decoder.Decode(&data);nil!=err{
		log.Panicf("deccod the inv struct failed! %v\n",err)
	}
	for _,hash := range data.Hashes{
		sendGetData(data.AddrFrom,hash)
	}
}