## 近期
* 支持返回非json格式的http response
* json schema
* 支持api带数组的
* 开发测试http server，postman或者apifox的测试脚本
* pre/post func支持
* 本地kv存储

## 中期
* api version
* 对输入数组的支持
* 支持switch default case
* 支持load API时候，检验private信息，load FLOW时候，检验API信息，load schedule时候，检验FLOW和API
* 支持在http response中访问origin中的值
* 在JSON，HTTP处理错误时能够返回HTTP错误给调用方
* 支持http请求retry，timeout
 **"timeoutPolicy": "TIME_OUT_WF",
 **"retryLogic": "FIXED",
 **"retryDelaySeconds": 600,
 **"responseTimeoutSeconds": 3600
* 增加plugin框架，并支持Prometheus，本地log，本地file log，基于kafka的JSON输出
## 任务池
* 暴露管理API，动态日志等级
* 支持单API并发限制，令牌桶限制
* 支持对private文件秘钥加解密
* 支持在parameters中引用之前的http错误码
* 性能提升，使用fasthttp，fastjson/json-iterator等
* 支持json文件动态下载并reload（全更新）
* go async pool
* graceful shutdown
* 参数有效性检查
* 支持异步，循环加异步
* 支持基于user的private
## 需要考虑
* Opentracing，Skywalking
* 多SSL证书
* 熔断，降级
* API健康检查
* 支持API调用websocket，gRPC，Dubbo，redis，kafka
* Open API ：支持使用open api配置网关
* URL Scheme
* 是否考虑集成进
    * https://github.com/go-kratos/kratos
    * https://github.com/eolinker/apinto

## 已经完成：
* 支持从远端http下载压缩包，解压作为conf，支持压缩包密码
* json文件load一次，反复使用
* json文件合法性检查
* 支持SWITCH/loop循环命令
* 支持loop并发调用
* 支持flow，schedule，loop并发调用
* 支持func获取value
    * 支持无输入参数 utc 
    * 支持有输入参数 md5
    * 支持template使用func（FuncMap）
    * 支持从.so动态注册函数
* 支持token缓存
* 支持flow覆盖api中的private
* 支持在url中使用private
* 支持200OK + error code转错误码