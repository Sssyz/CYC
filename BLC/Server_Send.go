package BLC

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

// 请求发送文件

// 发送请求
func sendMessage(to string,msg []byte ){
	fmt.Println("向服务器发送请求...")
	//1. 连接上服务器
	conn,err := net.Dial(PROTOCAL,to)
	if err != nil{
		log.Panicf("connect to server[%s] failed! %v\n",to,err)
	}
	defer conn.Close()
	//要发送的数据
	_,err = io.Copy(conn,bytes.NewReader(msg))
	if err!=nil{
		log.Panicf("add the data to conn failed! %v\n",err)
	}
}

// 区块链版本验证
func sendVersion(toAddress string,bc *BlockChain){
	// 1.获取当前节点区块高度
	height := bc.GetHeight()
	// 2.组装生成version
	versionData := Version{int(height),nodeAddress}
	fmt.Println("当前地址：",nodeAddress)
	// 3.数据的序列化
	data := gobEncode(versionData)
	// 4.将命令与版本组装成完整的请求
	request := append(commandToByte(CMD_VERSION),data...)
	fmt.Println("请求地址：",toAddress)
	// 5. 发送请求
	sendMessage(toAddress,request)
}

// 从指定节点同步数据
func sendGetBlocks(toAddress string){
	//1.生成数据
	data := gobEncode(GetBlocks{nodeAddress})

	//2.组装请求
	request := append(commandToByte(CMD_GETBLOCKS),data...)
	//3.发送请求
	sendMessage(toAddress,request)
}

// 发送获取指定区块请求
func sendGetData(toAddress string,hash []byte){
	//1.生成数据
	data := gobEncode(GetData{nodeAddress,hash})

	//2.组装请求
	request := append(commandToByte(CMD_GETDATA),data...)
	//3.发送请求
	sendMessage(toAddress,request)
}

// 向其他节点展示
func sendInv(toAddress string,hashed [][]byte){
	//1.生成数据
	data := gobEncode(Inv{nodeAddress,hashed})

	//2.组装请求
	request := append(commandToByte(CMD_INV),data...)
	//3.发送请求
	sendMessage(toAddress,request)
}

// 发送区块信息
func sendBlock(toAddress string,block []byte){
	//1.生成数据
	data := gobEncode(BlockData{nodeAddress,block})

	//2.组装请求
	request := append(commandToByte(CMD_BLOCK),data...)
	//3.发送请求
	sendMessage(toAddress,request)
}