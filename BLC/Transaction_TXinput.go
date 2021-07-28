 package BLC

 import (
	 "bytes"
 )

 type TXInput struct {
	//1.交易的ID
	TXHash []byte
	//2.存储在TXOUtput里的索引（该索引存在于对应TxHash的transaction中）
	Vout int
	//3.数字签名
	Signature []byte
	//4.公钥
	Pubkey []byte
}
//判断input是否是address的
func (txInput *TXInput)UnLockRipemd160Hash(ripemd160Hash []byte) bool{
	publicKey := Ripemd160Hash(txInput.Pubkey)

	return bytes.Compare(publicKey,ripemd160Hash)==0
}