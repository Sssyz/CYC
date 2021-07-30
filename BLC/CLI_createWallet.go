package BLC

import (
	"fmt"
)

func(cli *CLI) createWallet(nodeID string){
	wallets := NewWallets(nodeID)

	wallets.CreateNewWallet(nodeID)

	fmt.Println(wallets.WalletsMap)
}