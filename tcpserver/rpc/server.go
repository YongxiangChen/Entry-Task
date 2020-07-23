package rpc

import (
	"log"
	"net"
	"reflect"
)

type Server struct {
	addr string //连接地址
	funcs map[string]reflect.Value //健是函数名，值是reflect.ValueOf(fn)，其中fn是注册的函数
}

func NewServer(addr string) *Server {
	return &Server{addr:addr, funcs:make(map[string]reflect.Value)}
}

// 注册函数
func (s *Server) Register(name string, fn interface{}) {
	if _, ok := s.funcs[name]; ok {
		// 这个函数已经添加过了
		return
	}
	s.funcs[name] = reflect.ValueOf(fn)
}


// 运行服务器
func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Println("创建连接失败")
		return err
	}
	for {
		// 接收客户端连接
		conn, err := listener.Accept()
		if err != nil {
			log.Println("连接异常")
		}
		// 将连接放到新线程中处理
		go handelConn(conn, s)
	}
}

func handelConn(conn net.Conn, s *Server) {
	for {
		session := NewSession(conn)
		// 从连接中读数据
		bytedata, err := session.Read()
		if err != nil {
			log.Println("连接异常")
			return
		}
		// 把数据解码，成为Rpcdata类型
		rpcdata, err := Decode(bytedata)
		if err != nil {
			log.Println("解码异常")
			return
		}
		// 从映射中找到函数，它是reflect.Value类型
		fn, ok := s.funcs[rpcdata.Name]
		if !ok {
			log.Println("函数找不到")
			return
		}
		// 将入参转换为reflect.Value类型，然后放入[]reflect.Value中
		inArgs := make([]reflect.Value, 0, len(rpcdata.Args))
		for _, v := range rpcdata.Args {
			inArgs = append(inArgs, reflect.ValueOf(v))
		}
		// 调用函数，返回的是[]reflect.Value
		returnData := fn.Call(inArgs)
		log.Printf("调用函数 %q 成功!", rpcdata.Name)
		// 构造RpcData的Args成员数据
		outArgs := make([]interface{}, 0, len(returnData))
		for _, rv := range returnData {
			outArgs = append(outArgs, rv.Interface())
		}
		// 构造RpcData
		rspdata := RpcData{Name: rpcdata.Name, Args: outArgs}
		// 编码
		bytedata, err = Encode(rspdata)
		if err != nil {
			log.Println("编码失败")
			return
		}
		// 发送数据
		err = session.Write(bytedata)
		if err != nil {
			log.Println("数据发送失败")
			return
		}
	}
}