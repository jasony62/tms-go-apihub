# 基础结构

value和param是两个基础结构，用在多处。

value定义为:

| 字段名称 | 是否必选 | 数据类型 | 描述 | 
| -- | -- | -- | -- | 
| from |  必选 | String | 获取参数值的位置,支持:</br>`literal`  直接从content里获取;</br>`query` http query);</br>`header` http header;</br>`private` 从秘钥文件读取;</br>`origin` 原始报文body中的json;</br>`env` 系统env;</br>`heap` 从原始报文和处理结果获取;</br>`json` 根据json生成字符串;</br>`jsonRaw` 根据json生成json结构体;</br>`template` 从content中根据模板生成;</br>`func` hub.FuncMap内部定义函数的名称。 |
| content | 可选 | String | 参数名称，或者函数名称，或者template的内容。|
| args | 可选 | String | from为func时，func的输入参数，多个参数时需要以空格分割，如："args": "apikey X-CurTime X-Param"。 |
| json | 可选 | String | json的输入值,支持`.origin.`访问输入json，.vars.访问在parameters定义的值，支持采用template的FuncMap的方式直接调用hub.FuncMapForTemplate内部定义的函数(例如"template": "{{md5 .vars.apikey .vars.XCurTime .vars.XParam}}")。如果入参名字含有字符-，则需要定义一个新的vars，去掉原名字中的-。 |

param定义为:

| 字段名称 | 是否必选 | 数据类型 | 描述 | 
| -- | -- | -- | -- |
| name | 必选 | String |  变量名称    |
| value | 必选 | Object | 值，value结构体 |
# PRIVATE
private用于存储秘钥类信息，和API相对独立，不同部署，值会发生变化。

| 字段名称 | 是否必选 | 数据类型 | 描述 |  
| -- | -- | -- | -- |
| name | 必选 | String | private名称称。|                                                  
| value | 必选 | String | private值。|
# HTTPAPI
建议所有需要从template，func中传入的参数，都独立定义在parameter中，并使用origin。

| 字段名称 | 是否必选 | 数据类型 | 描述 | 
| -- | -- | -- | -- |
| id | 必选 | String | HTTPAPI，而是根据指定 定义的标识。 |
| url | 可选 | String | HTTPAPI，而是根据指定 的目标地址。不包括任何查询参数。 |
| dynamicUrl | 可选 | Object |  当url为空时，必须提供这个结构，用来动态生成URL（比如路径中含有appId），结构为标准的value结构。 |
| private | 可选 | String | HTTPAPI，而是根据指定 秘钥文件名。| 
| description | 可选 | String | HTTPAPI，而是根据指定 的描述。 |
| method | 必选 | String | HTTP 请求方法，支持`POST`和`GET`。 |
| requestContentType | 必选 | String | json映射为`application/json`，form映射为`application/x-www-form-urlencoded`，origin为取输入报文的ContentType，并直接转发输入报文的http body，none表示没有body,其他值则直接写入ContentType|
| args | 可选 | Object[] |  HTTP 请求的参数。 |
| &nbsp; &nbsp; &nbsp; &nbsp;-- in | 必选 | String | 参数位置。支持:</br>`query`;</br>`header`;</br>`body`;</br> `vars`。</br>前三者的值除了会放到发送报文里，还可以在模板通过.vars.访问，vars表示只进入.vars|
| &nbsp; &nbsp; &nbsp; &nbsp;-- name | 必选 | String | 参数名称。 | 
| &nbsp; &nbsp; &nbsp; &nbsp;-- value | 必选 | Object | 参数值，标准value结构。 |
| cache | 可选 | Object | HTTP请求是否支持缓存模式，如果支持，在过期时间内，将不会再向服务器请求，而是直接返回缓存内容。 |
| &nbsp; &nbsp; &nbsp; &nbsp; -- format | 必选 | String | 指定过期时间的解析格式。分为秒“second”和具体时间格式，如：“20060102150405” |
| &nbsp; &nbsp; &nbsp; &nbsp;-- expire | 必选 | Object | 指定过期时间的获取位置，标准value结构。 |
| &nbsp; &nbsp; &nbsp; &nbsp;-- from | 必选 | String | 差异：获取过期时间的位置，是从header域中获取的话，则设置为“header”，如果从body中获取，则设置为“template” |

目前系统并未使用`id`字段定位选择的 HTTPAPI，而是根据指定 HTTPAPI 定义文件的名称。

# API
用于FLOW和SCHEDULE中
| 字段名称 | 是否必选 | 数据类型 | 描述 |  
| -- | -- | -- | -- |
| name | 必选 | String | API的名称。|
| description | 可选 | String | API的描述。| 
| command | 必选 | String | API名称。|
| private | 可选 | String | 可以用于计算value和覆盖api内部的private。|
| resultKey | 可选 | String | 执行结果保存时的索引名称，origin,vars,result,loop为保留值不可使用。      |
| args | 可选 | Object[] | api的输入参数,为param结构体|
| origin | 可选 | Object[] | 进行tempalte替换时，origin数据，为param结构体|

# FLOW
| 字段名称 | 是否必选 | 数据类型 | 描述 |  
| -- | -- | -- | -- |
| name | 必选 | String | FLOW的名称。|
| description | 可选 | String | FLOW的描述。| 
| private | 可选 | String | API 秘钥文件名用于覆盖内层。 |
| steps  | 必选 | Object[] | 串行调用API的步骤。为API结构体。   |

# SCHEDULE
| 字段名称 | 是否必选 | 数据类型 | 描述 |  
| -- | -- | -- | -- |
| name | 必选 | String | SCHEDULE 定义的标识。 |
| description | 可选 | String | SCHEDULE 的描述。|
| concurrentNum | 可选 | Int | 最大允许的并行执行的数量。 |
| steps | -- | Object[] | schedule任务列表。 |
| &nbsp; &nbsp; &nbsp; &nbsp;-- type | 必选 | String | `api`;</br>`loop`;</br>`switch`| 
| &nbsp; &nbsp; &nbsp; &nbsp;-- mode | 可选 | String | 执行模式:</br>`normal`;</br>`concurrent`;</br>`background` |
| &nbsp; &nbsp; &nbsp; &nbsp;-- private | 可选 | String | API 秘钥文件名用于覆盖内层。   | 
|&nbsp; &nbsp; &nbsp; &nbsp;-- api | 可选 | Object | API结构体，type为api时执行。 |
|&nbsp; &nbsp; &nbsp; &nbsp;-- control | 可选 | Object | control结构体，type为loop和switch时执行。 |

control结构体定义为：
| 字段名称 | 是否必选 | 数据类型 | 描述 |  
| -- | -- | -- | -- |
| name | 必选 | String | FLOW的名称。|
| description | 可选 | String | FLOW的描述。|
| private | 可选 | String | API 秘钥文件名用于覆盖内层。 |
| &nbsp; &nbsp; &nbsp; &nbsp;-- resultKey | 可选 | String |  在API或者FLOW 执行结果对应的名称，在loop时将索引保存在.loop.resultKey,便于后续引用(如{{index .origin.cities .loop.myloop}}), origin,vars,result,loop为保留值不可使用。 |
| key | 可选 | Object |  switch时为要检查的值，loop时为循环的次数，标准from结构。 |
| concurrentNum | 可选 | Int |  最大允许的并行执行的数量。 |
| concurrentLoopNum | 可选 | Int |  最大允许的loop内并行执行的数量。 |
| steps | -- | object[] | schedule任务列表。 | 
|cases | 可选 | Object[] | switch时检查的case。 | 
| &nbsp; &nbsp; &nbsp; &nbsp;-- value | 必选 | String | 上层的key等于本字段则执行tasks。 |
| &nbsp; &nbsp; &nbsp; &nbsp;-- concurrentNum | 可选 | Int |  最大允许的并行执行的数量。 | 
| &nbsp; &nbsp; &nbsp; &nbsp;-- steps | 可选 | Object[] | 结构同上层的tasks，为tasks的自身嵌套。 | 
|}|
# RIGHT
| 字段名称 | 是否必选 | 数据类型 | 描述 |  
| -- | -- | -- | -- |
| type | 必选 | String | 权限文件对应的执行类型，httpapi，flow，schedule | 
| right | 必选 | String | 权限类型：</br>`public`（所有人都允许调用）;</br>`internal`（只允许内部调用，不允许外部调用）;</br>`whitelist`（只有list中的才允许访问）;</br>`blacklist`（非list中的才允许访问） |
| list | 必选 | Object[] | `user list`数组 |

