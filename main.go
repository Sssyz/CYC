package main

import (
	"pulicChain/part67-UTXOSet-update/BLC"
)

func main(){

	//创世区块链
	//blockchain := BLC.CreateBlockchainWithGenesisBlock()
	// 创建失败,读取已有区块链
	//if blockchain==nil{
	//	blockchain = BLC.ReloadBlockChain()
	//}
	//defer blockchain.DB.Close()
	//fmt.Println(blockchain)
	cli := BLC.CLI{}

	cli.Run()


}