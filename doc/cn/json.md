# 基础结构
value和param是两个基础结构，用在多处。
value定义为
| 字段           | 用途  | 类型     | 必选 |
| ------------- | ----- | -------- | ---- |
| from      | 获取参数值的位置,支持`literal`(直接从content里获取)，`query`(http query),`header`(http header),`private`(从秘钥文件读取),`origin`(原始报文body中的json),"env"(系统env)，StepResult(从原始报文和处理结果获取)，json(根据json生成字符串)，jsonRaw(根据json生成json结构体)，template(从content中生成),`func`(hub.FuncMap内部定义函数的名称)。               |     string     |  是    |
| content      | 参数名称，或者函数名称，或者template的内容。                 |     string     | 否     |
| args      | from为func时，func的输入参数，多个参数时需要以空格分割，如："args": "apikey X-CurTime X-Param"                 | string         |   否   |
| json  | json的输入值,支持.origin.访问输入json，.vars.访问在parameters定义的值，支持采用template的FuncMap的方式直接调用hub.FuncMapForTemplate内部定义的函数(例如"template": "{{md5 .vars.apikey .vars.XCurTime .vars.XParam}}")。如果入参名字含有字符-，则需要定义一个新的vars，去掉原名字中的-|    object      |  否    |
param定义为
| 字段           | 用途  | 类型     | 必选 |
| ------------- | ----- | -------- | ---- |
| name      | 变量名称               |     string     |  是    |
| value    | 值，value结构体。    | object   | 是   |
# PRIVATE
private用于存储秘钥类信息，和API相对独立，不同部署，值会发生变化。
| 字段           | 用途| 类型     | 必选 |
| -------------- | ------------- | -------- | ---- |
| name           | private名称称。                                                                                                                       | string   | 是   |
| value    | private值。                                                                                                                       | string   | 是   |
# HTTPAPI
建议所有需要从template，func中传入的参数，都独立定义在parameter中，并使用origin。

| 字段          | 用途 | 类型     | 必选 |
| -------------- | ------------- | ----- | ---- |
| id            | HTTPAPI，而是根据指定 定义的标识。  | string   | 是   |
| url           | HTTPAPI，而是根据指定 的目标地址。不包括任何查询参数。 | string   | 否   |
| dynamicUrl    | 当url为空时，必须提供这个结构，用来动态生成URL（比如路径中含有appId），结构为标准的value结构。    | object   | 否   |
| private       | HTTPAPI，而是根据指定 秘钥文件名。| string   | 否   |
| description   | HTTPAPI，而是根据指定 的描述。 | string   |  否    |
| method        | HTTP 请求方法，支持`POST`和`GET`。                                                                    | string   | 是   |
|requestContentType | json映射为`application/json`，form映射为`application/x-www-form-urlencoded`，origin为取输入报文的ContentType，并直接转发输入报文的http body，none表示没有body,其他值则直接写入ContentType|string |是|
|  |  |  |  |
| args    | HTTP 请求的参数。                                                                                     | object[] |    否  |
| --in          | 参数位置。支持`query`，`header`,`body`, `vars`。前三者的值除了会放到发送报文里，还可以在模板通过.vars.访问，vars表示只进入.vars| string| 是   |
| --name        | 参数名称。 | string   | 是   |
| --value        | 参数值，标准value结构。          | object   | 是   |  |  |  |  |
| cache | HTTP请求是否支持缓存模式，如果支持，在过期时间内，将不会再向服务器请求，而是直接返回缓存内容。 | object | 否 |
| --format | 指定过期时间的解析格式。分为秒“second”和具体时间格式，如：“20060102150405” | string | 是 |
| --expire | 指定过期时间的获取位置，标准value结构。 | object | 是 |
| ----from | 差异：获取过期时间的位置，是从header域中获取的话，则设置为“header”，如果从body中获取，则设置为“template” | string | 是 |
|          |                                        |        |    |

目前系统并未使用`id`字段定位选择的 HTTPAPI，而是根据指定 HTTPAPI 定义文件的名称。

# API
用于FLOW和SCHEDULE中
| 字段           | 用途  | 类型     | 必选 |
| ------------- | ----- | -------- | ---- |
| name           | API的名称。| string   | 是   |
| description    | API的描述。| string   | 否   |
| command    | API名称。| string   | 是   |
| private    | 可以用于计算value和覆盖api内部的private。| string   | 否   |
|resultKey    | 执行结果保存时的索引名称，origin,vars,result,loop为保留值不可使用。      | string   | 否   |
| args  | api的输入参数,为param结构体|    object[]      |  否    |
| origin  | 进行tempalte替换时，origin数据，为param结构体|    object[]      |  否    |

# FLOW
| 字段           | 用途  | 类型     | 必选 |
| ------------- | ----- | -------- | ---- |
| name           | FLOW的名称。| string   | 是   |
| description    | FLOW的描述。| string   | 否   |
| private       | API 秘钥文件名用于覆盖内层。 | string   | 否   |
| steps          | 串行调用API的步骤。为API结构体。   | object[] | 是   |

# SCHEDULE
| 字段           | 用途  | 类型     | 必选 |
| ------------- | ----- | -------- | ---- |
| name            | SCHEDULE 定义的标识。                                         | string   | 是   |
| description   | SCHEDULE 的描述。                        | string   |  否    |
| concurrentNum           | 最大允许的并行执行的数量。                               | int   | 否   |
|          |      |
| steps    | schedule任务列表。                                object[] |      |
| --type          | api, loop, switch| string   | 是   |
| --mode          | 执行模式，normal，concurrent，background三种   | 否   |
| --private          | API 秘钥文件名用于覆盖内层。   | 否   |
|--api   | API结构体，type为api时执行。      | object   |  否    |
|--control   | control结构体，type为loop和switch时执行。      | object   |  否    |

control结构体定义为
| 字段           | 用途  | 类型     | 必选 |
| ------------- | ----- | -------- | ---- |
| name           | FLOW的名称。| string   | 是   |
| description    | FLOW的描述。| string   | 否   |
| private       | API 秘钥文件名用于覆盖内层。 | string   | 否   |
|--resultKey   |  在API或者FLOW 执行结果对应的名称，在loop时将索引保存在.loop.resultKey,便于后续引用(如{{index .origin.cities .loop.myloop}}), origin,vars,result,loop为保留值不可使用。      | string   |  否    |
|key   |  switch时为要检查的值，loop时为循环的次数，标准from结构。      | object   |  否    |
|concurrentNum   |  最大允许的并行执行的数量。      | int   |  否    |
|concurrentLoopNum   |  最大允许的loop内并行执行的数量。      | int   |  否    |
| steps    | schedule任务列表。                                object[] |      |
|          |      |
|cases   | switch时检查的case。                       | object[]   | 否   |
| --value   | 上层的key等于本字段则执行tasks。                       | string   | 是   |
| --concurrentNum   |  最大允许的并行执行的数量。      | int   |  否    |
| --steps   | 结构同上层的tasks，为tasks的自身嵌套。                       | object[]   | 否   |

# RIGHT
| 字段    | 用途                                                         | 类型     | 必选 |
| ------- | ------------------------------------------------------------ | -------- | ---- |
| type    | 权限文件对应的执行类型，httpapi，flow，schedule              | string   | 是   |
| right   | 权限类型：public（所有人都允许调用），internal（只允许内部调用，不允许外部调用），whitelist（只有list中的才允许访问），blacklist（非list中的才允许访问） | string   | 是   |
| list    | user list数组                                                | object[] | 是   |
