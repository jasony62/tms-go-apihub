支持api列表
| 名称           | 入参  | 用途     |
| ------------- | ----- | -------- | 
| 启动相关 |  |  | 
| loadConf | url（可选） 远端文件地址，password（可选）下载文件的密码 | 当不提供url时候，直接从base目录读取，否则下载解压后再读取 | 
| apiGateway | host（可选，默认0.0.0.0）监听地址，port（可选，默认8080），监听端口， bucket（可选，默认false），是否使用bucket功能| 启动apigateway，注意这个api不会返回 | 
| 执行json文件 |  |  | 
| httpApi | name httpapi名字，private（可选）秘钥文件名称| 执行httpapi，发送http请求 | 
| flowApi | name flow名字，private（可选）秘钥文件名称| 执行flow | 
| scheduleApi | name schedule名字，private（可选）秘钥文件名称| 执行schedule | 
| http相关 |  |  | 
| httpResponse | type（json，html，或者其他） http response名称，key，从那个result获取，type为json时转换为string，其他则直接按照string发送| 回应http response | 
| 辅助类 |  |  | 
| checkStringsEqual | 任意| 检查数组中所有name和value是否都相等 | 
| checkStringsNotEqual | 任意| 检查数组中所有name和value是否都不相等 | 
| createJson | key origin入参中的name| 创建一个新的json结构体，并且存放在resultKey | 
| createHtml | type（local则从content中获取，resource则从resource目录获取），content html内容或者resource文件名| 生成html页面，并且存放在resultKey | 
