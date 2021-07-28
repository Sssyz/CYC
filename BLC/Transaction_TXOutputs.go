package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TXOutputs struct {
	UTXOS []*UTXO
}
// 将Outputs序列化->[]byte
func (txOutPuts *TXOutputs) Serialize()[]byte{
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(txOutPuts)
	if err!=nil{
		log.Panic(err)
	}
	return result.Bytes()
}

// 反序列化->Outputs
func DeSerializeTxOutputs(txOutputsByte []byte) *TXOutputs{
	var txOutputs TXOutputs

	decoder := gob.NewDecoder(bytes.NewReader(txOutputsByte))
	err:=decoder.Decode(&txOutputs)
	if err !=nil {
		panic(err)
	}
	return &txOutputs
}