package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"
	"time"
)

type Transaction struct {

	//1. 交易Hash
	TxHash []byte
	//2. 输入
	Vins []*TXInput
	//3. 输出
	Vouts []*TXOutput

}
func (tranction *Transaction) HashTransaction(){

	tranction.TxHash = tranction.Hash()
}

func(tx *Transaction) IsCoinBaseTransaction() bool{
	return len(tx.Vins[0].TXHash)==0&&tx.Vins[0].Vout==-1
}




//1. Transaction的创建分两种情况
//（1）创世区块
func NewCoinbaseTransaction(address string) *Transaction{
	//代表输入（消费记录）
	txInput := &TXInput{[]byte{},-1,nil,[]byte{}}
	//添加一笔钱
	txOutput := NewTXOutput(10,address)
	//txOutput := &TXOutput{10,address}
	txCoinbase := &Transaction{[]byte{},[]*TXInput{txInput},[]*TXOutput{txOutput}}

	//设置Hash值
	txCoinbase.HashTransaction()

	return txCoinbase
}
//（2）转账时的transaction

func NewSimpleTransaction(from string,to string,amount int64,utxoSet *UTXOSet,txs []*Transaction)*Transaction{

	var txInputs []*TXInput
	var txOutputs []*TXOutput
	wallets := NewWallets()
	wallet := wallets.WalletsMap[from]
	money,spendableUTXODic := utxoSet.FindSpendableUTXOs(from,amount,txs)
	//消费

	for hash,indexArray:=range spendableUTXODic{
		txHashBytes,_ := hex.DecodeString(hash)
		for _,index := range indexArray{
			txInput := &TXInput{txHashBytes,index,nil,wallet.PublicKey}

			txInputs = append(txInputs,txInput)
		}

	}


	//转账
	txOutput := NewTXOutput(int64(amount),to)
	//txOutput := &TXOutput{int64(amount),to}
	txOutputs = append(txOutputs,txOutput)
	//找零
	change:=int64(money)-int64(amount)
	//如果找零为0，不用存储
	if  change>int64(0){
		//txOutput = &TXOutput{int64(money)-int64(amount),from}
		txOutput = NewTXOutput(int64(money)-int64(amount),from)
		txOutputs = append(txOutputs,txOutput)
	}

	tx := &Transaction{[]byte{},txInputs,txOutputs}
	tx.HashTransaction()
	//进行数字签名
	utxoSet.BlockChain.SignTransaction(tx,wallet.PrivateKey,txs)
	return tx

}
func (tx *Transaction ) Sign(private ecdsa.PrivateKey,prevTxs map[string]Transaction)  {
	if tx.IsCoinBaseTransaction(){
		return
	}
	for _,vin := range tx.Vins{
		if prevTxs[hex.EncodeToString(vin.TXHash)].TxHash==nil{
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID,vin := range txCopy.Vins{
		prevTx := prevTxs[hex.EncodeToString(vin.TXHash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].Pubkey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		// 签名代码
		r,s,err := ecdsa.Sign(rand.Reader,&private,txCopy.TxHash)
		if err != nil{
			log.Panic(err)
		}

		signature := append(r.Bytes(),s.Bytes()...)
		tx.Vins[inID].Signature = signature
	}



}

func(tx *Transaction)TrimmedCopy() Transaction{
	var Vins []*TXInput

	for _,in := range tx.Vins{
		Vins = append(Vins,&TXInput{in.TXHash,in.Vout,nil,nil})
	}
	return Transaction{tx.TxHash,Vins,tx.Vouts}
}

func(tx *Transaction) Hash()[]byte{
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err!=nil{
		log.Panic(err)
	}
	hash :=sha256.Sum256(result.Bytes())

	// 挖矿奖励，防止txhash重复，添加时间戳
	if tx.Vins[0].Vout==-1{
		ttime := time.Now().Unix()
		data := bytes.Join(
			[][]byte{
				hash[:],
				IntToHex(ttime),
			},
			[]byte{},
		)
		hash = sha256.Sum256(data)
	}
	return hash[:]
}
//验证数字签名
func(tx *Transaction) Verify(prevTxs map[string]Transaction) bool{

	if tx.IsCoinBaseTransaction(){
		return true
	}
	for _,vin := range tx.Vins{
		if prevTxs[hex.EncodeToString(vin.TXHash)].TxHash==nil{
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()
	for inID,vin := range txCopy.Vins{
		prevTx := prevTxs[hex.EncodeToString(vin.TXHash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].Pubkey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()

		//私钥id
		r := big.Int{}
		s := big.Int{}

		sigLen := len(tx.Vins[inID].Signature)


		r.SetBytes(tx.Vins[inID].Signature[:(sigLen/2)])
		s.SetBytes(tx.Vins[inID].Signature[(sigLen/2):])
		x := big.Int{}
		y := big.Int{}

		keyLen := len(tx.Vins[inID].Pubkey)

		x.SetBytes(tx.Vins[inID].Pubkey[:keyLen/2])
		y.SetBytes(tx.Vins[inID].Pubkey[keyLen/2:])

		rawPubKey := ecdsa.PublicKey{curve,&x,&y}
		if ecdsa.Verify(&rawPubKey,txCopy.TxHash,&r,&s)==false{
			return false
		}

	}
	return true
}