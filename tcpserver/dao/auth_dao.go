package dao

import (
	"entrytask1/tcpserver/model"
	redigo "github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)
var RedisClient *redigo.Pool

func init() {
	var (
		addr = "127.0.0.1:6379"
		password = ""
	)
	RedisClient = PoolInitRedis(addr, password)
}

// redis pool
func PoolInitRedis(server string, password string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     300,//空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   1000,//最大数
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}


type AuthRedis struct {
	Pool *redigo.Pool
}

func NewRedisClient() AuthRedis {
	return AuthRedis{RedisClient}
}

func (client AuthRedis) SetValueRedis(auth *model.Auth) error {
	rc := client.Pool.Get()
	defer rc.Close()
	//设置键的过期时间为1个小时
	_, err := rc.Do("Set", auth.Token, auth.Userid, "EX", "3600")
	if err != nil {
		return err
	}
	return nil
}

func (client AuthRedis) GetValueRedis(auth *model.Auth) (int, error) {
	rc := client.Pool.Get()
	defer rc.Close()
	//返回0是正常，返回1是错误
	val, err := redigo.String(rc.Do("Get", auth.Userid))
	if err != nil {
		return 1, err
	}
	//找到了健值，赋值给auth结构体
	auth.Userid, _ = strconv.Atoi(val)
	return 0, err
}