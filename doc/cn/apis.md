支持api列表
| 名称           | 入参  | 用途     |
| ------------- | ----- | -------- |
| 启动相关（用于main.json） |  |  |
| loadConf | url（可选） 远端文件地址，password（可选）下载文件的密码 | 当不提供url时候，直接从base目录读取，否则下载解压后再读取 |
| confValidator（完善中） | schema（必选） json schema路径 | 对所有json文件进行json和json schema检查 |
| apiGateway | host（可选，默认0.0.0.0）监听地址，port（可选，默认8080），监听端口， bucket（可选，默认false），是否使用bucket功能| 启动apigateway，注意这个api不会返回 |
| 执行json文件 |  |  |
| httpApi | name httpapi名字，private（可选）秘钥文件名称| 执行httpapi，发送http请求 |
| flowApi | name flow名字，private（可选）秘钥文件名称| 执行flow |
| scheduleApi | name schedule名字，private（可选）秘钥文件名称| 执行schedule |
| http相关 |  |  |
| httpResponse | type（json，html，或者其他） http response名称，key，从哪个result获取，type为json时转换为string，其他则直接按照string发送| 在使用httpapi时，默认发送json格式的HTTP rsponse，flow和schedule没有这个默认逻辑，必须调用这个API发送http response|
| 辅助类 |  |  |
| checkStringsEqual | 任意| 检查数组中所有name和value是否都相等，用于解决200 OK+ error 在response json里的问题，参考例子{"name": "0","value": {"from": "StepResult","content": "sendResult.errcode"}} |
| checkStringsNotEqual | 任意| 检查数组中所有name和value是否都不相等，用于检查http回应内的值是否有效，参考例子        {"name": "","value": {"from": "StepResult","content": "tokenResult.access_token"}} |
| createJson | key origin入参中的name| 创建一个新的json结构体，并且存放在resultKey |
| createHtml | type（local则从content中获取，resource则从resource目录获取），content html内容或者resource文件名| 生成html页面，并且存放在resultKey |
|setDefaultAccessRight| default deny, 没有right配置文件拒绝访问；access， 没有right配置文件允许访问(默认)| 检查是否具有运行权限 |
| checkRight | userKey 查询参数中的用户id关键字； name httpapi名字；type 是httpapi，flow，schedule | 检查是否具有运行权限 |
| storageStore | user 查询用户appID，可以配置在query中，也可以在header中；key origin入参中的name；index 需要存储的索引关键字； source 存储方式 “local”-本地结构存储；content 存储的内容， 如果是“json”，则需要存储origin中的数据，如果为其他，则直接存储 | 多租户支持，存储某用户的数据，后面用来获取 |
| storageLoad | index 需要读取的索引关键字； source 存储方式 “local”-本地结构存储；content 读取的内容， 如果是“json”，则需要将读取到内容解析为json，如果为其他，则直接返回 | 多租户支持，读取某用户之前存储的数据，用来回复相关用户 |

