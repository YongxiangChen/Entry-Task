package rpc

import (
	"entrytask1/httpserver/conf"
	"log"
	"net"
)

var pool chan net.Conn

func init() {
	pool = make(chan net.Conn, conf.PoolSize)
	for {
		if len(pool) > conf.PoolSize {
			log.Println("initialize pool success")
			break
		}
		conn, err := net.Dial("tcp", conf.RpcAddr)
		if err != nil {
			log.Println("connect RPC error", err)
			continue
		}
		pool <- conn
	}
}

func GetConn() net.Conn {
	select {
	case conn := <- pool:
		return conn
	}
}

func PutConn(conn net.Conn) {
	select {
	case pool <- conn:
		return
	}
}