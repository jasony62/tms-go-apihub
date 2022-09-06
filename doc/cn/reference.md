# 低代码平台
基于go的低代码平台
https://gitee.com/easygoadmin
https://gitee.com/easygoadmin/EasyGoAdmin_Gin_Layui


项目地址：https://github.com/YaoApp/yao
官方文档：https://yaoapps.com/doc
https://yaoapps.com/doc/d.%E5%A4%84%E7%90%86%E5%99%A8/c.%E7%BD%91%E7%BB%9C%E8%AF%B7%E6%B1%82


# API网关
基于https://zhuanlan.zhihu.com/p/464504713
|-|Kong|Traefik|Ambassador|Tyk|Zuul|apihub|
| -- | -- | -- | -- | -- |-- |-- |
| 基本 |  |  |  |  | ||
|主要用途|企业级API管理	|微服务网关	|微服务网关	|微服务网关	|微服务网关|企业级API管理|
|学习曲线|适中|simple|simple|适中|simple|simple|
|成本|开源/企业版|开源|开源/pro|开源/企业版|开源|开源|
|社区star|20742|21194|1719|4299|7186|5|
| 配置 |  |  |  |  | ||
|配置语言|	Admin Rest api, Text file(nginx.conf 等)	|TOML|	YAML(kubernetes annotation)	|Tyk REST API|	REST API，YAML静态配置|JSON|
|配置端点类型|	命令式	|声明式|声明式|命令式|命令式|声明式|
|拖拽配置|yes|no|no|no|no|二期|
|管理模式|	configurable|	decentralised, self-service	|decentralised, self-service	|decentralised, self-service	|decentralised, self-service|decentralised, self-service|
|可扩展性 |  |  |  |  | ||
|扩展功能|插件|自己实现|插件|插件|自己实现|二期|
|扩展方法|水平|水平|水平|水平|水平|N/A|
| 功能 |  |  |  |  | ||
|多API编排|no|no|no|no|no|yes|
|API并发调用|no|no|no|no|no|yes|
|服务发现|动态|动态|动态|动态|动态|二期|
|协议|http,https,websocket|	http,https,grpc,websocket|	http,https,grpc,websocket|	http,https,grpc,websocket|	http,https|http,http|
|基于|	kong+nginx|	traefik|envoy|tyk|zuul|自研|
|ssl终止|yes|yes|yes|yes|no|二期|
|websocket|yes|yes|yes|yes|no|二期|
|routing|host,path,method|host,path|host,path,header|host,path|path|
|限流|	yes	|no|yes|yes|需要开发|二期|
|熔断|	yes|yes|no|yes|需要其他组件|二期|
|重试|	yes|yes|no|yes|	yes|二期|
|健康检查|yes|no|no|yes|yes|二期|
|负载均衡算法|轮询，哈希|轮询，加权轮询	加权轮询|轮询|轮询，随机，加权轮询，自定义|二期|
|权限|	Basic Auth, HMAC, JWT, Key, LDAP, OAuth 2.0, PASETO, plus paid Kong Enterprise options like OpenID Connect|	basic|	yes|	HMAC，JWT，Mutual TLS，OpenID Connect，基本身份验证，LDAP，社交OAuth（例如GPlus，Twitter，Github）和传统基本身份验证提供程序	|开发实现|二期|
|tracing|yes|	yes|	yes|	yes|	需要其他组件|no|
|istio集成|	no|	no|	yes|no|	no|no|
|dashboard|	yes|yes|grafana,Prometheus|yes|no|grafana,Prometheus|
| 部署 |  |  |  |  | ||
|kubernetes	适中(k8s yaml,helm chart)|	easy|	easy|	适中(k8s yaml,helm chart)|适中(k8s yaml,helm chart)|适中(k8s yaml,helm chart)|
|Cloud IAAS|high|easy|N/A|easy|easy|N/A|
|Private Data Center|high|easy|N/A|easy|easy|二期|
|部署模式|金丝雀(企业版)|金丝雀|金丝雀，shadow|金丝雀|金丝雀|N/A|
|state|postgres,cassandra|kubernetes|kubernetes	redis|内存文件|内存|
![compare](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/pic/api_compare.png)
![5](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/pic/apigateway5.png)
![1](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/pic/apigateway1.jpg)
![2](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/pic/apigateway4.png)
![6](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/pic/apigateway6.png)
![7](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/pic/apigateway7.png)

https://www.restcloud.cn/restcloud/mycms/product-gateway.html
![3](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/pic/apigateway2.jpg)
![4](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/pic/apigateway3.jpg)

https://github.com/luraproject/lura
![a](https://www.krakend.io/images/KrakendFlow.png)

多个api网关比较
https://blog.csdn.net/pushiqiang/article/details/95726137

私有云开源解决方案：
Netflix Zuul，zuul是spring cloud的一个推荐组件，https://github.com/Netflix/zuul
Kong kong是基于Nginx+Lua进行二次开发的方案， https://konghq.com/
Tyk是2014年创建的开源API网关，甚至比AWS的API网关即服务功能还要早。Tyk用Golang编写并使用Golang自己的HTTP服务器。

公有云解决方案：
Amazon API Gateway，https://aws.amazon.com/cn/api-gateway/
阿里云API网关，https://www.aliyun.com/product/apigateway/
腾讯云API网关， https://cloud.tencent.com/product/apigateway

# 标准
[OpenAPI Specification](https://swagger.io/specification/)

https://netflix.github.io/conductor/configuration/workflowdef/
https://states-language.net/spec.html
