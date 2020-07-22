package rpc

import (
	"crypto/md5"
	"fmt"
	"time"
)


func GetToken(username string) string {
	t := fmt.Sprintf("%v", time.Now().UnixNano()) //获取时间戳
	oriData := username + t[len(t)-6:] //取用户名和时间戳后六位
	hs := md5.Sum([]byte(oriData)) //加密
	md5str := fmt.Sprintf("%x", hs) //换成16进制
	return md5str
}
