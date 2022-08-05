# 需要测试完成的工作
* 第一阶段 使用postman脚本对现有的接口进行测试，并在脚本中检验返回值
* 第二阶段 结合mongodbweb-json schema,根据test server，构建基础httpapi和flow的json
* 第三阶段 构建性能测试
* 第四阶段 根据特性 构建扩展测试用例
# postman测试

| | |
| -- | -- |
| 操作系统 | Linux |
| 运行环境 | nodejs、newman、go |

本章节主要介绍测试脚本的使用，即如何快速进行测试验证工作。

现将命令符操着过程编写为shell脚本，更方便提供给使用者进行黑盒测试。

并且提供两种测试脚本，双脚本即`apihub`程序启动和`postman`程序启动分别由两个脚本控制，方便调试人员分别调试启动。单脚本即将`apihub`程序和`postman`启动放在一个脚本控制，方便全流程自动测试。

    注：目前仅在天翼云、本地Linux环境测试运行，本地的Linux首次运行需要配置nodejs、newman、go环境。Windows环境未验证！

## 1. ！！！测前必看
postman测试过程中，需要注意postman发送地址和端口号与apihub监听地址和端口号要保持一致。

postman地址和端口号修改位置如下
```
./broker/34test_0623.postman_environment
```
*.postman_environment修改value位置地址和端口号即可。
```
{
	"key": "url",
	"value": "127.0.0.1:8080",
	"type": "default",
	"enabled": true
},
```

apihub地址和端口号修改位置如下
```
../example/main.json
```
main.json文件最下方host、port
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
## 2. 单脚本测试方式
具体shell脚本在./brker目录下，脚本名称`startup-postmanup.sh`。

脚本首先检查进程中是否存在`tms-go-apihub`应用程序，若存在则杀死进程。

然后检查当前目录下是否有可执行文件`tms-go-apihub`（apihub的可执行文件），若有则直接运行，若无则自动build可执行文件并运行。

最后等待2秒钟左右，检查当前目录下是否有`./*.postman_collection 和 *.postman_environment`文件，若有则直接运行，返回成功信息，若无打印错误信息提示用户。
## 3. 双脚本测试方式
具体shell脚本在./brker目录下，分别为start.sh、postmanup.sh

* `start.sh`首先检查进程中是否存在`tms-go-apihub`应用程序，若存在则杀死进程。检查当前目录下是否有可执行文件`tms-go-apihub`（apihub的可执行文件），若有则直接运行，若无则自动`build`可执行文件并运行。
* `postmanup.sh`检查当前目录下是否有`./*.postman_collection 和 *.postman_environment`文件，若有则直接运行，返回成功信息，若无打印错误信息提示用户。


## 4. 脚本参数修改
启动程序名称和位置或许无法适配默认shell脚本，为方便使用，打开对应shell脚本，shell头直接修改文件名和位置即可。
```
# ############################################################
# 可配置文件路径及名称：
# 
# apihub_addr：                 apihub应用程序相对位置
# conf_addr：                   json文件相对位置
# postman_collection_addr：     postman_collection文件相对位置
# postman_environment_addr：    postman_environment文件相对位置
# 
# ############################################################
apihub_addr="./tms-go-apihub"
conf_addr="../example/"
postman_collection_addr="./APIHUB_0623.postman_collection"
postman_environment_addr="./34test_0623.postman_environment"
# ############################################################
```
# 现有接口测试
将test.sh中关键内容转换为psotman脚本，满足自动化测试要求
```
curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/httpapi/amap_district"   
``` 
检查JSON均为固定值
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
检查JSON除reporttime， humidity，temperature， weather， winddirection， windpower不能为空，其余均为固定值
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
检查JSON除humidity，temperature， weather， winddirection， windpower不能为空，其余均为固定值
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
检查HTML中除humidity，temperature， weather， winddirection， windpower不能为空，其余均为固定值
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
检查JSON除sid不为空，其余均为固定值
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
检查JSON除sid不为空，其余均为固定值
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
检查JSON除sid不为空，其余均为固定值
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
检查json均为固定值
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
检查JSON除log_id不为空，其余均为固定值
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
curl -i -X POST -d '{"touser": "YangYue","msgtype": "text","agentid": "1000002", "content": "试试企业微信" }' "http://localhost:8080/flow/qywx_message_send"  
```  
检查JSON除msgid不为空，其余均为固定值
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
