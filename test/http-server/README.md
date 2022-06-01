
# fast-http-server
fast-http-server is a service which receives messages from http and handle it.

## how to run
cd test
cd fast-http-server
go build
./fast-http-server -addr="127.0.0.1:6060"


## examples
## 注册
curl -X GET -H 'Content-Type: application/json' -d '{"nonce":"abcd","utc":"1234","checksum":"279fc4ff795c5fb5047c27d9f23f2332"}' "http://localhost:6060/register?app=appid1"
//成功注册则返回 {token:"xxx"，expires:3600}
## 显示
curl -X GET -H "Authorization:tokenstring"'Content-Type: application/json' -d '{"content":"hello world!"}' "http://localhost:6060/echo?app=appid1
//返回: {"content":"hello world!"}
## 连接
curl -X GET -H "Authorization:tokenstring"'Content-Type: application/json' -d '{"param1":"hello","param2":"world"}' "http://localhost:6060/joint?app=appid1
//返回: {"content":"helloworld"}
## 根据空格分词
curl -X GET -H "Authorization:tokenstring"'Content-Type: application/json' -d '{"content":"hello world"}' "http://localhost:6060/split?app=appid1
//返回: {"content":"hello,world"}
