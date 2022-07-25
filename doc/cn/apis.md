
# 启动相关API（main.json）

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
./broker/main.go
```
### 1.3. API输入介绍
`welcome API` 输入数组`args`参数介绍：
| 输入name | 是否必选 | 数据类型 | content内容 | 描述 |
| -- | -- | -- | -- | -- |
| content | 可选 | String | 输出字符串 | 输出打印字符串`"welcome to use apihub"` ` |


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
| 200 | 获取信息成功。 |
## 2. schema检查（confValidator API）（完善中）
### 2.1 功能介绍
对所有需要导入的json文件进行json和json schema检查
### 2.2. 位置
```
./broker/apis/schema.go
```
### 2.3. API输入介绍
`confValidator API` 输入数组`args`参数介绍：
| 输入name | 是否必选 | 数据类型 | content内容 | 描述 |
| -- | -- | -- | -- | -- |
| schem | 必选 | String  | Path | json schema文件夹路径，默认代码库主目录位置。即"../schema" |


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
| 200 | 获取信息成功。 |
## 3. 加载Conf文件（loadConf API）
### 3.1. 功能介绍
从`--base`指定目录读取conf文件
### 3.2. 位置
```
./broker/apis/util.go
```
### 3.3. API输入介绍
`loadConf API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 数据类型 | content内容 | 描述 |
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
| 200 | 获取信息成功。 |

## 4. 普罗米修细启动（promStart API）
### 4.1. 功能介绍
启动prometheus（普罗米修斯）服务并注册`counter`和`histogram`。
### 4.2. 位置
```
./broker/apis/prometheus.go
```
### 4.3. API输入介绍
`promStart API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 数据类型 | content内容 | 描述 |
| -- | -- | -- | -- | -- | 
| promHost | 可选 | String | 0.0.0.0 ~ 255.255.255.255 | prometheus默认监听地址0.0.0.0 |
| promPort | 可选 | String | 1024 ~ 65535 | prometheus默认监听端口8000 |


示例：
```
{
  "name": "promStart",
  "command": "promStart",
  "description": "promStart",
  "args": [
    {
      "name": "promHost",
      "value": {
        "from": "literal",
        "content": "0.0.0.0"
      }
    },
    {
      "name": "promPort",
      "value": {
        "from": "literal",
        "content": "8000"
      }
    }
  ]
}
```
### 4.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | 获取信息成功。 |

## 5. API网关启动（apiGateway API）（完善中）
### 5.1. 功能介绍
启动`apigateway API`，注意这个api不会返回
### 5.2. 位置
```
./broker/apis/apigateway.go
```
### 5.3. API输入介绍
`apiGateway API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 数据类型 | content内容 | 描述 |
| -- | -- | -- | -- | -- |
| host | 可选 | String | 0.0.0.0 ~ 255.255.255.255 | 监听地址默认0.0.0.0 |
| port | 可选 | String | 1024 ~ 65535 | 监听端口默认8080 |
| bucket | 可选 | Bool | "ture";</br>"false" | 默认false，是否使用bucket功能 |
| pre | 可选 |  String | "APIGATEWAY_PRE";</br>"none";</br>JSON名称 | 默认_APIGATEWAY_PRE，pre flow json名字，none代表不执行 |
| httpApi | 可选 | String |" _APIGATEWAY_HTTPAPI";</br>"none";</br>JSON名称 | 默认_APIGATEWAY_HTTPAPI，执行httpapi的flow json脚本的名字 |
| postOK | 可选 | String | "_APIGATEWAY_POST_OK";</br>"none";</br>JSON名称 | 默认_APIGATEWAY_POST_OK，POST OK的flow json名字，none代表不执行 |
| postNOK | 可选 |String | "_APIGATEWAY_POST_NOK";</br>"none";</br>JSON名称 | 默认_APIGATEWAY_POST_NOK，POST NOK的flow json名字，none代表不执行 |


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
### 5.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | 获取信息成功。 |

## 6. 远端Conf下载（downloadConf API）（未上线）
### 6.1. 功能介绍
若本地无响应conf文件，可通过配置从远端服务器下载conf文件
### 6.2. 位置
无
### 6.3. API输入介绍
`downloadConf API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 数据类型 | content内容 | 描述 |
| -- | -- | -- | -- | -- |
| url | 必选 | String | Addr | 远端Conf文件地址 | 
### 6.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | 获取信息成功。 |

## 7. 解压远端压缩包（decompressZip）（未上线）
### 7.1. 功能介绍
若下载远端文件为压缩包，解压远端下载后的zip压缩包。
### 7.2. 位置
无
### 7.3. API输入介绍

`decompressZip API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 数据类型 | content内容 | 描述 |
| -- | -- | -- | -- | -- |
| file | 必选 | String | 解压文件名 | 需要解压的文件名称，通常与downloadConf这个中url中的文件名一致，均为.zip格式 | 
| password | 可选 | String | 密码字符串 | 解压密码 | 
| path | 可选 | String | Path | 默认使用`--base`目录，解压之后的存储目录 |
### 7.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | 获取信息成功。 |

# 执行json文件(conf文件)
## 1. HTTP请求（httpApi API）
### 1.1. 功能介绍
执行httpApi，发送http请求，执行一个API调用。
### 1.2. 位置
```
./broker/apis/httpapi.go
```
### 1.3. API输入介绍
`httpApi API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 数据类型 | content内容 | 描述 |
| -- | -- | -- | -- | -- |
| name | 必选 | String | JSON文件名 | httpapi名称，详见./example/http/apis/*.json，例如"amap_district"，指向amap_district.json文件 |
| internal | 可选 | Bool | "ture";</br>"false" | 判断是否为内部API |
| private | 可选 | String | 密码字符串 | httpapi密钥文件名称 |

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
| 200 | 获取信息成功。|
## 2. FLOW（flowApi API）
### 2.1. 功能介绍
执行flowApi，发送flow请求，执行一个调用流程，即编排流程。
### 2.2 位置
### 2.3. API输入介绍
`flowApi API`输入数组`args`参数介绍：
| 输入name | 是否必选 | 数据类型 | content内容 | 描述 |
| -- | -- | -- | -- | -- |
| name | 必选 | String | JSON文件名 | httpapi名称，详见./example/http/apis/*.json，例如"amap_district"，指向amap_district.json文件 |
| internal | 可选 | Bool | "ture";</br>"false" | 判断是否为内部API |
| private | 可选 | String | 密码字符串 | httpapi密钥文件名称 |

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
}
```
### 2.4. 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | 获取信息成功。|

| 名称| 入参  | 用途|
| -- | -- | -- |
| 启动相关（用于main.json） |  |  |
| decompressZip | file（必选） 需要解压的文件名称，通常与downloadConf这个中url中的文件名一致，均为.zip格式；<br />password（可选） 解压密码<br />path（可选） 解压之后的存储目录，不配则使用base目录 | 解压远端下载后的zip压缩包 |
| 执行json文件 |  |  |
| httpApi | name httpapi名字，private（可选）秘钥文件名称| 执行httpapi，发送http请求 |
| flowApi | name flow名字，private（可选）秘钥文件名称| 执行flow |
| scheduleApi | name schedule名字，private（可选）秘钥文件名称| 执行schedule |
| http相关 |  |  |
| httpResponse | type（json，html，或者其他） http response名称，key，从哪个result获取，type为json时转换为string，其他则直接按照string发送 code, 发送的HTTP code，不输入则为200| 在使用httpapi时，默认发送json格式的HTTP rsponse，flow和schedule没有这个默认逻辑，必须调用这个API发送http response|
| 辅助类 |  |  |
| checkStringsEqual | 任意| 检查数组中所有name和value是否都相等，用于解决200 OK+ error 在response json里的问题，参考例子{"name": "0","value": {"from": "heap","content": "sendResult.errcode"}} |
| checkStringsNotEqual | 任意| 检查数组中所有name和value是否都不相等，用于检查http回应内的值是否有效，参考例子        {"name": "","value": {"from": "heap","content": "tokenResult.access_token"}} |
| createJson | key origin入参中的name| 创建一个新的json结构体，并且存放在resultKey |
| createHtml | type（local则从content中获取，resource则从resource目录获取），content html内容或者resource文件名| 生成html页面，并且存放在resultKey |
| setDefaultAccessRight | default（可选） 设置默认执行权限，默认“access”，如果设置为“deny”，则拒绝访问 | 默认执行权限 |
| checkRight | userKey 查询参数中的用户id关键字； name httpapi名字；type 是httpapi，flow，schedule | 检查是否具有运行权限 |
| storageStore | user 查询用户appID，可以配置在query中，也可以在header中；key origin入参中的name；index 需要存储的索引关键字； source 存储方式 “local”-本地结构存储；content 存储的内容， 如果是“json”，则需要存储origin中的数据，如果为其他，则直接存储 | 多租户支持，存储某用户的数据，后面用来获取 |
| storageLoad | index 需要读取的索引关键字； source 存储方式 “local”-本地结构存储；content 读取的内容， 如果是“json”，则需要将读取到内容解析为json，如果为其他，则直接返回 | 多租户支持，读取某用户之前存储的数据，用来回复相关用户 |
| promStart | host（可选，默认0.0.0.0）监听地址，port（可选，默认8000），监听端口| 启动prometheus服务并注册counter和histogram |
| promHttpCounterInc | httpInOut（填"httpIn"或者"httpOut"）设置httpin或者httpout类型的统计，其他参数均是供prometheus统计的label | 增加prometheus counter和histogram的统计 |

