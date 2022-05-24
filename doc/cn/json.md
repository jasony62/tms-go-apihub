
## FROM
from是一个基础结构，用在多处
| 字段           | 用途                                                                                                                                     | 类型     | 必选 |
| -------------- | ---------------------------------------------------------------------------------------------------------------------------------------- | -------- | ---- |
| from      | 获取参数值的位置,支持`literal`(直接从content里获取)，`query`,`header`,`private`(从秘钥文件读取),`origin`(原始报文body中的json),StepResult(从原始报文和处理结果获取)，json(从json中生成内容)，template(从content中生成),`func`(hub.FuncMap内部定义函数的名称)。               |     string     |  是    |
| content      | 参数名称，或者函数名称，或者template的内容。                 |     string     | 否     |
| args      | from为func时，func的输入参数，多个参数时需要以空格分割，如："args": "apikey X-CurTime X-Param"                 | string         |   否   |
| json  | json的输入值,支持.origin.访问输入json，.vars.访问在parameters定义的值，支持采用template的FuncMap的方式直接调用hub.FuncMapForTemplate内部定义的函数(例如"template": "{{md5 .vars.apikey .vars.XCurTime .vars.XParam}}")。如果入参名字含有字符-，则需要定义一个新的vars，去掉原名字中的-|    object      |  否    |
## PRIVATE
| 字段           | 用途                                                                                                                                     | 类型     | 必选 |
| -------------- | ---------------------------------------------------------------------------------------------------------------------------------------- | -------- | ---- |
| name           | private名称称。                                                                                                                       | string   | 是   |
| value    | private值。                                                                                                                       | string   | 否   |

## API
建议所有需要从template，func中传入的参数，都独立定义在parameter中，并使用origin。

| 字段          | 用途                                                                                                  | 类型     | 必选 |
| ------------- | ----------------------------------------------------------------------------------------------------- | -------- | ---- |
| id            | API 定义的标识。                                                                                      | string   | 是   |
| url           | API 的目标地址。不包括任何查询参数。                                                                  | string   | 否   |
| dynamicUrl    | 当url为空时，必须提供这个结构，用来动态生成URL（比如路径中含有appId），结构为标准的from。不包括任何查询参数。    | object   | 否   |
| private       | API 秘钥文件名。                                                                                      | string   | 否   |
| description   | API 的描述。                                                                                          | string   |  否    |
| method        | HTTP 请求方法，支持`POST`和`GET`。                                                                    | string   | 是   |
|requestContentType | json映射为`application/json`，form映射为`application/x-www-form-urlencoded`，origin为取输入报文的ContentType，并直接转发输入报文的http body，none表示没有body,其他值则直接写入ContentType|string |是|
|               |                                                                                                       |          |      |
| parameters    | HTTP 请求的参数。                                                                                     | object[] |    否  |
| --in          | 参数位置。支持`query`，`header`,`body`, `vars`。前三者的值除了会放到发送报文里，还可以在模板通过.vars.访问，vars表示只进入.vars| string| 是   |
| --name        | 参数名称。                                                                                            | string   | 是   |
| --from        | 指定参数值的获取位置，标准from结构。          | object   | 否   |
| response      | 返回给调用方的内容。返回的内容统一为`application/json`格式。如果不指定，直接转发目标 API 返回的内容。 | object   | 否   |
| --json        | 返回调用方内容的模板（mustache），数组或对象。支持从被调用方返回的结果进行映射。                      | any      | 是   |
|               |                                                                                                |          |      |
| cache | HTTP请求是否支持缓存模式，如果支持，在过期时间内，将不会再向服务器请求，而是直接返回缓存内容。 | object | 否 |
| --format | 指定过期时间的解析格式。分为秒“second”和具体时间格式，如：“20060102150405” | string | 是 |
| --from | 指定过期时间的获取位置，标准from结构。 | object | 是 |
| ----from | 差异：获取过期时间的位置，是从header域中获取的话，则设置为“header”，如果从body中获取，则设置为“template” | string | 是 |
|          |                                        |        |    |
| respStatus | 指定回应body体中的状态码 | object | 否 |
| --from | 指定状态码的获取位置，标准from结构。 | object | 是 |
| --format | 指定状态码的解析格式。分为数字“number”和string | string | 是 |
| --expected | 状态码的期望正确值 | string | 是 |


目前系统并未使用`id`字段定位选择的 API，而是根据指定 API 定义文件的名称。

## FLOW

| 字段           | 用途                                                                                                                                     | 类型     | 必选 |
| -------------- | ---------------------------------------------------------------------------------------------------------------------------------------- | -------- | ---- |
| name           | FLOW的名称。| string   | 是   |
| description    | FLOW的描述。| string   | 否   |
| concurrentNum           | 最大允许的并行执行的数量。| int   | 否   |
| steps          | 调用API的步骤。每个步骤对应 1 个 API 的调用。API 必须是已定义。   | object[] | 是   |
| --name         | 步骤的名称。          | string   | 是   |
| --description  | 步骤的描述。       | string   | 否   |
| --resultKey    | 在上下文中 API 执行结果对应的名称，origin,vars,result,loop为保留值不可使用。      | string   | 否   |
| --concurrent   | 是否使用并行执行。                       | bool   | 否   |
| --api          | 步骤对应的 API 定义。                                    | object   | 是   |
| ----private    | 用于覆盖原始API中的private。                            | string   | 否   |
| ----id         | API 定义的 ID。                            | string   | 是   |
| ----parameters | 放在这里的定义会补充或者覆盖输入报文里的json参数。`from.from`可以指定为`StepResult`，代表从之前执行步骤的结果（和 resultKey）中提取数据。 | object[] | 否   |
|                |          |      ||
| ------name        | 参数名称。                          | string   | 是   |
| ------from        | 指定参数值的获取位，标准from结构置。                               | object   | 否   |
|               |              |          |      |
| --response     | 定义返回结果的模板。                                      | object   | 否   |
| ----type   | 返回什么格式的内容，json或者html。                                       | string      | 是   |
| ----from | 返回的内容或者template定义，标准from结构 | object | 否 |

## SCHEDULE
| 字段          | 用途                                                                                                  | 类型     | 必选 |
| ------------- | ----------------------------------------------------------------------------------------------------- | -------- | ---- |
| name            | SCHEDULE 定义的标识。                                         | string   | 是   |
| description   | SCHEDULE 的描述。                        | string   |  否    |
| concurrentNum           | 最大允许的并行执行的数量。                               | int   | 否   |
|          |      |
| tasks    | 任务列表。                                object[] |      |
| --type          | control，api, flow| string   | 是   |
| --name          | type为api, flow时，为对应的名字，control时可以是switch，loop| string   | 是   |
|--description   | task 的描述。      | string   |  否    |
|--resultKey   |  在API或者FLOW 执行结果对应的名称，在loop时将索引保存在.loop.resultKey,便于后续引用(如{{index .origin.cities .loop.myloop}}), origin,vars,result,loop为保留值不可使用。      | string   |  否    |
|--key   |  switch时为要检查的值，loop时为循环的次数，标准from结构。      | object   |  否    |
|--concurrentNum   |  control命令时最大允许的并行执行的数量。      | int   |  否    |
| --concurrent   | 是否使用并行执行。                       | bool   | 否   |
| --tasks   | control命令时的执行列表，结构同上层的tasks，为tasks的自身嵌套。                       | object[]   | 否   |
| --parameters   | flow和control时用于改写origin，同flow的parameters。        | object[]   | 否   |
| --private    | 用于覆盖原始API中的private。                            | string   | 否   | 
| --cases   | switch时检查的case。                       | object[]   | 否   |
| ----Value   | 上层的key等于本字段则执行tasks。                       | string   | 是   |
| ----concurrentNum   |  control命令时最大允许的并行执行的数量。      | int   |  否    |
| --tasks   | 结构同上层的tasks，为tasks的自身嵌套。                       | object[]   | 否   |