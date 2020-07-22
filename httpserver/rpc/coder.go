package rpc

import (
	"bytes"
	"encoding/gob"
)

type RpcData struct {
	Name string //函数名
	Args []interface{} //函数入参或返回值
}

func Encode(data RpcData) ([]byte, error) {
	var buf bytes.Buffer
	// 新建编码器
	enc := gob.NewEncoder(&buf)
	// 编码
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	// 返回字节切片
	return buf.Bytes(), nil
}

func Decode(data []byte) (RpcData, error) {
	var buf = bytes.NewBuffer(data)
	// 新建解码器
	dec := gob.NewDecoder(buf)
	var rpcdata RpcData
	// 解码
	err := dec.Decode(&rpcdata)
	if err != nil {
		return rpcdata, err
	}
	return rpcdata, nil
}
