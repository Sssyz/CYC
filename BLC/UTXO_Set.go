package BLC

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)
// 1.方法，功能：遍历整个数据库，读取所有未花费的UTXO，然后将所有的UTXO存储到数据库
// reset
// 去遍历数据库时 {}   []*TXOUTputs
// [string]*TXOutputs
// txHash,TXOutputs := range txOutputsMap{
//}
const utxoTableName = "utxoTableName"

type UTXOSet struct {
	BlockChain *BlockChain

}

//重置数据库表
func(utxoSet *UTXOSet)ResetUTXOSet(){

	err :=utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx)error{
		b := tx.Bucket([]byte(utxoTableName))
		if b!=nil{
			err := tx.DeleteBucket([]byte(utxoTableName))
			if err!=nil{
				log.Panic(err)
			}
		}
		b,_=tx.CreateBucket([]byte(utxoTableName))
		if b!=nil{
			txOutputsMap := utxoSet.BlockChain.FindUTXOMap()
			for keyHash,outs := range txOutputsMap{
				txHash,_ := hex.DecodeString(keyHash)
				b.Put(txHash,outs.Serialize())
			}
		}

		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
}
//通过余额表找到对应的utxo
func(utxoSet *UTXOSet) findUTXOForAddress(address string)[]*UTXO{
	var utxos []*UTXO
	utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b!=nil{
			c := b.Cursor()
			for k,v := c.First();k!=nil;k,v = c.Next(){
				//fmt.Printf("key=%s,value=%s",k,v)
				txOutputs := DeSerializeTxOutputs(v)
				for _,utxo := range txOutputs.UTXOS{
					if utxo.Output.UnLockWithAddress(address){
						utxos = append(utxos,utxo)
					}
				}
			}
		}

		return nil
	})
	return utxos
}
func(utxoSet *UTXOSet) GetBalance(address string)int64{

	UTXOS := utxoSet.findUTXOForAddress(address)
	var amount int64
	for _,utxo := range UTXOS{
		amount += utxo.Output.Value
	}
	return amount
}

//返回凑多少钱和对应的txoutput的tx的hash和index
func(utxoSet *UTXOSet) FindUnPackageSpendableUTXOS(from string,txs []*Transaction)[]*UTXO{
	var unUTXOs []*UTXO
	spentTXOutputs := make(map[string][]int)
	//var money int64 = 0
	for _,tx := range txs {
		if tx.IsCoinBaseTransaction() == false {

			for _, in := range tx.Vins {
				//是否解锁
				pubKeyHash := Base58Decode([]byte(from))
				ripemd160Hash := pubKeyHash[1:len(pubKeyHash)-4]
				if in.UnLockRipemd160Hash(ripemd160Hash) {

					key := hex.EncodeToString(in.TXHash)

					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}
			}
		}

	}
	fmt.Println(spentTXOutputs)
	for _,tx:=range txs{
		//若当前的txHash都没有被记录消费
		spentArray,ok:=spentTXOutputs[hex.EncodeToString(tx.TxHash)]

		if ok==false{

			for index,out := range tx.Vouts{

				if out.UnLockWithAddress(from) {
					utxo := &UTXO{tx.TxHash,index,out}
					unUTXOs = append(unUTXOs, utxo)
				}
			}

		} else{
			//Vouts
			for index,out := range tx.Vouts{
				//判断是否花费
				flag := false
				//是否解锁
				if out.UnLockWithAddress(from){
					//判断是否被消费
					for _,spentIndex := range spentArray{
						if spentIndex==index{
							flag = true
							break
						}
					}
					//遍历所有已记录花费，该outPut未花费
					if flag == false{
						utxo := &UTXO{tx.TxHash,index,out}
						unUTXOs = append(unUTXOs,utxo)
					}

				}

			}
		}
	}
	return unUTXOs



}
func(utxoSet *UTXOSet) FindSpendableUTXOs(from string,amount int64,txs []*Transaction)(int64,map[string][]int){
	unPackageUTXOS := utxoSet.FindUnPackageSpendableUTXOS(from,txs)
	spendableUTXO := make(map[string][]int)
	var money int64 = 0
	//未打包是否满足
	for _,UTXO := range unPackageUTXOS{
		money += UTXO.Output.Value
		txhash := hex.EncodeToString(UTXO.TxHash)
		spendableUTXO[txhash] = append(spendableUTXO[txhash],UTXO.Index)
		if money >= amount{
			return money,spendableUTXO
		}
	}
	//若是未打包不满足，需要迭代余额数据库
	utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {


		b := tx.Bucket([]byte(utxoTableName))
		if b!=nil{
			c := b.Cursor()
			loop:
			for k,v := c.First();k!=nil;k,v = c.Next(){
				//fmt.Printf("key=%s,value=%s",k,v)
				txOutputs := DeSerializeTxOutputs(v)
				for _,utxo := range txOutputs.UTXOS{
					if utxo.Output.UnLockWithAddress(from){
						money += utxo.Output.Value
						txhash := hex.EncodeToString(utxo.TxHash)
						spendableUTXO[txhash] = append(spendableUTXO[txhash],utxo.Index)
						if money >= amount{
							break loop
						}
					}
				}
			}
		}
		return nil
	})
	if money<amount{
		log.Panic("余额不足！")
	}
	return money,spendableUTXO

}

// 更新
func(utxoSet *UTXOSet) Update(){

	//从数据库中取最新的BLock删除对应的txinput 添加新增的txoutput
	block := utxoSet.BlockChain.Iterator().Next()

	ins := []*TXInput{}
	outsMap := make(map[string]*TXOutputs)

	for _,tx := range block.Txs {
		//找到删除数据
		for _, in := range tx.Vins {
			if (in.Vout != -1) {
				ins = append(ins, in)
			}
		}
	}
	for _,tx := range block.Txs {
		utxos := []*UTXO{}
		//找到待添加数据
		for index,out := range tx.Vouts{
			isUsed := false
			for _,in := range ins{

				if in.Vout==index && bytes.Compare(tx.TxHash,in.TXHash)==0{
					isUsed = true
				}else{
					continue
				}
			}
			//没被使用过才可以添加
			if isUsed == false{
				utxo:=&UTXO{tx.TxHash,index,out}
				fmt.Printf("utxos:%x %d\n",utxo.TxHash,utxo.Index)
				utxos = append(utxos,utxo)
			}
		}
		if len(utxos)>0{
			txHash := hex.EncodeToString(tx.TxHash)

			outsMap[txHash] = &TXOutputs{utxos}
		}

	}



	err:=utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b!=nil{
			//删除
			for _,in := range ins{
				txOutputsBytes:=b.Get(in.TXHash)
				if len(txOutputsBytes)==0{
					continue
				}
				txOutputs := DeSerializeTxOutputs(txOutputsBytes)
				utxos := []*UTXO{}
				//通过重新修改value值来
				for _,utxo := range txOutputs.UTXOS{
					if in.Vout == utxo.Index{
						if bytes.Compare(utxo.Output.Ripemd160Hash,Ripemd160Hash(in.Pubkey))==0{
							continue
						}else{
							utxos = append(utxos,utxo)
						}
					}
				}
				b.Delete(in.TXHash)

				if len(utxos)>0{
					outputs := &TXOutputs{utxos}
					b.Put(in.TXHash,outputs.Serialize())
				}
			}
			//增加数据
			for keyHash,outputs:=range outsMap{
				keyHashByte,_ := hex.DecodeString(keyHash)
				b.Put(keyHashByte,outputs.Serialize())
			}

		}



		return nil
	})
	if err!=nil{
		log.Panic(err)
	}


}