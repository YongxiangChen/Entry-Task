package conf

import "path/filepath"



const RpcAddr = ":8008" //rpc服务器地址

var BasePathHttp, _ = filepath.Abs(".") //项目根目录
var StaticPath = filepath.Join(BasePathHttp, "httpserver/static/") //static目录
var MediaPath = filepath.Join(BasePathHttp, "httpserver/media/") //media目录
