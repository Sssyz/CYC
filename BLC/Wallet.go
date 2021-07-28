package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const version = byte(0x00)
const addressChecksumLen = 4
//钱包 存储私钥公钥
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}
//创建一个钱包
func NewWallet() *Wallet{
	private,public := newKeyPair()
	fmt.Println(private)
		fmt.Println(public)
		return &Wallet{private,public}
}
//创建公钥私钥
func  newKeyPair() (ecdsa.PrivateKey,[]byte){
	//1.创建私钥
	curve := elliptic.P256()
	private,err := ecdsa.GenerateKey(curve,rand.Reader)
	if err!=nil{
		log.Panic(err)
	}
	//2.通过私钥产生公钥
	pubKey := append(private.PublicKey.X.Bytes(),private.PublicKey.Y.Bytes()...)
	return *private,pubKey
}
//判断地址有效性
func IsValidForAddress(address []byte) bool{

	version_public_checkSumBytes := Base58Decode(address)
	//fmt.Println(version_public_checkSumBytes)
	checkSumBytes := version_public_checkSumBytes[len(version_public_checkSumBytes)-addressChecksumLen:]
	version_ripemd160Hash := version_public_checkSumBytes[:len(version_public_checkSumBytes)-addressChecksumLen]
	checkBytes:=CheckSum(version_ripemd160Hash)
	//fmt.Println(checkSumBytes)
	//fmt.Println(checkBytes)
	if bytes.Compare(checkBytes,checkSumBytes)==0{
		return true
	}else{
		return false
	}

}
func (w *Wallet) GetAddress() []byte{
	//1.先将公钥进行HASH256 -> hash160
	ripemd160Hash := Ripemd160Hash(w.PublicKey)
	//2.version 160拼接
	version_ripemd160Hash := append([]byte{version},ripemd160Hash...)
	//3.两次256 按sumcheck取
	checkSumBytes := CheckSum(version_ripemd160Hash)
	bytes := append(version_ripemd160Hash,checkSumBytes...)
	//fmt.Println(bytes)
	return Base58Encode(bytes[:])
}
func Ripemd160Hash(publicKey []byte) []byte{
	//256
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)
	//160
	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)
	//fmt.Printf("%s\n",ripemd160.Sum(nil))
	return ripemd160.Sum(nil)
}
func CheckSum(payload []byte)[]byte{
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:addressChecksumLen]
}