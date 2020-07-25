package main

import (
	"encoding/json"
	"log"
	"strconv"
)

type codeData struct {
	FuncName string
	Args map[string]string
}

func Call(funcname string, inArgs map[string]string) (map[string]string, error) {
	// 编码
	reqData := codeData{FuncName: funcname, Args: inArgs}
	reqData_json, err := json.Marshal(reqData)
	if err != nil {
		log.Println("编码错误")
		return nil, err
	}
	// 发送与接收
	rspData_json := sendtoRpc(reqData_json)

	// 解码
	var rspData codeData
	err = json.Unmarshal(rspData_json, &rspData)
	if err != nil {
		log.Println("解码错误")
		return nil, err
	}

	return rspData.Args, nil

}

func sendtoRpc(jsonData []byte) []byte {
	// 解码
	var reqData codeData
	err := json.Unmarshal(jsonData, &reqData)
	if err != nil {
		log.Println("解码错误")
		return nil
	}

	var rspData codeData
	rspData.FuncName = reqData.FuncName

	// 调用方法
	if reqData.FuncName == "add" {
		one, _ := strconv.Atoi(reqData.Args["one"])
		two, _ := strconv.Atoi(reqData.Args["two"])
		outArgs := make(map[string]string)
		result := Add(one, two)
		outArgs["result"] = strconv.Itoa(result)
		rspData.Args = outArgs
	} else {
		one, _ := strconv.Atoi(reqData.Args["one"])
		two, _ := strconv.Atoi(reqData.Args["two"])
		outArgs := make(map[string]string)
		result := Sub(one, two)
		outArgs["result"] = strconv.Itoa(result)
		rspData.Args = outArgs
	}

	// 编码
	rspData_json, err := json.Marshal(rspData)
	if err != nil {
		log.Println("编码错误")
		return nil
	}
	return rspData_json
}

func Add(one int, two int) int {
	return one + two
}
func Sub(one int, two int) int {
	return one - two
}

type Rpc struct {
	Data interface{}
}

type Bird struct {
	Sound string
}

func main() {
	var rpc = Rpc{Data: Add}
	rpc.Data(1, 2)
}

/*
type Array struct {
	A, B int
}

type Bird struct {
	Cool Array
	C int
}

func main() {
	var arr = Bird{Array{1, 2}, 3}
	data, err := json.Marshal(arr)
	if err != nil {
		log.Println("error")
		return
	}
	log.Println(string(data))
	var obj Bird
	json.Unmarshal(data, &obj)
	log.Printf("%+v", obj)
	log.Printf("%T", obj.Cool)

}

 */
