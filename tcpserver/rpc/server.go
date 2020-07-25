package rpc

import (
	"log"
	"net"
)

type Server struct {
	addr string //连接地址
	funcMap map[string]func(map[string]string) map[string]string
}

func NewServer(addr string) *Server {
	return &Server{addr:addr, funcMap: make(map[string]func(map[string]string) map[string]string)}
}

func (s *Server) Register(name string, fn func(map[string]string) map[string]string) {
	if _, ok := s.funcMap[name]; ok {
		return
	}
	s.funcMap[name] = fn
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
		reqData, err := Decode(bytedata)
		if err != nil {
			return
		}

		var rspData RpcData
		rspData.FuncName = reqData.FuncName
		// 调用方法
		if _, ok := s.funcMap[reqData.FuncName]; !ok {
			return
		}
		fn := s.funcMap[reqData.FuncName]
		outArgs := fn(reqData.Args)
		rspData.Args = outArgs

		// 编码
		bytedata, err = Encode(rspData)
		if err != nil {
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