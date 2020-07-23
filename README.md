# Entry Task



## 一、接口文档

### RPC服务端接口

1. ``func NewServer(addr string) *Server {}``

   - 接口功能：
     新建server结构体新实例

   - 参数说明：

     | 参数 | 类型   | 说明         |
     | ---- | ------ | ------------ |
     | addr | string | 服务器的地址 |

   - 返回值说明：
     返回server结构体指针

   - 调用实例
     ``server = NewServer(":8080")``

2. ``func (s *Server) Register(name string, fn interface{}) {}``

   - 接口功能：
     注册RPC服务端的方法，供客户端使用

   - 参数说明：

     | 参数 | 类型        | 说明       |
     | ---- | ----------- | ---------- |
     | name | string      | 方法名     |
     | fn   | interface{} | 要注册方法 |

   - 返回值说明：
     无

   - 调用方式：
     ``s.Register("Authenticate",  Authenticate)``

3. ``func (s *Server) Run() error {}``

   - 接口功能：
     运行RPC服务器
   - 参数说明：
     无
   - 返回值说明：
     返回tcp创建连接的错误
   - 调用方式：
     ``s.Run()``

### RPC客户端接口

1. ``func NewClient(addr string) *Client {}``

   - 接口功能：
     新建Client结构体新实例

   - 参数说明：

     | 参数 | 类型   | 说明         |
     | ---- | ------ | ------------ |
     | addr | string | 服务器的地址 |

   - 返回值说明：
     返回client结构体指针

   - 调用实例
     ``client = NewClient(":8080")``

2. ``func (client *Client) RpcCall(name string, fpoint interface{}) {}``

   - 接口功能：
     调用RPC已注册的服务

   - 参数说明：

     | 参数   | 类型        | 说明             |
     | ------ | ----------- | ---------------- |
     | name   | string      | 服务（方法）名   |
     | fpoint | interface{} | 客户端函数的指针 |

   - 返回值说明：
     无

   - 调用示例：

     ```go
     // 先声明一个和服务端的函数参数和返回值都一样的函数
     var auth = func(name string, pw string) bool
     // 传入时注意传入指针
     client.RpcCall("Authenticate", &auth)
     // 调用
     ok := auth("123456789", "123456789")
     ```

### RPC服务端可供调用的服务

1. ``func Authenticate(username string, password string) (model.User, bool) {}``

   - 接口功能：
     验证用户名和密码

   - 参数说明：

     | 参数     | 类型   | 说明     |
     | -------- | ------ | -------- |
     | name     | string | 用户名   |
     | password | string | 用户密码 |

   - 返回值说明：

     | 类型       | 说明           |
     | ---------- | -------------- |
     | model.User | User结构体对象 |
     | bool       | 验证通过与否   |

2. ``func SetToken(user model.User) (string, error) {}``

   - 接口功能：
     根据user生存token，存入redis并返回

   - 参数说明：

     | 参数 | 类型       | 说明           |
     | ---- | ---------- | -------------- |
     | user | model.User | User结构体对象 |

   - 返回值说明：

     | 类型   | 说明  |
     | ------ | ----- |
     | string | token |
     | error  | 错误  |

3. ``func VerifyToken(tk string) (model.User, int) {}``




## 二、性能测试报告

### 未优化版本

<img src="/Users/yongxiangchen/Library/Application Support/typora-user-images/image-20200723101532914.png" alt="image-20200723101532914" style="zoom:50%;" />

- 分析：并发量为200，模拟200个不同用户，每个用户访问2次，QPS=200

<img src="/Users/yongxiangchen/Library/Application Support/typora-user-images/image-20200723101410291.png" alt="image-20200723101410291" style="zoom:50%;" />

- 分析：增加每个用户访问量到10次，QPS=363

### 增加数据库索引，使用sql的连接池

<img src="/Users/yongxiangchen/Library/Application Support/typora-user-images/image-20200723105520607.png" alt="image-20200723105520607" style="zoom:50%;" />

- 分析：并发量为200，模拟200个不同用户，每个用户访问10次，QPS=444

### 

