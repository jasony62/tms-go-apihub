## 近期
* 继续修订readme
  - 【readme.md】中应该站在用户的角度解释一下，每个文件有什么用，是解决什么问题的
  - 【apis.md】应该提供一下示例。api的描述应该更规范些，参照一下互联网上开放api的描述。
  - 关于内容模版应该有个说明，说明一些模板的基本逻辑

* json schema
  - 根据 json 生成 schema — https://tooltt.com/json2schema/
  - 校验 jsonschema 是否合法 — https://www.jsonschemavalidator.net/
  -  校验某个 json 是否符合指定的 json schema 定义 — https://www.jsonschemavalidator.net/
  解决：
  - 部分字段没有
  - 部分字段需要设置enum
  - required： url or dynamicUrl, 
  - from为json 或者jsonRaw时required json，其他required content
  - requestContentType   可选json，form，orgin，none，或者其他任意值

* postman或者apifox的测试脚本
* 支持http请求retry，timeout
 - "timeoutPolicy": "TIME_OUT_WF",
 - "retryLogic": "FIXED",
 - "retryDelaySeconds": 600,
 - "responseTimeoutSeconds": 3600
* 根据http response content type查看是否需要分解json

## 中期
* 支持更好的http error msg返回
* 支持query数组转json数组
* http put
* 支持switch default case
* 支持load API时候，检验private信息，load FLOW时候，检验API信息，load schedule时候，检验FLOW和API
* 在JSON，HTTP处理错误时能够返回HTTP错误给调用方
* 增加plugin框架，并支持Prometheus，本地log，本地file log，基于kafka的JSON输出
* 控制流
https://yaoapps.com/doc/d.%E5%A4%84%E7%90%86%E5%99%A8/g.%E6%B5%81%E7%A8%8B%E6%8E%A7%E5%88%B6
xiang.flow.IF	IF 流程控制	查看
xiang.flow.Return	返回输入数据	查看
xiang.flow.Throw	抛出异常并结束程序	查看
* 研究gin router
https://gitee.com/easygoadmin/EasyGoAdmin_Gin_Layui/blob/master/router/router.go
* 增强内置函数
https://goframe.org/pages/viewpage.action?pageId=1114270
## 任务池
* 暴露管理API，动态日志等级
* 支持单API并发限制，令牌桶限制
* 支持对private文件秘钥加解密
* 支持在parameters中引用之前的http错误码
* 性能提升，fastjson/json-iterator等
* 支持json文件动态下载并reload（全更新）
* OpenTelemetry
 https://goframe.org/pages/viewpage.action?pageId=38575612
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
###未分类
* 重构，将api降级为httpapi，api将泛指一起
* 支持从远端http下载压缩包，解压作为conf，支持压缩包密码
* json文件load一次，反复使用
* json文件合法性检查
* 支持flow覆盖api中的private
* 多租户异步回调：本地kv存储

###schedule
* schedule支持SWITCH/loop循环命令
* schedule支持background，并发运行模式
* schedule loop支持并发调用

###httpapi
* httpapi支持token缓存
* httpapi支持动态url
* httpapi支持200OK + error code转错误码
* httpapi支持返回非json格式的http response
* httpapi支持api version
* httpapi性能提升，南向使用fasthttp
* 多租户支持,权限管理

###util
* 支持func
     * 支持无输入参数 utc
     * 支持有输入参数 md5
     * 支持template和json使用func（FuncMap）
     * 支持从.so动态注册函数

###test
* test 开发测试http server

