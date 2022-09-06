# 前言
apigateway是当前apihub的主功能，用于执行灵活可编排的API功能。本文主要聚焦于HTTPAPI、FLOW、SCHEDULE的定义和执行，三者的定义和执行也是实现低代码模式的关键。

具体如何编排HTTPAPI、FLOW、SCHEDULE对应的json文件，在后续章节详细介绍，见

[JSON SCHEMA定义](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/json.md)

[API相关接口](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/apis.md)

# 定义和执行 HTTPAPI
若要实现对单个API的调用和测试，需要在`httpapi`（`按照当前目录结构，具体为./example/httpapis 目录`）目录中存放`{Id}.json`文件。

即每个 HTTPAPI 定义对应一个文件，文件名（不含扩展名`.json`）必须和 API 定义的 ID 一致。

对低代码的使用者而言，可以理解为通过apihub程序通过路由`/httpapis/{Id}`的方式调用指定的 API。例如：

```
curl -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/httpapi/amap_district"
```

    注意：
    1. curl命令：用来请求 Web 服务器，它的名字就是客户端（client）的 URL 工具的意思。
    2. 对于这个功能，只支持APPID在query，同时返回的是json格式的回应。

curl指令可做如下理解：

* -H 参数添加 HTTP 请求的标头，`"Content-Type: application/json"`格式的参数请求方式处理，即向后台发送数据的格式必须为json字符串
* -d 发送 POST 请求的数据体，HTTP 请求会自动加上标头`Content-Type : application/x-www-form-urlencoded`。并且会自动将请求转为 POST 方法，参数也可以读取本地文本文件的数据，向服务器发送。示例中`'{"city": "北京"}'`即请求城市为北京。
* `"http://localhost:8080/httpapi/amap_district"`，因为目前仍在本地测试，json文件仍在本地，即通过`localhost:8080`（apihub监听的地址和端口号）进行本地http通信，`/httpapi/amap_district`指向`/httpapis/amap_district.json`文件。

# 定义和执行 FLOW

若要实现对多个API调用，即flow的方式实现某项功能，需要在`flow`（`按照当前目录结构，具体为./example/flows 目录`）中存放`{Id}.json`文件。

每个 FLOW 定义对应一个文件，文件名（不含扩展名`.json`）必须和 FLOW 定义的 ID 一致。

apihub程序通过路由`/flow/{Id}`的方式调用指定的 FLOW。例如：

```
curl -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/flow/amap_city_weather"
```
# 定义和执行 SCHEDULE
若要实现多个API组成的计划，即需要实现多个API通过嵌套、循环等方式实现某项功能，需要在`schedule`（`按照当前目录结构，具体为./example/schedules 目录`）中存放`{Id}.json`文件。每个 SCHEDULE 定义对应一个文件，文件名（不含扩展名`.json`）必须和 SCHEDULE 定义的 ID 一致。

apihub程序通过路由`/schedule/{scheduleId}`调用指定的 FLOW。例如：

```
curl  -H "Content-Type: application/json" -d '{"cities":["sh", "bj", "sh", "sh"], "image":"https://img.zcool.cn/community/01ff2059770a25a8012193a37c7695.jpg"}' "http://localhost:8080/schedule/amap_qywx"
```
