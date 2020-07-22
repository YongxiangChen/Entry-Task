package server

import (
	"entrytask1/httpserver/conf"
	"entrytask1/httpserver/rpc"
	"entrytask1/tcpserver/model"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 登陆
func login(w http.ResponseWriter, r *http.Request) {
	HttpLog(r.Method, r.URL.Path, r.RemoteAddr)
	if r.Method == "GET" {
		var tpath = filepath.Join(conf.StaticPath, "login.html")
		t, _ := template.ParseFiles(tpath) //解析static的html文件，路径表示有点问题
		t.Execute(w, nil)
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Fatal("ParseForm ", err)
		}

		//构造user结构体
		u := model.User{
			Username: r.Form["username"][0],
			Password: r.Form["password"][0],
		}
		//backend验证用户名和密码，调用rpc服务
		var auth func(name string, pw string) (model.User, bool)
		conn, err := RpcService(conf.RpcAddr, "Authenticate", &auth)
		if err != nil {
			log.Println("RPC error：客户端建立连接失败")
		}
		u, ok := auth(r.Form["username"][0], r.Form["password"][0])
		conn.Close()
		if ok == false {
			//做一些登陆不通过的事
			fmt.Println("登陆失败")
			fmt.Fprint(w, "登陆失败")
			return
		}

		//调用rpc服务，设置token，存储到redis
		var settoken func(user model.User) (string, error)
		conn, err = RpcService(conf.RpcAddr, "SetToken", &settoken)
		if err != nil {
			log.Println("RPC error：客户端建立连接失败")
		}
		tk, err := settoken(u)
		conn.Close()
		if err != nil {
			log.Println("连接redis错误")
		}
		w.Header().Set("Set-Cookie", "Token:"+tk)

		//设置重定向
		w.Header().Set("Location", "/userhome")//跳转地址设置
		w.WriteHeader(302)
	}
}

//进入个人主页
func userhome(w http.ResponseWriter, r *http.Request) {
	HttpLog(r.Method, r.URL.Path, r.RemoteAddr)
	if r.Method == "GET" {
		// 对在cookie中对token进行验证
		cookie := string(r.Header.Get("Cookie")) //[]unit8 to string
		var tk string
		if cookie == "" {
			tk = ""
		} else {
			tk = strings.Split(cookie, ":")[1]
		}

		// 调用Rpc服务，验证tk，确认用户是否登陆
		var verify func(tk string) (model.User, int)
		conn, err := RpcService(conf.RpcAddr, "VerifyToken", &verify)
		if err != nil {
			log.Println("RPC error：客户端建立连接失败")
		}
		user, ok := verify(tk)
		conn.Close()
		if ok == 1 {
			//未登陆
			fmt.Fprintf(w, "unauthorized")
			return
		} else if ok == 2 {
			log.Printf("数据库错误")
			return
		}

		//登陆成功
		log.Println("User log in: ", user)
		var tpath = filepath.Join(conf.StaticPath, "userhome.html")
		t, _ := template.ParseFiles(tpath)
		// 补充昵称
		t.Execute(w, user.Nickname)
	} else {
		if r.Header["Content-Type"][0] == "application/x-www-form-urlencoded" {
			// 用户是修改昵称

		}
		// 用户上传图片
		r.ParseMultipartForm(32 << 20) // 位运算，32MB
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)

		dirname := filepath.Join(conf.MediaPath, handler.Filename) //路径和文件名拼接
		f, err := os.OpenFile(dirname, os.O_WRONLY|os.O_CREATE, 0666) //打开目标文件等待写入，这里后期把文件名换成用户名相关的
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

// 主页
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "欢迎来到，用户管理系统")
}

// 打印格式
func HttpLog(method string, path string, addr string) {
	log.Println(method, " --> ", path, " By ", addr)
}

// 调用rpc的服务
func RpcService(addr string, funcname string, fpoint interface{}) (net.Conn, error) {
	// 创建客户端连接
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return conn, err
	}
	// 创客户端
	client := rpc.NewClient(conn)
	// 定义函数调用原型
	// 客户端调用rpc
	client.RpcCall(funcname, fpoint)
	return conn, nil
}