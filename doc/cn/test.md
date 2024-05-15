# 需要测试完成的工作

- 第一阶段 使用 postman 脚本对现有的接口进行测试，并在脚本中检验返回值
- 第二阶段 结合 mongodbweb-json schema,根据 test server，构建基础 httpapi 和 flow 的 json
- 第三阶段 构建性能测试
- 第四阶段 根据特性 构建扩展测试用例

# postman 测试

|          |                    |
| -------- | ------------------ |
| 操作系统 | Linux              |
| 运行环境 | nodejs、newman、go |

本章节主要介绍测试脚本的使用，即如何快速进行测试验证工作。

现将命令符操着过程编写为 shell 脚本，更方便提供给使用者进行黑盒测试。

并且提供两种测试脚本，双脚本即`apihub`、`http-server`程序启动和`postman`程序启动分别由两个脚本控制，方便调试人员分别调试启动。

单脚本即将`apihub`程序、`http-server`程序、`postman`测试脚本放在一个 Shell 脚本控制，方便全流程自动测试。

    注：目前仅在天翼云、本地Linux环境测试运行，本地的Linux首次运行需要配置nodejs、newman、go环境。Windows环境未验证！
    http-server程序模拟API网关调用的外部server返回结果

## 1. ！！！测前必看

若本机 Linux 未运行或配置过 postman 脚本，可参考如下链接进行环境配置

[postman 脚本环境配置](https://blog.csdn.net/szm1234/article/details/126345866)

`postman`测试过程中，需要注意`postman`发送地址和端口号与 apihub 监听地址和端口号要保持一致。

`postman`地址和端口号修改位置如下

```
./example/34test_0623.postman_environment
```

\*.postman_environment 修改 value 位置地址和端口号即可。

```
{
	"key": "url",
	"value": "127.0.0.1:8080",
	"type": "default",
	"enabled": true
},
```

`apihub`地址和端口号修改位置如下

```
../example/main.json
```

`main.json`文件最下方`host、port`

```
{
  "name": "apiGateway",
  "command": "apiGateway",
  "description": "apiGateway",
  "args": [
    {
      "name": "host",
      "value": {
        "from": "literal",
        "content": "127.0.0.1"
      }
    },
    {
      "name": "port",
      "value": {
        "from": "literal",
        "content": "8080"
      }
    }
  ]
}
```

`http-server`程序暂时不需要修改 IP 和端口号，保持默认即可。

## 2. 单脚本测试方式

具体 shell 脚本在`./brker`目录下，脚本名称`startup-postmanup.sh`。

脚本首先检查进程中是否存在`tms-go-apihub`、`http-server`应用程序，若存在则杀死进程。

然后检查当前目录下是否有可执行文件`tms-go-apihub`（apihub 的可执行文件）、`http-server`（http-server 程序模拟 API 网关调用的外部 server 返回结果），若有则直接运行，若无则自动 build 可执行文件并运行。

最后等待 2 秒钟左右，检查当前目录下是否有`./*.postman_collection 和 *.postman_environment`文件，若有则直接运行，返回成功信息，若无打印错误信息提示用户。

## 3. 双脚本测试方式

具体 shell 脚本在./example 目录下，分别为 start.sh、postmanup.sh

- `start.sh`首先检查进程中是否存在`tms-go-apihub`、`http-server`应用程序，若存在则杀死进程。检查当前目录下是否有可执行文件`tms-go-apihub`（apihub 的可执行文件）、`http-server`，若有则直接运行，若无则自动`build`可执行文件并运行。
- `postmanup.sh`检查当前目录下是否有`./*.postman_collection 和 *.postman_environment`文件，若有则直接运行，返回成功信息，若无打印错误信息提示用户。

## 4. 脚本参数修改

启动程序名称和位置或许无法适配默认 shell 脚本，为方便使用，打开对应 shell 脚本，shell 头直接修改文件名和位置即可。

```
##############################################################
# 可配置文件路径及名称：
#
# apihub_app：                 apihub应用程序相对位置
# conf_path：                  json文件夹相对位置
#
# postman_collection_app：     postman_collection文件相对位置
# postman_environment_app：    postman_environment文件相对位置
#
# httpserver_app：             httpserver应用程序相对位置
# httpserver_ip：              httpserver应用程序默认监听IP和端口号
#
##############################################################
####################自定义位置#################################

apihub_app="../broker/tms-go-apihub"
conf_path="./"

postman_collection_app="./APIHUB_0623.postman_collection"
postman_environment_app="./34test_0623.postman_environment"

httpserver_app="../test/http-server/http-server"
httpserver_path="../test/http-server/"
httpserver_ip="127.0.0.1:6060"

##############################################################
##############################################################
```

# 现有接口测试

将 test.sh 中关键内容转换为 psotman 脚本，满足自动化测试要求

```
curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/httpapi/amap_district"
```

检查 JSON 均为固定值

```
{
    "count": "2",
    "districts": [
        {
            "adcode": "110000",
            "center": "116.407387,39.904179",
            "citycode": "010",
            "districts": [],
            "level": "province",
            "name": "北京市"
        },
        {
            "adcode": "110100",
            "center": "116.405285,39.904989",
            "citycode": "010",
            "districts": [],
            "level": "city",
            "name": "北京城区"
        }
    ],
    "info": "OK",
    "infocode": "10000",
    "status": "1",
    "suggestion": {
        "cities": [],
        "keywords": []
    }
}
```

```
curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/httpapi/amap_district/v1?appID=001"
```

同上

```json
curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/httpapi/amap_district/v1?appID=003"
```

返回`HTTP/1.1 403 Forbidden`

```
curl -i -H "Content-Type: application/json" -d '{"city": "110100"}' "http://localhost:8080/httpapi/amap_weather"
```

检查 JSON 除 reporttime， humidity，temperature， weather， winddirection， windpower 不能为空，其余均为固定值

```
{
    "count": "1",
    "info": "OK",
    "infocode": "10000",
    "lives": [
        {
            "adcode": "110100",
            "city": "北京城区",
            "humidity": "34",
            "province": "北京",
            "reporttime": "2022-06-07 13:41:20",
            "temperature": "25",
            "weather": "多云",
            "winddirection": "西",
            "windpower": "≤3"
        }
    ],
    "status": "1"
}
```

```
curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/flow/amap_city_weather"
```

检查 JSON 除 humidity，temperature， weather， winddirection， windpower 不能为空，其余均为固定值

```
{
    "data": {
        "humidity": "34",
        "region": "北京",
        "temperature": "25",
        "weather": "多云",
        "winddirection": "西",
        "windpower": "≤3"
    },
    "errCode": "1"
}
```

```
curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/flow/amap_city_weather_html"
```

检查 HTML 中除 humidity，temperature， weather， winddirection， windpower 不能为空，其余均为固定值

```
<html>
    <head>
        <title>
            Hello API
        </title>
    </head>
    <body>
        <p>
            status:1
        </p>
        <p>
            region:北京
        </p>
        <p>
            weather:多云
        </p>
        <p>
            temperature:25
        </p>
        <p>
            winddirection:西
        </p>
        <p>
            windpower:≤3
        </p>
        <p>
            humidity:34
        </p>
    </body>
</html>
```

```

curl -i -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/httpapi/kdxf_nlp_cws"
```

检查 JSON 除 sid 不为空，其余均为固定值

```
{
    "code": "0",
    "data": {
        "word": [
            "北京",
            "的",
            "天气"
        ]
    },
    "desc": "success",
    "sid": "ltp00074c51@dx241015fa33141aba00"
}
```

```
curl -i -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/httpapi/kdxf_nlp_pos"
```

检查 JSON 除 sid 不为空，其余均为固定值

```
{
    "code": "0",
    "data": {
        "pos": [
            "ns",
            "u",
            "n"
        ]
    },
    "desc": "success",
    "sid": "ltp00074cdc@dx241015fa337a1aba00"
}
```

```
curl -i -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/httpapi/kdxf_nlp_ke"
```

检查 JSON 除 sid 不为空，其余均为固定值

```
{
    "code": "0",
    "data": {
        "ke": [
            {
                "score": "0.739",
                "word": "天气"
            },
            {
                "score": "0.696",
                "word": "北京"
            }
        ]
    },
    "desc": "success",
    "sid": "ltp00074db0@dx241015fa34241aba00"
}
```

```
curl -i -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/flow/kdxf_nlp"
```

检查 json 均为固定值

```
{
    "data": {
        "pos": [
            "ns",
            "u",
            "n"
        ],
        "word": [
            "北京",
            "的",
            "天气"
        ]
    },
    "errCode": "0"
}
```

```
curl -i -X POST -d '{"content": "https://img.zcool.cn/community/01ff2059770a25a8012193a37c7695.jpg"}'  "http://localhost:8080/flow/baidu_image_classify"
```

检查 JSON 除 log_id 不为空，其余均为固定值

```
{
    "log_id": 1534054924618786784,
    "result": [
        {
            "name": "短毛猫",
            "score": "0.487635"
        },
        {
            "name": "英国短毛猫",
            "score": "0.323795"
        },
        {
            "name": "家猫",
            "score": "0.0873927"
        },
        {
            "name": "欧洲短毛猫",
            "score": "0.0221565"
        },
        {
            "name": "短毛家猫",
            "score": "0.0219908"
        },
        {
            "name": "美国短毛猫",
            "score": "0.0096925"
        }
    ]
}
```

```
curl -i -X POST -d '{"touser": "XXXXX","msgtype": "text","agentid": "1000002", "content": "试试企业微信" }' "http://localhost:8080/flow/qywx_message_send"
```

检查 JSON 除 msgid 不为空，其余均为固定值

```
{
    "errcode": 0,
    "errmsg": "ok",
    "msgid": "Dv0oBVNA9p2BIWPODPqgkhelohZZgrHQ_GN54CQh_-BJOgdwoZeHMfGeV9OVjEGFjkGQ1TptUqNmTXpFdVOD1g"
}
```

```
curl -i -H "Content-Type: application/json" -d '{"cities":["sh", "bj", "sh", "sh"], "image":"https://img.zcool.cn/community/01ff2059770a25a8012193a37c7695.jpg"}' "http://localhost:8080/schedule/amap_qywx"
```

同上
