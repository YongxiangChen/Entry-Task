package rpc


/*
import (
	"entrytask1/httpserver/conf"
	"log"
	"net"
)


var pool chan net.Conn

func init() {
	pool = make(chan net.Conn, conf.PoolSize)
	for {
		if len(pool) > conf.PoolSize - 1 {
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
	log.Println(len(pool))
}

func GetConn() net.Conn {
	conn := <- pool
	return conn
}

func PutConn(conn net.Conn) {
	pool <- conn
}

 */