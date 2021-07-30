package BLC
func(cli *CLI) printChain(nodeID string){
	blockChain := GetBlockChainObject(nodeID)
	defer blockChain.DB.Close()
	blockChain.PrintChain()
}