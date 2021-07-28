package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//256位Hash里面前面之少16个0
const  targetBit = 16
type ProofOfWork struct {
	Block *Block //当前要验证的区块
	target *big.Int
}

// 创建新的工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork{
	//1.big.Int对象 1

	//1. 创建一个初始值为1的target
	target := big.NewInt(1)
	target = target.Lsh(target,256-targetBit)
	return &ProofOfWork{block,target}
}

//挖矿
func(proofOfWork *ProofOfWork) Run() ([]byte,int64){
	//1.将BLock凭借成字节数组
	//2.生成hash
	//3.判断hash有效性，如果满足条件，跳出循环。
	nonce := 0
	var hashInt big.Int
	for {
		//准备数据
		dataBytes:=proofOfWork.prepareData(nonce)
		//生成hash
		hash := sha256.Sum256(dataBytes)
		fmt.Printf("\r%x",hash)

		//将hash存储到hashint
		hashInt.SetBytes(hash[:])
		//fmt.Println(hashInt)
		//判断hashINT是否小于target
		if proofOfWork.target.Cmp(&hashInt)==1{
			//成立跳出循环
			fmt.Printf("\n")
			return hash[:],int64(nonce)
		}
		nonce = nonce+1
	}

}

// 数据拼接->[]byte
func(pow *ProofOfWork) prepareData(nonce int)[]byte{
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.HashTransactions(),
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(nonce)),
			IntToHex(int64(pow.Block.Height)),
		},
		[]byte{},
		)
	return data
}

// 判断hash是否有效->bool
func (proofOfWork *ProofOfWork) isValid() bool{

	// 1.proofOfWork.Block.Hash
	// 2.proofOfWork.Target
	var hashInt big.Int
	hashInt.SetBytes(proofOfWork.Block.Hash)

	if proofOfWork.target.Cmp(&hashInt)==1{
		return true
	}else{
		return false
	}

}
