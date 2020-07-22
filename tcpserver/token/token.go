package token

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

// 根据用户名生成token
func GetToken(username string) string {
	// 获取时间戳并转为字符串,取其后6位
	timeUnixNano := time.Now().UnixNano()
	timeToString := strconv.Itoa(int(timeUnixNano))
	timeToString = timeToString[len(timeToString)-6:]
	// 用户名和时间戳合在一起
	usernameTime := username + timeToString
	// md5编码
	data := []byte(usernameTime)
	token := md5.Sum(data)
	return fmt.Sprintf("%x", token)
}

