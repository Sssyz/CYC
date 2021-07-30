package BLC

import "fmt"

// 打印所有钱包地址
func(cli *CLI) addressList(nodeID string)  {

	fmt.Println("打印所有钱包地址:")
	wallets :=  NewWallets(nodeID)
	for address,_:= range wallets.WalletsMap{
		fmt.Println(address)
	}
}