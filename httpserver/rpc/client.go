package rpc

import (
	"net"
)

type Client struct {
	conn net.Conn
}

// 构造新客户端
func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

// 通过Rpc调用方法
// client.RpcCall("login", &req)
func (client *Client) RpcCall(name string, inArgs map[string]string) (map[string]string, error) {
	//编码
	var reqData RpcData
	reqData.FuncName = name
	reqData.Args = inArgs
	reqData_json, err := Encode(reqData)
	if err != nil {
		return nil, err
	}

	// 发送数据
	session := NewSession(client.conn)
	err = session.Write(reqData_json)
	if err != nil {
		return nil, err
	}

	// 接收数据
	rspData_json, err := session.Read()
	if err != nil {
		return nil, err
	}

	// 解码
	rspData, err := Decode(rspData_json)
	if err != nil {
		return nil, err
	}
	return rspData.Args, nil
}