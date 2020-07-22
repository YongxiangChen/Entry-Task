// 新建网络会话
// 网络会话的读和写

package rpc

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// 处理连接会话

// 会话对象结构体
type Session struct {
	conn net.Conn
}

// 传输数据存储方式
// 字节数组， 添加4个字节的头，用来存储数据的长度

// 会话构造函数
func NewSession(conn net.Conn) *Session {
	return &Session{conn: conn}
}

// 从连接中读取数据
func (s *Session) Read() (data []byte, err error) {
	// 读取数据header数据
	header := make([]byte, 4)
	_, err = s.conn.Read(header)
	if err != nil {
		fmt.Printf("read conn header data failed, err: %v\n", err)
		return
	}
	// 根据header数据（body的长度）读取body数据
	hlen := binary.BigEndian.Uint32(header)
	data = make([]byte, hlen)
	_, err = io.ReadFull(s.conn, data)
	if err != nil {
		fmt.Printf("read conn body data failed, err: %v\n", err)
		return
	}
	return
}

// 向连接中写入数据
func (s *Session) Write(data []byte) (err error) {
	// 创建数据字节切片
	buf := make([]byte, 4+len(data))
	// 向header写入数据长度
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	// 写入body内容
	copy(buf[4:], data)
	// 写入连接数据
	_, err = s.conn.Write(buf)
	if err != nil {
		fmt.Printf("write conn data failed, err: %v\n", err)
		return
	}
	return
}