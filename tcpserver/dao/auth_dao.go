package dao

import (
	"context"
	"entrytask1/tcpserver/model"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type AuthRedis struct {
	*redis.Client
}

var ctx = context.Background()

func NewRedisClient() AuthRedis {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	return AuthRedis{rdb}
}

func (client AuthRedis) SetValueRedis(auth *model.Auth) error {
	//设置键的过期时间为2个小时
	err := client.SetNX(ctx, auth.Token, auth.Userid, 2*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client AuthRedis) GetValueRedis(auth *model.Auth) (int, error) {
	//返回0是正常，返回1是连接redis错误， 返回2是查不到键
	val, err := client.Get(ctx, auth.Token).Result()
	if err != nil {
		return 1, err
	}
	if err == redis.Nil {
		return 2, err
	}
	//找到了健值，赋值给auth结构体
	auth.Userid, _ = strconv.Atoi(val)
	return 0, err
}