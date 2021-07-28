package BLC
//创建创世区块
func(cli *CLI) createGenesisBlockChain(address string){


	blockchain:=CreateBlockchainWithGenesisBlock(address)
	defer blockchain.DB.Close()
	utxoSet := &UTXOSet{blockchain }
	utxoSet.ResetUTXOSet()

}
