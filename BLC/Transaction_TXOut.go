package BLC

import (
	"bytes"
)

type TXOutput struct {
	Value int64
	Ripemd160Hash []byte //用户名

}
func (out *TXOutput) Lock(address string){
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-4]
	out.Ripemd160Hash = pubKeyHash
}
//判断output是否是address的
func (txOutput *TXOutput)UnLockWithAddress(address string) bool{
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-4]

	return bytes.Compare(pubKeyHash,txOutput.Ripemd160Hash)==0
}

func NewTXOutput(value int64,address string) *TXOutput{
	txo := &TXOutput{value,nil}
	txo.Lock(address)
	return txo
}