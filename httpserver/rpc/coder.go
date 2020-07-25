package rpc

import (
	"encoding/json"
	"log"
)

type RpcData struct {
	FuncName string //函数名
	Args map[string]string
}

func Encode(data RpcData) ([]byte, error) {
	dataJson, err := json.Marshal(data)
	if err != nil {
		log.Println("编码错误")
		return nil, err
	}
	return dataJson, nil
}

func Decode(dataJson []byte) (RpcData, error) {
	var data RpcData
	err := json.Unmarshal(dataJson, &data)
	if err != nil {
		log.Println("解码错误")
		return RpcData{}, err
	}

	return data, nil
}
