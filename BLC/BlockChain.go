package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

const dbName = "blockchain_%s.db"
const blockTableName = "blocks"
type BlockChain struct {
	Tip []byte //最新区块的hash
	DB *bolt.DB
}

func GetBlockChainObject(nodeID string) *BlockChain {
	if !dbExsits(nodeID){
		fmt.Println("创世区块不存在！")
		os.Exit(1)
	}
	var Tip []byte
	dbName := fmt.Sprintf(dbName,nodeID)
	fmt.Println(dbName)
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockTableName))
		if b!=nil{
			Tip = b.Get([]byte("l"))
		}
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
	//fmt.Println("读取最新区块链成功！")
	//读取区块链
	return &BlockChain{Tip,db}
}

//1. 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock(address string,nodeID string) *BlockChain{
	//判断是否存在数据库
	if(dbExsits(nodeID)){
		fmt.Println("创世区块已存在！")
		//在数据库中读取最新区块链
		os.Exit(1)

	}

	// 当数据库不存在时，创建创世区块链
	fmt.Println("正在创建创世区块。。。")
	dbName := fmt.Sprintf(dbName,nodeID)
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	//defer db.Close()
	var genesisBlock *Block
	err = db.Update(func(tx *bolt.Tx)error{

		b,err := tx.CreateBucket([]byte(blockTableName))
		if err!=nil{
			log.Panic(err)

		}

		if b!=nil{
			//创建一个coinbase transaction
			txCoinbase := NewCoinbaseTransaction(address)
			//创建创世区块
			genesisBlock = CreateGenesisBlock([]*Transaction{txCoinbase})
			//将创世区块放入数据库
			err=b.Put(genesisBlock.Hash,genesisBlock.Serialize())
			if err!=nil{
				log.Panic(err)
			}
			//存储最新的区块的hash
			err=b.Put([]byte("l"),genesisBlock.Hash)

			if err!=nil{
				log.Panic(err)
			}
		}

		return nil
	})
	fmt.Println("创建创世区块成功！")
	return &BlockChain{genesisBlock.Hash,db}
}

// 增加区块到区块链中
func(blc *BlockChain) AddBlockToBlockchain(txs []*Transaction){
	var height int64
	var preHash []byte
	//获取新增区块的height和preHash
	fmt.Println("开始挖矿。。。")
	err:=blc.DB.View(func(tx *bolt.Tx)error{
		b := tx.Bucket([]byte(blockTableName))
		if b!=nil{
			//blockHash := b.Get([]byte("l"))
			block := DeSerialize(b.Get(blc.Tip))
			height = block.Height+1
			preHash = block.Hash
		}
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
	// 创建新区块并添加数据库

	err = blc.DB.Update(func(tx *bolt.Tx) error{
		b := tx.Bucket([]byte(blockTableName))
		if b!=nil{
			newBlock := NewBlock(txs,height,preHash)
			newBlockByte := newBlock.Serialize()
			//添加区块信息值数据库
			err :=b.Put(newBlock.Hash,newBlockByte)
			if err!=nil{
				log.Panic(err)
			}

			//更新区块链的Tip以及数据库的l
			blc.Tip = newBlock.Hash
			b.Put([]byte("l"),newBlock.Hash)
			fmt.Println("挖矿成功！")
		}

		return nil
	})
	if err!=nil{
		log.Panic(err)
	}

}

// 遍历输出所有区块的信息
func (blc *BlockChain) PrintChain(){

	blockchainIterator := blc.Iterator()
	for {
		block := blockchainIterator.Next()
		fmt.Printf("Height: %d\n",block.Height)
		fmt.Printf("PrevBlockHash: %x\n",block.PrevBlockHash)

		fmt.Printf("TimeStamp: %s\n",time.Unix(block.Timestamp,0).Format("2006-01-02 15:04:05"))
		fmt.Printf("Hash: %x\n",block.Hash)
		fmt.Printf("Nonce: %d\n",block.Nonce)
		fmt.Println("Txs:")
		for _,tx := range block.Txs{
			fmt.Printf("%x\n",tx.TxHash)
			fmt.Println("Vins:")
			for _,in:=range tx.Vins{
				fmt.Printf("%x\n",in.TXHash)
				fmt.Printf("%d\n",in.Vout)
				//fmt.Printf("%x\n",in.Signature)
				fmt.Printf("%x\n",in.Pubkey)
			}
			fmt.Println("Vouts:")
			for _,out:=range tx.Vouts{
				fmt.Printf("%d\n",out.Value)
				fmt.Printf("%v\n",out.Ripemd160Hash)
			}
		}
		fmt.Println("----------------------------------------------------------------------------------")
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0))==0{
			break
		}

	}
}

// 判断数据库是否存在
func dbExsits(nodeID string) bool{
	//生成不同节点的数据库文件
	dbName := fmt.Sprintf(dbName,nodeID)
	if _,err := os.Stat(dbName);os.IsNotExist(err){
		return false
	}
	return true
}
//如果一个地址所对应的TXout未花费，那么这个就应该被添加到数组中
func (blockchain *BlockChain) UnUTXOs(address string,txs []*Transaction)[]*UTXO{

	var unUTXOs []*UTXO
	spentTXOutputs := make(map[string][]int)
	//
	//pubKeyHash := Base58Decode([]byte(address))
	//ripemd160Hash := pubKeyHash[1:len(pubKeyHash)-4]
	//fmt.Printf("转换后%v\n",ripemd160Hash)
	// 处理未加入数据库中的交易
	for _,tx := range txs {
		if tx.IsCoinBaseTransaction() == false {

			for _, in := range tx.Vins {
				//是否解锁
				pubKeyHash := Base58Decode([]byte(address))
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

				if out.UnLockWithAddress(address) {
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
				if out.UnLockWithAddress(address){
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



	//迭代数据库
	blockIterator := blockchain.Iterator()
	for{
		block:=blockIterator.Next()

		fmt.Println(block)
		for _,tx:=range block.Txs {

			//txHash
			//Vins
			if tx.IsCoinBaseTransaction() == false {
				pubKeyHash := Base58Decode([]byte(address))
				ripemd160Hash := pubKeyHash[1:len(pubKeyHash)-4]
				for _, in := range tx.Vins {
					//是否解锁
					if in.UnLockRipemd160Hash(ripemd160Hash) {

						key := hex.EncodeToString(in.TXHash)

						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}
				}
			}

		}
		fmt.Printf("%v\n",spentTXOutputs)
		for _,tx:=range block.Txs{
			//若当前的txHash都没有被记录消费
			spentArray,ok:=spentTXOutputs[hex.EncodeToString(tx.TxHash)]
			//fmt.Printf("ok is %s",ok)
			if ok==false{

				for index,out := range tx.Vouts{

					if out.UnLockWithAddress(address) {

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
					if out.UnLockWithAddress(address){
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



		//终止遍历条件
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0))==0{
			break
		}
	}

	return unUTXOs
}
//去查找可用的output,即将被消费
func(blockchain *BlockChain) FindSpendableUTXOs(from string,amount int,txs []*Transaction)(int64,map[string][]int){
	//1.获取所有UTXOs
	utxos:=blockchain.UnUTXOs(from,txs)
	spendAbleUTXO := make(map[string][]int)
	//2.遍历utxos
	var value int64
	for _,utxo := range utxos{
		value = value + utxo.Output.Value
		hash := hex.EncodeToString(utxo.TxHash)
		spendAbleUTXO[hash] = append(spendAbleUTXO[hash],utxo.Index)
		if value>=int64(amount){
			break
		}
	}
	if value < int64(amount){
		fmt.Printf("%s's fund isnt enough\n",from)
		os.Exit(1)
	}
	return value,spendAbleUTXO
}
//挖掘新的区块
func (blockchain *BlockChain)MineNewBlock(from []string,to []string,amount []string,nodeId string){
	//fmt.Println(from)
	//fmt.Println(to)
	//fmt.Println(amount)
	//1.通过相关算法建立交易数组
	//main.exe send -from "['liyuechun']" -to "['zhangqiang']" -amount "['2']"
	utxoSet := &UTXOSet{blockchain}
	var txs []*Transaction
	var block *Block

	//奖励from的第一个（先添加奖励余额，该余额可以被使用）
	tx := NewCoinbaseTransaction(from[0])
	txs = append(txs,tx)

	//处理所有的交易
	for index,_ := range from{
		amountint,_ := strconv.Atoi(amount[index])
		//可能有多比交易，之前的交易还未存储到数据库中，在新建交易时，需要考虑已有的未保存的交易，因此传入txs
		tx := NewSimpleTransaction(from[index],to[index],int64(amountint),utxoSet,txs,nodeId)
		txs = append(txs,tx)
	}




	//blockchain.AddBlockToBlockchain(txs)
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b :=tx.Bucket([]byte(blockTableName))
		if b!=nil{
			hash:=b.Get([]byte("l"))
			blockBytes := b.Get(hash)
			block = DeSerialize(blockBytes)
		}

		return nil
	})
	//在建立新区块之前，要进行数字签名的验证
	for _,tx:= range txs{
		if blockchain.VerifyTransaction(tx,txs)==false{
			log.Panic("签名失败！")
			//os.Exit(1)
		}
	}


	//2. 建立新的区块
	block =NewBlock(txs,block.Height+1,block.Hash)

	////3.存储到数据库
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockTableName))
		if b!=nil{
			b.Put(block.Hash,block.Serialize())
			b.Put([]byte("l"),block.Hash)
			blockchain.Tip = block.Hash
		}

		return nil
	})

}
func(blockchain *BlockChain)GetBalance(address string) int64{
	utxos := blockchain.UnUTXOs(address,[]*Transaction{})
	fmt.Println(utxos)
	var amount int64
	for _,utxo := range utxos{
		amount += utxo.Output.Value
	}
	return amount



}
// 数字签名
func(bc *BlockChain)SignTransaction(tx *Transaction,private ecdsa.PrivateKey,txs []*Transaction){
	if tx.IsCoinBaseTransaction(){
		return
	}
	prevTxs := make(map[string]Transaction)

	for _,vin := range tx.Vins{
		prevTx,err := bc.FindTransaction(vin.TXHash,txs)
		if err!=nil{
			log.Panic(err)
		}
		prevTxs [hex.EncodeToString(prevTx.TxHash)] = prevTx
	}

	tx.Sign(private,prevTxs)
}
//找签名相关交易
func(bc *BlockChain) FindTransaction(txHash []byte,txs []*Transaction)(Transaction,error){
	for _,tx:=range txs{
		if bytes.Compare(tx.TxHash,txHash)==0{
			return *tx,nil
		}
	}

	bci := bc.Iterator()
	var hashInt big.Int
	for{
		block := bci.Next()
		for _,tx:=range block.Txs{
			if bytes.Compare(tx.TxHash,txHash)==0{
				return *tx,nil
			}
		}

		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0))==0{
			break
		}
	}
	return Transaction{},errors.New("Transaction is not found")
}
func(bc *BlockChain) VerifyTransaction(tx *Transaction,txs []*Transaction) bool{
	if tx.IsCoinBaseTransaction(){
		return true
	}
	prevTxs := make(map[string]Transaction)
	for _,vin := range tx.Vins{
		fmt.Println("1111")
		prevTx,err := bc.FindTransaction(vin.TXHash,txs)
		if err!=nil{
			log.Panic(err)
		}
		prevTxs [hex.EncodeToString(prevTx.TxHash)] = prevTx
	}
	return tx.Verify(prevTxs)
}

//[string]*TXoutputs
func(blc *BlockChain)FindUTXOMap()map[string]*TXOutputs{
	blcIterator := blc.Iterator()
	//存储已花费的utxo信息
	spentUTXOsMap := make(map[string][]*TXInput)
	utxoMap := make(map[string]*TXOutputs)
	//1.spentUTXOsMap := make(map[string][]int)
	for{
		block := blcIterator.Next()

		for i := len(block.Txs)-1;i>=0;i--{
			txOutputs := &TXOutputs{[]*UTXO{}}
			tx := block.Txs[i]
			txHash := hex.EncodeToString(tx.TxHash)
			//coinbase
			//添加记录已花费的
			if tx.IsCoinBaseTransaction()==false{
				for _,txInput:=range tx.Vins{
					txInputHash := hex.EncodeToString(txInput.TXHash)
					spentUTXOsMap[txInputHash] = append(spentUTXOsMap[txInputHash],txInput)
					//1.spentUTXOsMap[txInputHash] = append(spentUTXOsMap[txInputHash],txInput.Vout)
				}
			}
			for index,out := range tx.Vouts{

				txInputs := spentUTXOsMap[txHash]

				if len(txInputs) > 0 {
					//是否消费，默认没有
					flag := false
					for _,in := range txInputs{
						outPublicKey:= out.Ripemd160Hash
						inPublicKey := in.Pubkey
						if bytes.Compare(outPublicKey,Ripemd160Hash(inPublicKey))==0 && index==in.Vout{
							flag = true
							break
						}
					}
					if flag == false{
						utxo := &UTXO{tx.TxHash,index,out}
						txOutputs.UTXOS = append(txOutputs.UTXOS,utxo)
					}
				}else{
					utxo := &UTXO{tx.TxHash,index,out}
					txOutputs.UTXOS = append(txOutputs.UTXOS,utxo)

				}
			}
			//设置键值对
			utxoMap[txHash] = txOutputs

		}
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0))==0{
			break
		}
	}
	return utxoMap
}