package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)
const walletFile = "Wallets_%s.dat"
type Wallets struct {
	WalletsMap map[string]*Wallet
}

// 创建钱包集合
func NewWallets(nodeID string) *Wallets{
	wallets := &Wallets{}
	wallets.WalletsMap = make(map[string]*Wallet)
	wallets.LoadFromFile(nodeID)
	return wallets
}

// 创建一个新钱包
func (w *Wallets) CreateNewWallet(nodeID string){
	wallet := NewWallet()
	fmt.Printf("Address:%s\n",wallet.GetAddress())
	w.WalletsMap[string(wallet.GetAddress())]=wallet
	// 把所有数据保存起来
	w.SavetoFile(nodeID)
}

//根据地址获得钱包对象
func(ws *Wallets) GetWallet(address string) *Wallet{
	return ws.WalletsMap[address]
}

//加载钱包文件
func(ws *Wallets) LoadFromFile(nodeId string) error{
	walletFile := fmt.Sprintf(walletFile,nodeId)
	if _,err := os.Stat(walletFile);os.IsNotExist(err){
		return err
	}

	fileContent,err := ioutil.ReadFile(walletFile)
	if err!=nil{
		log.Panic(err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder :=gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err!=nil{
		log.Panic(err)
	}
	ws.WalletsMap = wallets.WalletsMap
	return nil
}
//保存钱包，写入文件
func(ws *Wallets) SavetoFile(nodeId string){
	var content bytes.Buffer

	//注册 为了可以序列化任何类型
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err:=encoder.Encode(ws)
	if err != nil{
		log.Panic(err)
	}
	// 覆盖源文件  无法遍历文件中的键值对，无法用键值对存储
	walletFile := fmt.Sprintf(walletFile,nodeId)
	err = ioutil.WriteFile(walletFile,content.Bytes(),0644)
	if err!=nil{
		log.Panic(err)
	}
}