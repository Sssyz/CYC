package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {

}
func printUsage(){
	fmt.Println("Usage:")
	fmt.Println("\t createwallet --创建钱包")
	fmt.Println("\t addresslist --输出所有钱包地址")
	fmt.Println("\t createblockchian -address address --余额")
	fmt.Println("\t send -from FROM -to TO -amount AMOUNT -转账")
	fmt.Println("\t printchain --输出区块信息")
	fmt.Println("\t getbalance -addresss --输出账号余额")
	fmt.Println("\t test -addresss --测试")
	fmt.Println("\t set_id -port PORT --设置端口节点号")
	fmt.Println("\t port --访问的节点号")
	fmt.Println("\t start --启动节点")
}
func isVaildArgs(){
	if len(os.Args)<2{
		printUsage()
		os.Exit(1)
	}
}
//func(cli *CLI)  addBlock(txs []*Transaction){
//	blockChain := GetBlockChainObject()
//	defer blockChain.DB.Close()
//	blockChain.AddBlockToBlockchain(txs)
//}





func(cli *CLI) Run(){
	nodeId := GetEnvNodeId()

	isVaildArgs()
	testCmd := flag.NewFlagSet("test",flag.ExitOnError)
	addresslistCmd :=flag.NewFlagSet("addresslist",flag.ExitOnError)
	startCmd := flag.NewFlagSet("start",flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet",flag.ExitOnError)

	sendBlockCmd := flag.NewFlagSet("send",flag.ExitOnError)
	flagFrom := sendBlockCmd.String("from","","转账源地址")
	flagTo := sendBlockCmd.String("to","","转账目的地址")
	flagAmount := sendBlockCmd.String("amount","","转账金额")

	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	//节点号设置
	setNodeIdCmd := flag.NewFlagSet("set_id",flag.ExitOnError)
	//端口号参数
	flagPortArg := setNodeIdCmd.String("port","","设置节点ID")
	getBalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)
	getBalanceWithAddress :=getBalanceCmd.String("address","","查询账号")

	createBlockChainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	flagCreateBlockChainWithAddress :=createBlockChainCmd.String("address","","创建创世区块的地址")

	switch os.Args[1] {
	case "start":
		err := startCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "set_id":
		err := setNodeIdCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "test":
		err := testCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "addresslist":
		err := addresslistCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}
	//节点启动服务
	if startCmd.Parsed(){
		cli.startNode(nodeId)
	}

	if setNodeIdCmd.Parsed(){
		if *flagPortArg==""{
			fmt.Println("请输入端口...")
			os.Exit(1)
		}
		cli.SetNodeId(*flagPortArg)
	}

	if sendBlockCmd.Parsed(){
		if *flagFrom ==""||*flagTo ==""||*flagAmount ==""{
			printUsage()
			os.Exit(1)
		}
		//fmt.Println(*flagAddBlockData)
		//cli.addBlock([]*Transaction{})
		//fmt.Println(*flagFrom)
		//fmt.Println(*flagTo)
		//fmt.Println(*flagAmount)
		//fmt.Println(JSONToArray(*flagFrom))
		//fmt.Println(JSONToArray(*flagTo))
		//fmt.Println(JSONToArray(*flagAmount))
		from :=JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		for index,fromAddress := range from{
			isValidFrom := IsValidForAddress([]byte(fromAddress))
			isValidTo := IsValidForAddress([]byte(to[index]))
			if isValidFrom==false{
				fmt.Printf("From地址%s无效!\n",fromAddress)
			}else if isValidTo == false{
				fmt.Printf("To地址%s无效!\n",to[index])
			}
		}
		amount := JSONToArray(*flagAmount)
		//三个数组长度需要相等
		if len(from)!=len(to)||len(from)!=len(amount)||len(to)!=len(amount){
			println("the numbers of from,to and amount must be euqal!")
			printUsage()

		}
		cli.send(from,to,amount,nodeId)
	}
	if printChainCmd.Parsed(){

		//fmt.Println("输出所有区块信息")
		cli.printChain(nodeId)
	}
	if createBlockChainCmd.Parsed(){
		if IsValidForAddress([]byte(*flagCreateBlockChainWithAddress))==false{
			fmt.Println("地址%s无效...",*flagCreateBlockChainWithAddress)
			printUsage()
			os.Exit(1)
		}
		//fmt.Println("创建创世区块")
		cli.createGenesisBlockChain(*flagCreateBlockChainWithAddress,nodeId)
	}
	if getBalanceCmd.Parsed(){
		if IsValidForAddress([]byte(*getBalanceWithAddress))==false{
			fmt.Printf("地址%s无效...",*getBalanceWithAddress)
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceWithAddress,nodeId)
	}
	// 创建钱包
	if createWalletCmd.Parsed(){

		cli.createWallet(nodeId)
	}
	if addresslistCmd.Parsed(){

		cli.addressList(nodeId)
	}
	if testCmd.Parsed(){
		cli.TestMethod(nodeId)
	}

}
