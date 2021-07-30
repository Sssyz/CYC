package BLC

import "fmt"

func(cli *CLI) TestMethod(nodeId string){
	blockchain := GetBlockChainObject(nodeId)

	defer blockchain.DB.Close()

	utxoMap := blockchain.FindUTXOMap()
	fmt.Println(utxoMap)
	utxoSet := UTXOSet{blockchain}
	utxoSet.ResetUTXOSet()
	fmt.Println("重置余额数据库成功！")
}