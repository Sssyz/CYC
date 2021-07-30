package BLC
//创建创世区块
func(cli *CLI) createGenesisBlockChain(address string,nodeId string){


	blockchain:=CreateBlockchainWithGenesisBlock(address,nodeId)
	defer blockchain.DB.Close()
	utxoSet := &UTXOSet{blockchain }
	utxoSet.ResetUTXOSet()

}
