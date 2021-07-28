package BLC

import "fmt"

// 打印所有钱包地址
func(cli *CLI) addressList()  {

	fmt.Println("打印所有钱包地址:")
	wallets :=  NewWallets()
	for address,_:= range wallets.WalletsMap{
		fmt.Println(address)
	}
}