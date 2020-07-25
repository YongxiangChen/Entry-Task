package server

import (
	"entrytask1/easypool"
	"entrytask1/httpserver/conf"
	"entrytask1/httpserver/rpc"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var pool easypool.Pool

func init() {
	factory := func() (net.Conn, error) { return net.Dial("tcp", "localhost:8008") }
	config := &easypool.PoolConfig{
		InitialCap:  800,
		MaxCap:      1500,
		MaxIdle:     800,
		Idletime:    10 * time.Second,
		MaxLifetime: 10 * time.Minute,
		Factory:     factory,
	}

	var err error
	pool, err = easypool.NewHeapPool(config)
	if err != nil {
		log.Printf("err:%v\n", err)
		return
	}
	log.Println("pool success")
}

type htmlDetail struct {
	Nickname string
	ImgUrl string
}

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

		// 用户名和密码不为空？
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "" || password == "" {
			//做一些登陆不通过的事
			fmt.Println("登陆失败")
			fmt.Fprint(w, "登陆失败")
			return
		}

		//从pool中拿连接
		conn, err := pool.Get()
		if err != nil {
			log.Printf("err:%v\n", err)
			return
		}

		//调用rpc服务，backend验证用户名和密码
		var reqData = make(map[string]string)
		reqData["username"] = username
		reqData["password"] = password
		rspData, err := RpcService(conn, "Authenticate", reqData)
		if err != nil {
			if _, ok := err.(net.Error); ok {
				// 直接关闭conn，不放回pool
				log.Println("invaild conn")
				conn.(*easypool.PoolConn).MarkUnusable()
				conn.Close()
			} else {
				//做一些登陆不通过的事
				log.Println("登陆失败")
				fmt.Fprint(w, "登陆失败")
				conn.Close()
				return
			}
		}

		// 判断状态码
		if rspData["code"] != "0" {
			//做一些登陆不通过的事
			fmt.Println("登陆失败")
			fmt.Fprint(w, "登陆失败")
			conn.Close()
			return
		}

		// 调用rpc服务，设置token
		userid := rspData["userid"]
		reqData = make(map[string]string)
		reqData["userid"] = userid
		reqData["username"] = username
		rspData, err = RpcService(conn, "SetToken", reqData)
		if err != nil {
			if _, ok := err.(net.Error); ok {
				// 直接关闭conn，不放回pool
				log.Println("invaild conn")
				conn.(*easypool.PoolConn).MarkUnusable()
				conn.Close()
			} else {
				//做一些登陆不通过的事
				log.Println("登陆失败")
				fmt.Fprint(w, "登陆失败")
				conn.Close()
				return
			}
		}

		// cookie设置token
		tk := rspData["token"]
		w.Header().Set("Set-Cookie", "Token:"+tk)

		//设置重定向
		w.Header().Set("Location", "/userhome")//跳转地址设置
		w.WriteHeader(302)
		conn.Close()
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

		//从pool中拿连接
		conn, err := pool.Get()
		if err != nil {
			log.Printf("err:%v\n", err)
			return
		}

		// 调用Rpc服务，验证tk，确认用户是否登陆
		var reqData = make(map[string]string)
		reqData["token"] = tk
		rspData, err := RpcService(conn, "VerifyToken", reqData)
		if err != nil {
			if _, ok := err.(net.Error); ok {
				// 直接关闭conn，不放回pool
				log.Println("invaild conn")
				conn.(*easypool.PoolConn).MarkUnusable()
				conn.Close()
			} else {
				//做一些登陆不通过的事
				log.Println("登陆失败")
				fmt.Fprint(w, "登陆失败")
				conn.Close()
				return
			}
		}

		// 判断状态码
		code := rspData["code"]
		if code != "0"{
			//未登陆
			fmt.Fprintf(w, "unauthorized")
			conn.Close()
			return
		}


		//登陆成功
		username := rspData["username"]
		nickname := rspData["nickname"]
		log.Printf("User log in: username=%v nickname=%v", username, nickname)

		// 构造待补充的信息
		var detail = htmlDetail{Nickname: nickname, ImgUrl: username}

		// 解析html
		var tpath = filepath.Join(conf.StaticPath, "userhome.html")
		t, _ := template.ParseFiles(tpath)
		// 补充昵称
		t.Execute(w, detail)
		//conn放回pool
		conn.Close()
		return
	} else {
		if r.Header["Content-Type"][0] == "application/x-www-form-urlencoded" {
			// 用户是修改昵称
			cookie := string(r.Header.Get("Cookie")) //[]unit8 to string
			var tk string
			if cookie == "" {
				return
			} else {
				tk = strings.Split(cookie, ":")[1]
			}

			// 解析参数
			err := r.ParseForm()
			if err != nil {
				log.Fatal("ParseForm ", err)
			}
			nickname := r.Form["nickname"][0]

			//从pool中拿连接
			conn, err := pool.Get()
			if err != nil {
				log.Printf("err:%v\n", err)
				return
			}

			//调用rpc服务
			var reqData = make(map[string]string)
			reqData["nickname"] = nickname
			reqData["token"] = tk
			rspData, err := RpcService(conn, "ChangeNickname", reqData)
			if err != nil {
				if _, ok := err.(net.Error); ok {
					// 直接关闭conn，不放回pool
					log.Println("invaild conn")
					conn.(*easypool.PoolConn).MarkUnusable()
					conn.Close()
				} else {
					//做一些登陆不通过的事
					log.Println("登陆失败")
					fmt.Fprint(w, "登陆失败")
					conn.Close()
					return
				}
			}

			// 判断状态码
			if rspData["code"] != "0" {
				log.Println("修改昵称失败")
			}

			//设置重定向
			w.Header().Set("Location", "/userhome")//跳转地址设置
			w.WriteHeader(302)
			conn.Close()
			return

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

		cookie := string(r.Header.Get("Cookie")) //[]unit8 to string
		var tk string
		if cookie == "" {
			return
		} else {
			tk = strings.Split(cookie, ":")[1]
		}

		//从pool中拿连接
		conn, err := pool.Get()
		if err != nil {
			log.Printf("err:%v\n", err)
			return
		}

		// 调用Rpc服务，验证tk，确认用户是否登陆
		var reqData = make(map[string]string)
		reqData["token"] = tk
		rspData, err := RpcService(conn, "VerifyToken", reqData)
		if err != nil {
			if _, ok := err.(net.Error); ok {
				// 直接关闭conn，不放回pool
				log.Println("invaild conn")
				conn.(*easypool.PoolConn).MarkUnusable()
				conn.Close()
			} else {
				//做一些登陆不通过的事
				log.Println("登陆失败")
				fmt.Fprint(w, "登陆失败")
				conn.Close()
				return
			}
		}

		// 判断状态码
		code := rspData["code"]
		if code != "0"{
			//未登陆
			fmt.Fprintf(w, "unauthorized")
			conn.Close()
			return
		}

		username := rspData["username"]

		format := strings.Split(handler.Filename, ".")
		filename := username + "." + format[len(format)-1]

		dirname := filepath.Join(conf.MediaPath, filename) //路径和文件名拼接
		f, err := os.OpenFile(dirname, os.O_WRONLY|os.O_CREATE, 0666) //打开目标文件等待写入，这里后期把文件名换成用户名相关的
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		conn.Close()
	}
}

// 图片显示
func showImg(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query()["name"][0]
	// 图片格式
	imgFormat := [3]string{".img", ".png", ".jpg"}
	// 判断是哪种格式
	for _, v := range imgFormat {
		path := username + v
		dirname := filepath.Join(conf.MediaPath, path) //路径和文件名拼接
		_, err := os.Stat(dirname)
		if err != nil {
			continue
		}
		w.Header().Set("Content-Type", "image")
		http.ServeFile(w, r, dirname)
		return
	}
	http.NotFound(w, r)
}

// 主页
func index(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "hello")
	//从pool中拿连接
	conn, err := pool.Get()
	if err != nil {
		log.Printf("err:%v\n", err)
		return
	}

	reqData := make(map[string]string)
	reqData["userid"] = "1"
	reqData["username"] = "123456789"
	rspData, err := RpcService(conn, "SetToken", reqData)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			// 直接关闭conn，不放回pool
			log.Println("invaild conn")
			conn.(*easypool.PoolConn).MarkUnusable()
			conn.Close()
		} else {
			//做一些登陆不通过的事
			log.Println("登陆失败")
			fmt.Fprint(w, "登陆失败")
			conn.Close()
			return
		}
	}

	// cookie设置token
	tk := rspData["token"]
	fmt.Fprint(w, tk)
	conn.Close()
}

// 打印格式
func HttpLog(method string, path string, addr string) {
	log.Println(method, " --> ", path, " By ", addr)
}

// 调用rpc的服务
func RpcService(conn net.Conn, funcname string, inArgs map[string]string) (map[string]string, error) {
	// 创客户端
	client := rpc.NewClient(conn)
	// 定义函数调用原型
	// 客户端调用rpc
	outArgs, err := client.RpcCall(funcname, inArgs)
	return outArgs, err
}