package rpc

import (
	"entrytask1/tcpserver/dao"
	"entrytask1/tcpserver/model"
	"entrytask1/tcpserver/token"
	"log"
)

// 验证用户名和密码
func Authenticate(username string, password string) (model.User, bool) {
	userdb := dao.NewUserDB()
	//defer userdb.Close()

	//先根据用户名查出有无用户
	var user = &model.User{Username: username}
	ok, err := userdb.UserQueryByName(user)
	if err != nil {
		if ok == 1 {
			//运行正常，查无此人
			return model.User{} ,false
		} else {
			//运行异常
			log.Println("mysql error!")
			return model.User{} ,false
		}
	}
	//查到了用户，对比密码是否一致
	if password != user.Password {
		return model.User{}, false
	}

	return *user, true
}



//存储token到redis并且返回token值
func SetToken(user model.User) (string, error) {
	tk := token.GetToken(user.Username)
	auth := model.Auth{
		Token: tk,
		Userid: user.Id,
	}
	client := dao.NewRedisClient()
	defer client.Close()
	err := client.SetValueRedis(&auth)
	if err != nil {
		return "", err
	}

	return tk, nil
}

// 修改用户昵称
func ChangeNickname(tk string, name string) int {
	client := dao.NewRedisClient()
	auth := model.Auth{Token: tk}
	ok, _ := client.GetValueRedis(&auth)
	if ok != 0 {
		return ok
	}
	user := model.User{Id: auth.Userid, Nickname: name}
	userDB := dao.NewUserDB()
	ok, _ = userDB.UserUpdateById(&user)
	return ok
}

//验证token，返回user对象
func VerifyToken(tk string) (model.User, int) {
	client := dao.NewRedisClient()
	auth := model.Auth{Token: tk}
	ok, err := client.GetValueRedis(&auth)
	defer client.Close()
	switch ok {
	case 1:
		//找不到token，未登陆
		log.Println("未登陆:, ", err)
		return model.User{}, 1
	case 2:
		//redis连接错误
		log.Println("redis error: ", err)
		return model.User{}, 2
	}

	user := model.User{Id: auth.Userid}
	// 根据id查
	userDB := dao.NewUserDB()
	ok, err = userDB.UserQueryById(&user)
	switch ok {
	case 1:
		//找不到user，id不匹配
		log.Println("找不到user:, ", err)
		return model.User{}, 1
	case 2:
		//mysql连接错误
		log.Println("mysql error: ", err)
		return model.User{}, 2
	}
	return user, 0
}