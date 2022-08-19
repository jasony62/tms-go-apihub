
# 前言

本设计通过API编排的方式实现了对API的调用，同时API编排实现了低代码的方式实现。

本设计中低代码实现方式体现在API网关程序与API配置信息解耦。

若要实现API编排，仅需通过程序规范好的JSON文件格式进行编排即可，不涉及具体功能实现，上述过程即本设计的低代码实现方式。

本章节主要聚焦API编排过程中，在低代码方式编写JSON文件时所需的API说明，体现在API功能介绍、源码位置、输入参数说明、状态码说明，帮助使用者快速上手编排。

API列表如下：

表1：启动相关API

| API名称 | 功能简述 |
| -- | -- |
| welcome | 启动界面 |
| confValidator | json schema检查 |
| loadConf |  加载Conf文件 |
| apiGateway | API网关启动 |
| downloadConf  | 远端Conf下载 |
| decompressZip | 解压远端压缩包 |

表2：执行相关API

| API名称 | 功能简述 |
| -- | -- |
| httpApi  | HTTP请求 |
| httpResponse  | HTTP返回结果 |
| flowApi  | FLOW请求 |
| scheduleApi | SCHEDULE请求 |
| createJson  | 创建json结构体 |
| createHtml  | 创建html |


表3：辅助功能API

| API名称 | 功能简述 |
| -- | -- |
| checkStringsEqual  | 检查name、value是否相等 |
| checkStringsNotEqual  | 检查name、value是否不相等 |
| storageStore  | 存储数据 |
| storageLoad  | 加载存储数据 |
| setDefaultAccessRight  | 默认执行权限 |
| checkRight  | 检查权限 |
| fillBaseInfo | 添加基本信息 |
| logToFile | 将日志写入文件，默认目录放在log目录中（不支持Gin框架输出的日志） |

表4：普罗米修斯相关API
| API名称 | 功能简述 |
| -- | -- |
| promStart | 普罗米修斯启动 |
| promHttpCounterInc  | 普罗米修斯http统计 |

版本说明：
| 版本 | 修订人 | 说明 |
| v0.202206 | wangbinbupt |  |
| v0.20220728 | Sheng-ZM | 整理API说明，添加API使用范例 |
# 启动相关API

程序调用`API`既有外部`API`也有内部`API`，但无所谓内外，对于`command`名称调用的API都是指向的一个`API`，仅仅是地址不同。

首先对json文件中名称进行介绍，如json.md所示。

`main.json`中各启动对象均已API形式进行调用，例如：
```
{
  "name": "main",
  "description": "main for apigateway",
  "steps": [
    {
      "name": "welcome",
      "command": "welcome",
      "description": "welcome",
      "args": [
        {
          "name": "content",
          "value": {
            "from": "literal",
            "content": "welcome to use apihub"
          }
        }
      ]
    },
    ... ...
    ... ...
}
```
* `"name": "welcome"`，表示当前对象名称；

* `"command": "welcome"`，表示当前对象调用叫做`welcome`的API接口；

* `"args": [{"name": "content","value": {"from": "literal","content": "welcome to use apihub"}}`，表示将`args`关键字的`json`数组内容并发的输入到内部的`welcome API`接口。

## 1. 启动界面（welcome API）
### 1.1. 功能介绍
apihub程序启动后，首次调用conf配置文件夹时，屏幕打印输`出welcome to use apihub`字符串，用于提示用户程序开始读取conf文件夹API配置信息。
### 1.2. 位置
```
./broker/main.go/
```
### 1.3. API输入介绍
`welcome API` 输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "content" | 可选 | literal | 输出字符串 | 输出打印字符串`"welcome to use apihub"` ` |


示例：
```
{
    "name": "welcome",
    "command": "welcome",
    "description": "welcome",
    "args": [
      {
        "name": "content",
        "value": {
          "from": "literal",
          "content": "welcome to use apihub"
        }
      }
    ]
}
```
### 1.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |
## 2. schema检查（confValidator API）（完善中）
### 2.1 功能介绍
对所有需要导入的json文件进行json和json schema检查
### 2.2. 位置
```
./broker/apis/schema.go
```
### 2.3. API输入介绍
`confValidator API` 输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "schem" | 必选 | literal | Path | json schema文件夹路径，默认代码库主目录位置。即"../schema" |


示例：
```
{
    "name": "confValidator",
    "command": "confValidator",
    "description": "confValidator",
    "args": [
      {
        "name": "schema",
        "value": {
          "from": "literal",
          "content": "../schema"
        }
      }
    ]
}
```
### 2.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |
## 3. 加载Conf文件（loadConf API）
### 3.1. 功能介绍
从`--base`指定目录读取conf文件
### 3.2. 位置
```
./broker/apis/util.go
```
### 3.3. API输入介绍
`loadConf API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| 无 | 必选 | 无 | 无 | 无 |

示例：
```
{
  "name": "loadConf",
  "command": "loadConf",
  "description": "loadConf"
}
```
### 3.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |



## 4. API网关启动（apiGateway API）（完善中）
### 4.1. 功能介绍
启动`apigateway API`，注意这个api不会返回
### 4.2. 位置
```
./broker/apis/apigateway.go
```
### 4.3. API输入介绍
`apiGateway API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "host" | 可选 | literal | 0.0.0.0 ~ 255.255.255.255 | 监听地址默认0.0.0.0 |
| "port" | 可选 | literal | 1024 ~ 65535 | 监听端口默认8080 |
| "bucket“ | 可选 | literal | "true";</br>"false"; | 默认false，是否使用bucket功能 |
| "pre" | 可选 | literal | "APIGATEWAY_PRE";</br>"none";</br>"JSON名称"; | 默认_APIGATEWAY_PRE，pre flow json名字，none代表不执行 |
| "httpApi" | 可选 | literal |" _APIGATEWAY_HTTPAPI";</br>"none";</br>"JSON名称"; | 默认_APIGATEWAY_HTTPAPI，执行httpapi的flow json脚本的名字 |
| "postOK" | 可选 | literal | "_APIGATEWAY_POST_OK";</br>"none";</br>J"JSON名称"; | 默认_APIGATEWAY_POST_OK，POST OK的flow json名字，none代表不执行 |
| "postNOK" | 可选 | literal | "_APIGATEWAY_POST_NOK";</br>"none";</br>"JSON名称"; | 默认_APIGATEWAY_POST_NOK，POST NOK的flow json名字，none代表不执行 |


示例：
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
        "content": "0.0.0.0"
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
### 4.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |

## 5. 远端Conf下载（downloadConf API）（未上线）
### 5.1. 功能介绍
若本地无响应conf文件，可通过配置从远端服务器下载conf文件
### 5.2. 位置
无
### 5.3. API输入介绍
`downloadConf API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| url | 必选 | literal | Addr | 远端Conf文件地址 | 
### 5.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |

## 6. 解压远端压缩包（decompressZip）（未上线）
### 6.1. 功能介绍
若下载远端文件为压缩包，解压远端下载后的zip压缩包。
### 6.2. 位置
无
### 6.3. API输入介绍

`decompressZip API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| file | 必选 | literal | 解压文件名 | 需要解压的文件名称，通常与downloadConf这个中url中的文件名一致，均为.zip格式 | 
| password | 可选 | literal | 密码字符串 | 解压密码 | 
| path | 可选 | literal | "Path" | 默认使用`--base`目录，解压之后的存储目录 |
### 6.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |

# 执行json文件
## 1. HTTP请求（httpApi API）
### 1.1. 功能介绍
执行httpApi，发送http请求，执行一个API调用。
### 1.2. 位置
```
./broker/apis/httpapi.go
```
### 1.3. API输入介绍
`httpApi API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "name" | 必选 | literal | "httpapi文件名" | httpapi名称，详见./example/httpapis/*.json，例如"amap_district"，指向amap_district.json文件 |
| "internal" | 可选 | literal | "true";</br>"false"; | 判断是否为内部API |
| "private" | 可选 | literal | "密钥文件名" | httpapi密钥文件名称 |

示例：
```
{
   "name": "city_adcode",
   "command": "httpApi",
   "description": "查询城市的区域码",
   "args": [
     {
       "name": "name",
       "value": {
         "from": "literal",
         "content": "amap_district"
       }
     }
   ],
   "resultKey": "adcodeResult"
},
```
### 1.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |
## 2. HTTP返回结果（httpResponse API）
### 2.1. 功能介绍
在使用httpapi时，默认发送json格式的HTTP rsponse，使用flowapi和scheduleapi没有这个默认逻辑，必须调用这个API发送http response。
### 2.2. 位置
```
./broker/apis/http.go
```
### 2.3. API输入介绍
`httpResponse API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "type" | 可选 | literal | "json";</br>"html";</br>"其他"; | http response名称 |
| "key" | 可选 | literal | "NOK_result";</br>"httpapi_result";</br>"merge_result";</br>"..." | 上文"resultKey"获取，type为json时key转换为String，其他则直接按照string发送 |
| code | 可选 | template | {{ .resultKey.code }} | 上文resultKey.code获取，例如：{{.NOK_result.code}}，生成template内容， 发送的HTTP code，不输入默认200 |

示例：
```
{
  "name": "response",
  "command": "httpResponse",
  "description": "返回结果",
  "args": [
    {
      "name": "type",
      "value": {
        "from": "literal",
        "content": "json"
      }
    },
    {
      "name": "key",
      "value": {
        "from": "literal",
        "content": "NOK_result"
      }
    },
    {
      "name": "code",
      "value": {
        "from": "template",
        "content": "{{.NOK_result.code}}"
      }
    }
  ]
}
```

### 2.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |
## 3. FLOW（flowApi API）
### 3.1. 功能介绍
执行flowApi，发送flow请求，执行一个调用流程，即编排流程。
### 3.2 位置
```
./broker/core/flow.go
```
### 3.3. API输入介绍
`flowApi API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "name" | 必选 | literal | "flow文件名" | flow名称，详见./example/flows/*.json，例如"amap_city_weather_base"，指向amap_city_weather_base.json文件 |
| internal | 可选 | literal | "true";</br>"false"; | 判断是否为内部API |
| private | 可选 | literal | "密钥文件名" | flowapi密钥文件名称 |

示例：
```
{
  "name": "city_adcode",
  "command": "flowApi",
  "description": "查询城市天气",
  "args": [
    {
      "name": "name",
      "value": {
        "from": "literal",
        "content": "amap_city_weather_base"
      }
    }
  ],
  "resultKey": "weatherResult"
},
```
### 3.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |
## 4. SCHEDULE（scheduleApi API）
### 4.1. 功能介绍

执行schedule，发送schedule请求，执行一个计划流程。

### 4.2. 位置
```
./broker/core/schedule.go
```
### 4.3. API输入介绍
`scheduleApi API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "name" | 必选 | literal | "schedule文件名" | scheduleApi名称，详见./example/schedules/*.json |
| private | 可选  | literal | "密钥文件名" | schedule密钥文件名称 |

示例：
```

```
### 4.4. 状态码 
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |

## 5. 创建json结构体（createJson API）
### 5.1. 功能介绍
创建一个新的json结构体，并且存放在resultKey。
### 5.2. 位置
```
./broker/apis/json.go
```
### 5.3. API输入介绍
`createJson API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "key" | 必选 | literal | "merge_result";</br>任意; | 后文origin的name名称 |

`createJson API`输入数组`origin`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "上文content内容" | 必选 | jsonRaw | json内容 | 根据json内容生成json结构体 |

示例：
```
{
  "name": "merge_result",
  "command": "createJson",
  "description": "合并收到的结果",
  "resultKey": "merged",
  "args": [
    {
      "name": "key",
      "value": {
        "from": "literal",
        "content": "merge_result"
      }
    }
  ],
  "origin": [
    {
      "name": "merge_result",
      "value": {
        "from": "jsonRaw",
        "json": {
          "errCode": "{{.weatherResult.status}}",
          "data": {
            "region": "{{(index .weatherResult.lives 0).province}}",
            "weather": "{{(index .weatherResult.lives 0).weather}}",
            "temperature": "{{(index .weatherResult.lives 0).temperature}}",
            "winddirection": "{{(index .weatherResult.lives 0).winddirection}}",
            "windpower": "{{(index .weatherResult.lives 0).windpower}}",
            "humidity": "{{(index .weatherResult.lives 0).humidity}}"
          }
        }
      }
    }
  ]
}
```
### 5.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |


## 6. 创建html（createHtml API）
### 6.1. 功能介绍
生成html页面，并且存放在resultKey。
### 6.2. 位置
```
./broker/apis/html.go
```
### 6.3. API输入介绍
`createHtml API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "type" | 必选 | literal | "local";</br>"resource" | local表示从content中获取;</br>resource表示从resource目录获取; |
| "content" | 必选 | literal | html内容或者resource文件名 |


示例：
```
{
  "name": "merge_result",
  "command": "createHtml",
  "description": "创建html",
  "resultKey": "mergedHtml",
  "args": [
    {
      "name": "type",
      "value": {
        "from": "literal",
        "content": "local"
      }
    },
    {
      "name": "content",
      "value": {
        "from": "literal",
        "content": "<html><head><title>Hello API</title></head><body><p>status:{{.weatherResult.status}}</p><p>region:{{(index .weatherResult.lives 0).province}}</p><p>weather:{{(index .weatherResult.lives 0).weather}}</p><p>temperature:{{(index .weatherResult.lives 0).temperature}}</p><p>winddirection:{{(index .weatherResult.lives 0).winddirection}}</p><p>windpower:{{(index .weatherResult.lives 0).windpower}}</p><p>humidity:{{(index .weatherResult.lives 0).humidity}}</p></body></html>"
      }
    }
  ]
},
```
### 6.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |


# 辅助功能API
## 1. 检查name、value是否相等（checkStringsEqual API）
### 1.1. 功能介绍
检查数组中所有name和value是否都相等，用于解决200 OK + error 在response json里的问题
### 1.2. 位置
```
./broker/api/string.go
```
### 1.3. API输入介绍
`checkStringsEqual API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "0";</br>"任意" | 必选 | -- | "sendResult.errcode";</br>"任意" | 检查数组中所有name和value是否都相等，例如"name":"0"，与本json文件上文resultKey获取的sendResult.errcode对比，若sendResult.errcode与"0"相等，则API调用成功 |

示例：
```
{
  ... ...
  ... ...
  "resultKey": "sendResult"
},    
{
  "name": "checkResult",
  "description": "检查返回值",
  "command": "checkStringsEqual",
  "args": [
    {
      "name": "0",
      "value": {
        "from": "heap",
        "content": "sendResult.errcode"
      }
    }
  ]
}
```
### 1.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |

## 2. 检查name、value是否不相等（checkStringsNotEqual API）
### 2.1. 功能介绍
检查数组中所有name和value是否都不相等，用于检查http回应内的值是否有效。
### 2.2. 位置
```
./broker/api/string.go
```
### 2.3. API输入介绍
`checkStringsNotEqual API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "";</br>"任意" | 可选 | heap | "tokenResult.access_token";</br>"任意" | 检查数组中所有name和value是否不相等，例如"name":""，与本json文件上文resultKey获取的tokenResult.access_token对比，若tokenResult.access_token不为""，则API返回的token有值，即token获取成功 |

示例： 
```
{
  ... ...
  ... ...
  "args": [
    {
      ... ...
      }
    }
  ],
  "resultKey": "tokenResult"
},
{
  "name": "checkTokenResult",
  "description": "检查token是否有效",
  "command": "checkStringsNotEqual",
  "args": [
    {
      "name": "",
      "value": {
        "from": "heap",
        "content": "tokenResult.access_token"
      }
    }
  ]
}
```
### 2.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |


## 3. 存储数据（storageStore API）
### 3.1. 功能介绍
引入storage层，负责存储和加载，引入两个API storageStore，storageLoad，本小结介绍storageStore API。

多租户支持，存储某用户的数据，后面用来获取。
### 3.2. 位置
```
./broker/apis/storage.go
```
### 3.3. API输入介绍
`storageStore API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "user" | 可选 | query | "appID" | 用户ID关键字，默认为`""`，可通过外部输入，例如curl请求命令`curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/flow/amap_city_weather?appID=001"`中`appID=001`为query类型的输入 |
| "key" | 可选 | literal | "任意" | 索引key，后文origin输入参数中的name |
| "index" | 可选 | template | {{.}} | 需要存储的索引关键字，`"{{(index .adcodeResult.districts 0).adcode}}"`解析为:`(index .adcodeResult.districts 0)`返回index后面的第一个参数的某个索引对应的元素值，其余的参数为索引值，即返回上文中adcodeResult.districts索引；读取上文adcodeResult.districts索引的0索引值的adcode |
| "source" | 可选 | literal | "local";</br>"mongodb"; | 存储方式 :</br>“local”-本地结构存储;</br>"mongodb"-远程数据库存储 |
| "content" | 可选 | template | "json";</br>"...";  | 存储的内容，如果是“json”，则需要存储origin中的数据；如果为其他，则直接存储 |

`storageStore API`输入`origin`参数介绍
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "上文content内容" | 必选 | jsonRaw | json内容 | 根据json内容生成json结构体 |

示例：
```
{
  ... ...
  "description": "查询城市的区域码",
  ... ...
  "resultKey": "adcodeResult"
},
... ...
{
  "name": "storage_store",
  "command": "storageStore",
  "description": "保存查询到的城市码",
  "args": [
    {
      "name": "user",
      "value": {
        "from": "query",
        "content": "appID"
      }
    },
    {
      "name": "key",
      "value": {
        "from": "literal",
        "content": "store_result"
      }
    },
    {
      "name": "index",
      "value": {
        "from": "template",
        "content": "{{(index .adcodeResult.districts 0).adcode}}"
      }
    },
    {
      "name": "source",
      "value": {
        "from": "literal",
        "content": "local"
      }
    },
    {
      "name": "content",
      "value": {
        "from": "template",
        "content": "json"
      }
    }
  ],
  "origin": [
    {
      "name": "store_result",
      "value": {
        "from": "jsonRaw",
        "json": {
          "errCode": "{{.weatherResult.status}}",
          "data": {
            "region": "{{(index .weatherResult.lives 0).province}}",
            "weather": "{{(index .weatherResult.lives 0).weather}}",
            "temperature": "{{(index .weatherResult.lives 0).temperature}}",
            "winddirection": "{{(index .weatherResult.lives 0).winddirection}}",
            "windpower": "{{(index .weatherResult.lives 0).windpower}}",
            "humidity": "{{(index .weatherResult.lives 0).humidity}}"
          }
        }
      }
    }
  ],      
  "resultKey": "storaged"
},
```
### 4.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |

## 5. 加载存储数据（storageLoad API）
### 5.1. 功能介绍
多租户支持，读取某用户之前存储的数据，用来回复相关用户。
### 5.2. 位置
```
./broker/apis/storage.go
```
### 5.3. API输入介绍
`storageLoad API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "key" | 可选 | literal | "任意" | 索引key |
| "index" | 可选 | template | {{.}} | 需要读取的索引关键字，`"{{(index .adcodeResult.districts 0).adcode}}"`解析为:`(index .adcodeResult.districts 0)`返回index后面的第一个参数的某个索引对应的元素值，其余的参数为索引值，即返回上文中adcodeResult.districts索引；读取上文adcodeResult.districts索引的0索引值的adcode |
| "source" | 可选 | literal | "local";</br>"mongodb"; | 存储方式 :</br>“local”-本地结构存储;</br>"mongodb"-远程数据库存储 |
| "content" | 可选 | template | {{.}} | 读取的内容，如果是“json”，则需要将读取到内容解析为json；如果为其他，则直接返回 |


示例：
```
{
  ... ...
  "description": "查询城市的区域码",
  ... ...
  "resultKey": "adcodeResult"
},
... ...
{
  "name": "storage_load",
  "command": "storageLoad",
  "description": "查询城市码",
  "args": [
    {
      "name": "key",
      "value": {
        "from": "literal",
        "content": "load_result"
      }
    },
    {
      "name": "index",
      "value": {
        "from": "template",
        "content": "{{(index .adcodeResult.districts 0).adcode}}"
      }
    },
    {
      "name": "source",
      "value": {
        "from": "literal",
        "content": "local"
      }
    },
    {
      "name": "content",
      "value": {
        "from": "template",
        "content": "json"
      }
    }
  ],     
  "resultKey": "loaded"
},
```
### 5.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |

## 6. 默认执行权限（setDefaultAccessRight API）
### 6.1. 功能介绍
设置默认执行权限。
### 6.2. 位置
```
./broker/apis/right.go
```
### 6.3. API输入介绍
`setDefaultAccessRight API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "default" | 可选 | literal | "access";</br>"deny"; | 设置默认执行权限，默认“access”，如果设置为“deny”，则拒绝访问 |

示例：
```
```
### 6.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |

## 7. 检查权限（checkRight API）
### 7.1. 功能介绍
检查是否具有运行权限。
### 7.2. 位置
```
./broker/apis/right.go
```
### 7.3. API输入介绍
`checkRight API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "user" | 可选 | literal | "appID" | 查询参数中的用户id关键字 |
| "name" | 可选 | template |  {{.*.*}} | template规范的httpapi，flow，schedule名 |
| "type" | 可选 | literal | {{.*.*}} | type 是httpapi，flow，schedule |

示例：
```
{
  "name": "check_right",
  "command": "checkRight",
  "description": "查询执行权限",
  "args": [
    {
      "name": "user",
      "value": {
        "from": "query",
        "content": "appID"
      }
    },
    {
      "name": "name",
      "value": {
        "from": "template",
        "content": "{{.base.root}}"
      }
    },
    {
      "name": "type",
      "value": {
        "from": "literal",
        "content": "{{.base.type}}"
      }
    }
  ]
}
```
### 7.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |

## 8. 添加基本信息（fillBaseInfo API）
### 8.1. 功能介绍
执行api gateway pre流程，首先添加基本信息
### 8.2. 位置
```
./broker/apis/base.go
```
### 8.3. API输入介绍
`fillBaseInfo API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| user | 可选 | query | "appID" | 用户ID关键字,query类型数据，为外部输入，例如curl请求命令`curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/flow/amap_city_weather?appID=001"`中appID=001为query类型的输入 |
| uuid | 可选 | header | "uuid" | uuid请求关键字，可输入，也可默认，默认自动创建唯一标识字符串，例如`"0e67c5cf-ef93-4a49-8aad-c2dabe7ea20d"`，方便回溯调用情况 |

示例：
```
{
  "name": "fillBaseInfo",
  "command": "fillBaseInfo",
  "description": "添加基本信息",
  "args": [
    {
      "name": "user",
      "value": {
        "from": "query",
        "content": "appID"
      }
    },
    {
      "name": "uuid",
      "value": {
        "from": "header",
        "content": "uuid"
      }
    }
  ]
},
```
### 8.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |

## 9. 将日志写入文件，默认目录放在log目录中（logOutput API）

### 9.1. 功能介绍

将日志写入文件，默认目录放在log目录中，目前支持klog输出的日志，GIN框架生成的日志暂时不能写入日志文件。

### 9.2. 位置

```
./broker/apis/file.go
```

### 9.3. API输入介绍

`logToFile API`输入数组`args`参数介绍：

| 输入name       | 是否必选 | 获参位置 | value内容                                                    | 描述                                                         |
| -------------- | -------- | -------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| filepath       | 可选     | literal  | 日志写入文件的路径                                           | 默认../log/目录                                              |
| filename       | 可选     | literal  | 日志文件名称                                                 | 默认不写日志，将content字段置位空                            |
| logformat      | 可选     | literal  | 日志显示类型，“logfmt”和“json”两种格式                       | 默认为logfmt                                                 |
| loglevel       | 可选     | literal  | 屏幕及日志文件输出打印级别，一共四个级别，分别填写字符串："debug","info","warn","error" | 默认为info，内容区分大小写                                   |
| fileMaxSize    | 可选     | literal  | 每个日志文件最大值                                           | 默认50M                                                      |
| fileMaxBackups | 可选     | literal  | 一共最多存储多少个日志文件                                   | 默认100                                                      |
| maxDays        | 可选     | literal  | 日志文件默认保存多少天                                       | 默认10天                                                     |
| compress       | 可选     | literal  | 日志文件是否压缩                                             | 默认false                                                    |
| stdout         | 可选     | literal  | 日志是否输出到屏幕                                           | 默认true                                                     |
| logwithlevel   | 可选     | literal  | 日志是否分级别打印不同文件                                   | 默认为false，如果为true则，error及以上打印一个文件，全部日志打印一个文件 |

示例：

```
    {
      "name": "logOutput",
      "command": "logOutput",
      "description": "logOutput",
      "args": [
        {
          "name": "filepath",
          "value": {
            "from": "literal",
            "content": "../log/"
          }
        },
        {
          "name": "filename",
          "value": {
            "from": "literal",
            "content": ""
          }
        },
        {
          "name": "logformat",
          "value": {
            "from": "literal",
            "content": "logfmt"
          }
        },
        {
          "name": "loglevel",
          "value": {
            "from": "literal",
            "content": "warn"
          }
        },
        {
          "name": "fileMaxSize",
          "value": {
            "from": "literal",
            "content": "50"
          }
        },
        {
          "name": "fileMaxBackups",
          "value": {
            "from": "literal",
            "content": "100"
          }
        },        
        {
          "name": "maxDays",
          "value": {
            "from": "literal",
            "content": "10"
          }
        },
        {
          "name": "compress",
          "value": {
            "from": "literal",
            "content": "false"
          }
        },
        {
          "name": "stdout",
          "value": {
            "from": "literal",
            "content": "true"
          }
        },        
        {
          "name": "logwithlevel",
          "value": {
            "from": "literal",
            "content": "false"
          }
        }
      ]
    },
```

### 9.4. 状态码

| 状态码 | 描述                                    |
| ------ | --------------------------------------- |
| 200    | StatusOK，获取信息成功                  |
| 403    | StatusForbidden，获取信息失败           |
| 500    | StatusInternalServerError，获取信息失败 |

# 普罗米修斯相关API

## 1. 普罗米修斯启动（promStart API）
### 1.1. 功能介绍
启动prometheus（普罗米修斯）服务并注册`counter`和`histogram`。
### 1.2. 位置
```
./broker/apis/prometheus.go
```
### 1.3. API输入介绍
`promStart API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| host | 可选 | literal | 0.0.0.0 ~ 255.255.255.255 | prometheus默认监听地址0.0.0.0 |
| port | 可选 | literal | 1024 ~ 65535 | prometheus默认监听端口8000 |


示例：
```
{
  "name": "promStart",
  "command": "promStart",
  "description": "promStart",
  "args": [
    {
      "name": "host",
      "value": {
        "from": "literal",
        "content": "0.0.0.0"
      }
    },
    {
      "name": "port",
      "value": {
        "from": "literal",
        "content": "8000"
      }
    }
  ]
}
```
### 1.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |
## 2. 普罗米修斯http统计（promHttpCounterInc API）
### 2.1. 功能介绍
增加prometheus counter和histogram的统计。

当前在下面四个FLOW里面增加统计
* _APIGATEWAY_POST_OK.json
* _APIGATEWAY_POST_NOK.json
* _HTTPOK.json
* _HTTPNOK.json
### 2.2. 位置
```
./broker/apis/prometheus.go
```
### 2.3. API输入介绍
`promHttpCounterInc API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 获参位置 | value内容 | 描述 |
| -- | -- | -- | -- | -- |
| "httpInOut" | 可选 | literal | "httpIn";</br>"httpOut"; | 设置httpin或者httpout类型的统计，apigateway进入或出去的请求数目 |
| "duration" | 可选 | template | {{}} | http_in_duration_sec:apigateway进入的请求处理时间，0-10秒，每秒一个桶; </br>http_out_duration_sec:apigateway出去的请求处理时间，0-10秒，每秒一个桶；|
| "root";</br>"start";</br>"type";</br>"uuid";</br>"type";</br>"child";</br>"code";</br>"id";</br>"msg";</br>... ...| 可选 | literal | -- | prometheus统计的label：</br>type：apigateway入请求的类型，可以为httpapi，flow，schedule；</br>root：apigateway入请求的名称；</br>child：对外调用的httpapi的名称；</br>code：返回的HTTP回应code；|

示例：
```
{
  "name": "promHttp_In_CounterInc",
  "command": "promHttpCounterInc",
  "description": "promHttp_In_CounterInc",
  "args": [
    {
      "name": "httpInOut",
      "value": {
        "from": "literal",
        "content": "httpIn"
      }
    },
    {
      "name": "root",
      "value": {
        "from": "template",
        "content": "{{.base.root}}"
      }
    },
    {
      "name": "start",
      "value": {
        "from": "template",
        "content": "{{.base.start}}"
      }
    },
    {
      "name": "type",
      "value": {
        "from": "template",
        "content": "{{.base.type}}"
      }
    },
    {
      "name": "uuid",
      "value": {
        "from": "template",
        "content": "{{.base.uuid}}"
      }
    },
    {
      "name": "type",
      "value": {
        "from": "template",
        "content": "{{.base.type}}"
      }
    },
    {
      "name": "child",
      "value": {
        "from": "template",
        "content": "{{.stats.child}}"
      }
    },
    {
      "name": "code",
      "value": {
        "from": "template",
        "content": "{{.stats.code}}"
      }
    },
    {
      "name": "duration",
      "value": {
        "from": "template",
        "content": "{{.stats.duration}}"
      }
    },
    {
      "name": "id",
      "value": {
        "from": "template",
        "content": "{{.stats.id}}"
      }
    },
    {
      "name": "msg",
      "value": {
        "from": "template",
        "content": "{{.stats.msg}}"
      }
    }
  ]
}
```
### 2.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | StatusOK，获取信息成功 |
| 403 | StatusForbidden，获取信息失败 |
| 500 | StatusInternalServerError，获取信息失败 |
