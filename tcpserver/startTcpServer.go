package main

import (
	"entrytask1/tcpserver/rpc"
)

func number() string {
	return "hello"
}

func main() {
	// 设置地址
	addr := ":8008"

	// 创建服务端
	server := rpc.NewServer(addr)

	// 注册服务
	server.Register("Authenticate", rpc.Authenticate)
	server.Register("SetToken", rpc.SetToken)
	server.Register("ChangeNickname", rpc.ChangeNickname)
	server.Register("VerifyToken", rpc.VerifyToken)

	// 运行服务器
	server.Run()
}