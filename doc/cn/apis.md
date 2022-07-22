
# 启动相关API（main.json）

程序调用`API`既有外部`API`也有内部`API`，但无所谓内外，对于`command`名称调用的API都是指向的一个`API`，仅仅是地址不同。

首先对json文件中名称进行介绍，如表所示。

| JSON名称  | 数据类型 | 描述 |
| -- | -- | -- |
| name |  String | 对象名称 |
| command |  String | 当前对象调用API名称 |
| description | String | 当前对象描述信息 |
| steps |  String | 多个对象顺序执行 |
| args |  Array | 多个并列对象 |
| value |  String | 当前对象输入value值 |
| from |  String | 当前对象默认literal |
| content |  String | 当前对象输入内容 |

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
    ...
    ...
}
```
* `"name": "welcome"`，表示当前对象名称；

* `"command": "welcome"`，表示当前对象调用叫做`welcome`的API接口；

* `"args": [{"name": "content","value": {"from": "literal","content": "welcome to use apihub"}}`，表示将`args`关键字的`json`数组内容并发的输入到内部的`welcome API`接口。

## 1. welcome API
### 1.1 功能介绍
apihub程序启动后，首次调用conf配置文件夹时，屏幕打印输出welcome to use apihub字符串，用于提示用户程序开始读取conf文件夹API配置信息。
### 1.2 位置
```
./broker/main.go
```
### 1.3 API输入介绍
`welcome API`输入内容`args`介绍：
| 参数名称 | 是否必选 | 数据类型 | 内容描述 |
| -- | -- | -- | -- |
| name | 必选 | String | 输入名称`"content"` |
| content | 可选 | String | 输入打印字符串`"welcome to use apihub"` |

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
### 1.4 状态码
| 状态码 | 描述 |
| -- | -- |
| 200 | 获取信息成功。 |
## 2. confValidator API（完善中）
对所有需要导入的json文件进行json和json schema检查

`confValidator` API输入参数介绍：
| 参数名称 | 是否必选 | 数据类型 | 内容描述 |
| -- | -- | -- | -- |
| schema | 必选 | String | json schema文件夹路径|

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
## 本地配置加载（loadConf）

loadConf输入参数介绍：
| 参数名称 | 是否必选 | 数据类型 | 描述 |
| -- | -- | -- | -- |
| loadConf | 必选 | 无 | 从base目录读取conf文件 |

```
{
  "name": "loadConf",
  "command": "loadConf",
  "description": "loadConf"
}
```
## 普罗米修斯加载（promStart）

启动prometheus服务并注册counter和histogram

promStart输入参数介绍：
| 参数名称 | 是否必选 | 数据类型 | 描述 |
| -- | -- | -- | -- |
| promHost | 可选 | String | 监听地址默认0.0.0.0 |
| promPort | 可选 | String | 监听端口默认8000 |

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
## APIG加载（apiGateway）（完善中）

启动apigateway，注意这个api不会返回

apiGateway输入参数介绍：
| 参数名称 | 是否必选 | 数据类型 | 描述 |
| -- | -- | -- | -- |
| host | 可选 | String | 监听地址默认0.0.0.0 |
| port | 可选 | String | 监听端口默认8080 |
| bucket | 可选 | Bool | 默认false，是否使用bucket功能 |
| pre | 可选 |  String | 默认_APIGATEWAY_PRE，pre flow json名字，none代表不执行 |
| httpApi | 可选 | String | 默认_APIGATEWAY_HTTPAPI，执行httpapi的flow json名字 |
| postOK | 可选 | String | 默认_APIGATEWAY_POST_OK，POST OK的flow json名字，none代表不执行 |
| postNOK | 可选 |String  | 默认_APIGATEWAY_POST_NOK，POST NOK的flow json名字，none代表不执行 |
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
## 远端Conf下载（downloadConf）（未上线）

若本地无响应conf文件，可通过配置从远端服务器下载conf文件

downloadConf输入参数介绍：
| 参数名称 | 是否必选 | 数据类型 | 描述 |
| -- | -- | -- | -- |
| url | 必选 | String | 远端文件地址 | 


## 解压远端压缩包（decompressZip）（未上线）

若下载远端文件为压缩包，解压远端下载后的zip压缩包。

decompressZip输入参数介绍：

| 参数名称 | 是否必选 | 数据类型 | 描述 |
| -- | -- | -- | -- |
| file | 必选 | String | 需要解压的文件名称，通常与downloadConf这个中url中的文件名一致，均为.zip格式 | 
| password | 可选 | String | 解压密码 | 
| path | 可选 | String | 默认使用`--base`目录，解压之后的存储目录 |



# 执行json文件

## httpApi
执行httpapi，发送http请求。

httpApi具体

httpApi参数：
| 参数名称 | 是否必选 | 数据类型 | 描述 |
| -- | -- | -- | -- |
| name | 必选 | String | httpapi名称 |
| private | 可选 | String | httpapi密钥文件名称 |

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
}
```
## flowApi


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

