package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)
//转化为16进制字节数组
func IntToHex(x int64)[]byte{

	buff := new(bytes.Buffer)
	err := binary.Write(buff,binary.BigEndian,x)
	if err!=nil{
		log.Panic(err)
	}
	return buff.Bytes()
}

func JSONToArray(jsonString string)[]string{
	var sArr []string
	jsonString = strings.Replace(jsonString,`'`,`"`,-1)
	//fmt.Println(jsonString)
	if err := json.Unmarshal([]byte(jsonString),&sArr);err!=nil{
		log.Panic(err)
	}
	return sArr
}
//字节数组反转
func ReverseBytes(data []byte){
	for i,j := 0,len(data)-1;i<j;i,j=i+1,j-1{
		data[i],data[j] = data[j],data[i]
	}
}
//获取节点id
func GetEnvNodeId()string {

	nodeId:= os.Getenv("NODE_ID")
	if nodeId==""{
		fmt.Println("NODE_ID is not set...")
		os.Exit(1)
	}
	return nodeId
}

// gob编码
func gobEncode(data interface{})[]byte{
	var result bytes.Buffer

	enc := gob.NewEncoder(&result)
	err := enc.Encode(data)
	if nil!=err{
		log.Panicf("encode the data failed! %v\n",err)
	}
	return result.Bytes()
}

// 命令转换为请求([]byte)
func commandToByte(command string)[]byte{
	var bytes [CMMAND_LENGTH]byte
	for i,c:=range command{
		bytes[i]=byte(c)
	}
	return bytes[:]
}
// 反解析，把请求中的命令解析出
func bytesToCommand(bytes []byte) string{
	var command []byte
	for _,b := range bytes{
		if b!=0x00{
			command = append(command,b)
		}
	}
	return fmt.Sprintf("%s",command)

}