package rpc

import (
	"entrytask1/tcpserver/dao"
	"entrytask1/tcpserver/model"
	"entrytask1/tcpserver/token"
	"log"
	"strconv"
)

// 验证用户名和密码,返回code错误码
// reqData{"username", "password"}
// rspData{"userid", "code"}
func Authenticate(reqData map[string]string) map[string]string {
	userdb := dao.NewUserDB()
	//defer userdb.Close()

	//先根据用户名查出有无用户
	var user = &model.User{Username: reqData["username"]}
	ok, err := userdb.UserQueryByName(user)
	var rspData = make(map[string]string)
	if err != nil {
		if ok == 1 {
			//运行正常，查无此人
			rspData["code"] = "1"
			rspData["userid"] = "-1"
			return rspData
		} else {
			//运行异常
			rspData["code"] = "2"
			rspData["userid"] = "-1"
			log.Println("mysql error!")
			log.Println(err)
			return rspData
		}
	}
	//查到了用户，对比密码是否一致
	if reqData["password"] != user.Password {
		rspData["code"] = "3"
		rspData["userid"] = "-1"
		return rspData
	}

	rspData["userid"] = strconv.Itoa(user.Id)
	rspData["code"] = "0"
	return rspData
}



//存储token到redis并且返回token值
// reqData{"username", "userid"}
// rspData{"token"}
func SetToken(reqData map[string]string) map[string]string {
	tk := token.GetToken(reqData["username"])
	userid, _ := strconv.Atoi(reqData["userid"])
	auth := model.Auth{
		Token: tk,
		Userid: userid,
	}
	client := dao.NewRedisClient()
	defer client.Close()
	client.SetValueRedis(&auth)

	rspData := make(map[string]string)
	rspData["token"] = tk
	return rspData
}

// 修改用户昵称, 返回0是正常，返回1是连接redis错误， 返回2是查不到键
// reqData{"nickname", "token"}
// rspData{"code"}
func ChangeNickname(reqData map[string]string) map[string]string {
	client := dao.NewRedisClient()
	auth := model.Auth{Token: reqData["token"]}
	ok, _ := client.GetValueRedis(&auth)
	var rspData = make(map[string]string)
	if ok != 0 {
		rspData["code"] = strconv.Itoa(ok)
		return rspData
	}
	user := model.User{Id: auth.Userid, Nickname: reqData["nickname"]}
	userDB := dao.NewUserDB()
	ok, _ = userDB.UserUpdateById(&user)
	rspData["code"] = strconv.Itoa(ok)
	return rspData
}

//验证token，返回nickname, username, 错误码(1用户登陆问题，2数据库错误，0正常)
// reqData{"token"}
// rspData{"nickname", "username", "code"}
func VerifyToken(reqData map[string]string) map[string]string {
	client := dao.NewRedisClient()
	auth := model.Auth{Token: reqData["token"]}
	ok, err := client.GetValueRedis(&auth)
	defer client.Close()
	rspData := make(map[string]string)
	switch ok {
	case 1:
		//找不到token，未登陆
		rspData["nickname"] = ""
		rspData["username"] = ""
		rspData["code"] = "1"
		log.Println("未登陆:, ", err)
		return rspData
	case 2:
		//redis连接错误
		rspData["nickname"] = ""
		rspData["username"] = ""
		rspData["code"] = "2"
		log.Println("redis error: ", err)
		return rspData
	}

	user := model.User{Id: auth.Userid}
	// 根据id查
	userDB := dao.NewUserDB()
	ok, err = userDB.UserQueryById(&user)
	switch ok {
	case 1:
		//找不到user，id不匹配
		rspData["nickname"] = ""
		rspData["username"] = ""
		rspData["code"] = "1"
		log.Println("找不到user:, ", err)
		return rspData
	case 2:
		//mysql连接错误
		rspData["nickname"] = ""
		rspData["username"] = ""
		rspData["code"] = "2"
		log.Println("mysql error: ", err)
		return rspData
	}
	rspData["nickname"] = user.Nickname
	rspData["username"] = user.Username
	rspData["code"] = "0"
	return rspData
}