package BLC
func(cli *CLI) printChain(){
	blockChain := GetBlockChainObject()
	defer blockChain.DB.Close()
	blockChain.PrintChain()
}