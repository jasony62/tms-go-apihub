通过编排能力简化对 API 的调用。

# 启动

## 环境变量

| 环境变量           | 用途                           | 默认值  |
| ------------------ | ------------------------------ | ------- |
| TGAH_HOST          | 服务的主机名                   | 0.0.0.0 |
| TGAH_PORT          | 服务的端口号                   | 8080    |
| TGAH_BUCKET_ENABLE | API 和 FLOW 是否按 bucket 隔离 | no      |
| TGAH_API_DEF_PATH  | API 定义文件存放位置           | -       |
| TGAH_FLOW_DEF_PATH | 编排定义文件存放位置           | -       |

## 命令行

通过`--env`指定使用的环境变量文件。

```
run go . --env envfile
```

```
run build -o tms-gah-broker
```

```
./tms-gah-broker --env envfile
```

## docker

```
docker build -t tms/gah-broker .
```

```
docker run -it --rm --name tms-gah-broker -p 8080:8080 tms/gah-broker sh
```

```
docker compose build tms-gah-broker
```

```
docker compose up tms-gah-broker
```

## 安装插件

将插件代码复制到容器中

```
docker cp plugins tms-gah-broker:/home/tms-gah/plugins
```

进入容器编译插件

```
docker exec -it tms-go-broker sh
```

```
cd plugins
cd kdxfnlp
go build -buildmode=plugin -o kdxfnlp.so kdxfnlp.go
```

# 基础

## 定义 API

| 字段          | 用途                                                                                                  | 类型     | 必选 |
| ------------- | ----------------------------------------------------------------------------------------------------- | -------- | ---- |
| id            | API 定义的标识。                                                                                      | string   | 是   |
| url           | API 的目标地址。不包括任何查询参数。                                                                  | string   | 是   |
| private       | API 秘钥文件名。                                                                                      | string   | 是   |
| description   | API 的描述。                                                                                          | string   |      |
| method        | HTTP 请求方法，支持`POST`和`GET`。                                                                    | string   | 是   |
|               |                                                                                                       |          |      |
| parameters    | HTTP 请求的参数。                                                                                     | string[] |      |
| --in          | 参数位置。支持`query`和`header`。                                                                     | string   | 是   |
| --name        | 参数名称。                                                                                            | string   | 是   |
| --value       | 参数的值。                                                                                            | string   | 否   |
| --from        | 指定参数值的获取位置。                                                                                | object   | 否   |
| ----in        | 获取参数值的位置,支持`query`和`private`。                                                             |          |      |
| ----name      | 参数值所在位置的名称。                                                                                |          |      |
|               |                                                                                                       |          |      |
| requestBody   | 发送给被调用方的内容。如果不指定，直接转发调用方发送的内容。                                          | any      | 否   |
| --contentType | 默认为`application/json`，还支持`"application/x-www-form-urlencoded`。                                | string   | 否   |
| --content     | 发送给目标 API 的内容模板（mustache）。                                                               |          |      |
|               |                                                                                                       |          |      |
| response      | 返回给调用方的内容。返回的内容统一为`application/json`格式。如果不指定，直接转发目标 API 返回的内容。 | object   | 否   |
| --json        | 返回调用方内容的模板（mustache），数组或对象。支持从被调用方返回的结果进行映射。                      | any      | 是   |
|               |                                                                                                       |          |      |
| plugins       | 支持通过`plugin`对 API 定义进行改写。                                                                 | object[] | 否   |
| --path        | 插件文件的路径。                                                                                      | string   | 是   |

目前系统并未使用`id`字段定位选择的 API，而是根据指定 API 定义文件的名称。

## 编排 API

| 字段           | 用途                                                                                                                                     | 类型     | 必选 |
| -------------- | ---------------------------------------------------------------------------------------------------------------------------------------- | -------- | ---- |
| name           | API 调用流的名称。                                                                                                                       | string   | 是   |
| description    | API 调用流的描述。                                                                                                                       | string   | 否   |
| steps          | 调用流执行的步骤。每个步骤对应 1 个 API 的调用。API 必须是已定义。                                                                       | object[] | 是   |
| --name         | 步骤的名称。                                                                                                                             | string   | 是   |
| --description  | 步骤的描述。                                                                                                                             | string   | 是   |
| --resultKey    | 在上下文中 API 执行结果对应的名称。                                                                                                      | string   | 是   |
| --api          | 步骤对应的 API 定义。                                                                                                                    | object   | 是   |
| ----id         | API 定义的 ID。                                                                                                                          | string   | 是   |
| ----parameters | API 的参数定义，这里可以覆盖 API 定义中的参数定义。`from.in`可以指定为`StepResult`，代表从之前执行步骤的结果（和 resultKey）中提取数据。 | object[] | 否   |
|                |                                                                                                                                          |          |      |
| --response     | 定义返回结果的模板。                                                                                                                     | object   | 否   |
| ----json       | 统一返回 JSON 格式的内容。                                                                                                               | any      | 是   |

# 功能

## 定义和执行 API

需要在 API 定义存放目录中存在`{apiId}.json`文件。每个 API 定义对应一个文件，文件名（不含扩展名`.json`）必须和 API 定义的 ID 一致。

需要通过环境变量`TGAH_API_DEF_PATH`指定定义文件存放位置。

通过路由`/api/{apiId}`调用指定的 API。例如：

```
curl "http://localhost:8080/api/amap_district?city=北京"
```

## 定义和执行调用流 FLOW

需要通过环境变量`TGAH_FLOW_DEF_PATH`指定定义文件存放位置。每个 FLOW 定义对应一个文件，文件名（不含扩展名`.json`）必须和 API 定义的 ID 一致。

通过路由`/flow/{flowId}`调用指定的 FLOW。例如：

```
curl "http://localhost:8080/flow/amap_city_weather?city=北京"
```

## 数据转换模板

待补充

# 插件

插件需要在与主程序相同的环境进行编译。

# 隔离

使用`bucket`进行数据隔离。

# 参考

[OpenAPI Specification](https://swagger.io/specification/)

https://netflix.github.io/conductor/configuration/workflowdef/
