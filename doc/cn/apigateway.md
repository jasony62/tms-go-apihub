apigateway是当前apihub的主功能，用于执行灵活可编排的API功能。
## 定义和执行 HTTPAPI
需要在httpapi中存放`{Id}.json`文件。每个 HTTPAPI 定义对应一个文件，文件名（不含扩展名`.json`）必须和 API 定义的 ID 一致。
通过路由`/httpapi/{Id}`调用指定的 API。例如：

```
curl -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/httpapi/amap_district"
```

注意对于这个功能，只支持APPID在query，同时返回的是json格式的回应.
## 定义和执行 FLOW
需要在flow中存放`{Id}.json`文件。每个 FLOW 定义对应一个文件，文件名（不含扩展名`.json`）必须和 FLOW 定义的 ID 一致。
通过路由`/flow/{Id}`调用指定的 FLOW。例如：

```
curl -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/flow/amap_city_weather"
```
## 定义和执行 SCHEDULE
需要在schedule中存放`{Id}.json`文件。每个 SCHEDULE 定义对应一个文件，文件名（不含扩展名`.json`）必须和 SCHEDULE 定义的 ID 一致。

通过路由`/schedule/{scheduleId}`调用指定的 FLOW。例如：

```
curl  -H "Content-Type: application/json" -d '{"cities":["sh", "bj", "sh", "sh"], "image":"https://img.zcool.cn/community/01ff2059770a25a8012193a37c7695.jpg"}' "http://localhost:8080/schedule/amap_qywx"
```

## 值得参考的商业/开源方案
![compare](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/cn/api_compare.png)
![5](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/cn/apigateway5.png)
![1](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/cn/apigateway1.jpg)
![2](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/cn/apigateway4.png)
![6](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/cn/apigateway6.png)
![7](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/cn/apigateway7.png)

https://www.restcloud.cn/restcloud/mycms/product-gateway.html
![3](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/cn/apigateway2.jpg)
![4](https://github.com/wangbinbupt/tms-go-apihub/raw/main/doc/cn/apigateway3.jpg)

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
