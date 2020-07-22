package rpc

import (
	"log"
	"net"
	"reflect"
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
func (client *Client) RpcCall(name string, fpoint interface{}) {
	fn := reflect.ValueOf(fpoint).Elem()
	f := func(args []reflect.Value) (result []reflect.Value) {
		// 构造RpcData类型
		inArgs := make([]interface{}, 0, len(args))
		for _, v := range args {
			inArgs = append(inArgs, v.Interface())
		}
		rpcdata := RpcData{Name: name, Args: inArgs}
		// 编码
		bytedata, err := Encode(rpcdata)
		if err != nil {
			log.Println("RPC error: 编码错误")
			return
		}

		// 新建会话
		session := NewSession(client.conn)
		// 发送数据
		err = session.Write(bytedata)
		if err != nil {
			log.Println("RPC error: 发送数据失败")
			return
		}
		// 接收客户端数据
		bytedata, err = session.Read()
		if err != nil {
			log.Println("RPC error: 接收数据失败")
			return
		}
		//解码
		rpcdata, err = Decode(bytedata)
		if err != nil {
			log.Println("RPC error: 解码错误")
			return
		}

		// 处理服务器的数据
		outArgs := make([]reflect.Value, 0, len(rpcdata.Args))
		for i, v := range rpcdata.Args {
			// 数据特殊情况处理
			if v == nil {
				// reflect.Zero() 返回某类型的零值的value
				// .Out()返回函数输出的参数类型
				// 得到具体第几个位置的参数的零值
				outArgs = append(outArgs, reflect.Zero(fn.Type().Out(i)))
				continue
			}
			outArgs = append(outArgs, reflect.ValueOf(v))
		}
		return outArgs
	}

	v := reflect.MakeFunc(fn.Type(), f)
	fn.Set(v)
}