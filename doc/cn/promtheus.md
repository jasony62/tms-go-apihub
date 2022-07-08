当前是在下面四个FLOW里面增加统计
* _APIGATEWAY_POST_OK.json
* _APIGATEWAY_POST_NOK.json
* _HTTPOK.json
* _HTTPNOK.json

# 参数
| 字段           | 类型|解释  |
| ------------- | -----|----- |
|http_in|counter|apigateway进入的请求数目|
|http_in_duration_sec|histogram|apigateway进入的请求处理时间，0-10秒，每秒一个桶|
|http_out|counter|httpapi发出的请求数目|
|http_out_duration_sec|histogram|httpapi发出的请求处理时间，0-10秒，每秒一个桶|
# label
| 名称           | 解释  |
| ------------- | ----- |
|type|apigateway请求的类型，可以为httpapi，flow，schedule|
|root|apigateway请求的名称|
|child|当前执行的名称|
|code|返回的HTTP回应code|
