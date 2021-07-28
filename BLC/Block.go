package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

)

type Block struct {
	//1.区块高度
	Height int64
	//2.上一个区块的Hash
	PrevBlockHash []byte
	//3.交易数据
	Txs []*Transaction
	//4.时间戳
	Timestamp int64
	//5.Hash
	Hash []byte
	//6. Nonce
	Nonce int64
}
//func(block *Block)SetHash(){
//
//	// 1.height转化为字节数组
//	heightBytes := IntToHex(block.Height)
//
//
//	// 2.将时间戳转化为字节数组
//	timeString := strconv.FormatInt(block.Timestamp,2)
//	timeBytes := []byte(timeString)
//
//	// 3.拼接所有属性
//	blockBytes := bytes.Join([][]byte{heightBytes,block.PrevBlockHash,block.Data,timeBytes},[]byte{})
//
//	// 4，生成Hash
//	hash := sha256.Sum256(blockBytes)
//
//	block.Hash = hash[:]
//}
//1. 创建新的区块
func NewBlock(txs []*Transaction,height int64,prevBlockHash []byte) *Block{


	// 创建区块
	block := &Block{
		height,prevBlockHash,txs,time.Now().Unix(),nil,0,
	}

	// 设置Hash值
	//block.SetHash()
	//调用工作量证明方法，并且返回有效的Hash和Nonce值
	pow := NewProofOfWork(block)
	//挖矿验证
	hash,nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

//2. 单独写一个方法，生成创世区块

func CreateGenesisBlock(txs []*Transaction) *Block{

	return NewBlock(txs,1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}

//将交易数据换为字节数组-》[]byte
func (block *Block) HashTransactions() []byte{
	//var txHash []byte
	////合并所有TXs(*[]Transaction)中的TxHash
	//for _,tx:=range block.Txs{
	//	txHash = append(txHash,tx.TxHash...)
	//}
	//
	//return txHash
	var transactions [][]byte
	for _,tx := range block.Txs{
		transactions = append(transactions,tx.Hash())

	}
	mTree := NewMerkleTree(transactions)
	return mTree.RootNode.Data
}
// 将区块序列化->[]byte
func (block *Block) Serialize()[]byte{
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err!=nil{
		log.Panic(err)
	}
	return result.Bytes()
}

// 反序列化->Block
func DeSerialize(blockByte []byte) *Block{
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockByte))
	err:=decoder.Decode(&block)
	if err !=nil {
		panic(err)
	}
	return &block
}