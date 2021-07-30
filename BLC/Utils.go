package BLC

import (
	"bytes"
	"encoding/binary"
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