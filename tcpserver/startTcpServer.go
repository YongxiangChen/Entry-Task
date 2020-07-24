package main

import (
	"encoding/gob"
	"entrytask1/tcpserver/model"
	"entrytask1/tcpserver/rpc"
)

func number() string {
	return "hello"
}

func main() {
	// 注册结构体
	gob.Register(model.User{})
	// 设置地址
	addr := ":8008"

	// 创建服务端
	server := rpc.NewServer(addr)
	// 服务端注册服务
	server.Register("Authenticate", rpc.Authenticate)
	server.Register("SetToken", rpc.SetToken)
	server.Register("VerifyToken", rpc.VerifyToken)
	server.Register("ChangeNickname", rpc.ChangeNickname)
	server.Register("number", number)
	// 运行服务器
	server.Run()
}