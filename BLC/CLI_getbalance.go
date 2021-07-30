package BLC

import "fmt"

//查余额
func(cli *CLI) getBalance(address string,nodeId string){
	//fmt.Println("地址："+address)
	blockchian := GetBlockChainObject(nodeId)
	defer blockchian.DB.Close()
	utxoSet := &UTXOSet{blockchian}
	amount := utxoSet.GetBalance(address)
	fmt.Printf("%s一共有%d个Token",address,amount)
}