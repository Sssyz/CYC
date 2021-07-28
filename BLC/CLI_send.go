package BLC
//转账
func(cli *CLI) send(from []string,to []string,amount []string){

	blockchain:= GetBlockChainObject()
	defer blockchain.DB.Close()
	//1.数字签名 2.验证是否合法
	blockchain.MineNewBlock(from,to,amount)

	utxsSet := &UTXOSet{blockchain}
	//转账成功以后，需要更新一下
	utxsSet.Update()
}
