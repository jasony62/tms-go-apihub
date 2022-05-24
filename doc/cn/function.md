## 定义和执行 API

需要在 API 定义存放目录中存在`{apiId}.json`文件。每个 API 定义对应一个文件，文件名（不含扩展名`.json`）必须和 API 定义的 ID 一致。

需要通过环境变量`TGAH_API_DEF_PATH`指定定义文件存放位置。

通过路由`/api/{apiId}`调用指定的 API。例如：

```
curl -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/api/amap_district"
```

## 定义和执行 FLOW

需要通过环境变量`TGAH_FLOW_DEF_PATH`指定定义文件存放位置。每个 FLOW 定义对应一个文件，文件名（不含扩展名`.json`）必须和 FLOW 定义的 ID 一致。

通过路由`/flow/{flowId}`调用指定的 FLOW。例如：

```
curl -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/flow/amap_city_weather"
```
## 定义和执行 SCHEDULE

需要通过环境变量`TGAH_SCHEDULE_DEF_PATH`指定定义文件存放位置。每个 SCHEDULE 定义对应一个文件，文件名（不含扩展名`.json`）必须和 SCHEDULE 定义的 ID 一致。

通过路由`/schedule/{scheduleId}`调用指定的 FLOW。例如：

```
curl  -H "Content-Type: application/json" -d '{"cities":["sh", "bj", "sh", "sh"], "image":"https://img.zcool.cn/community/01ff2059770a25a8012193a37c7695.jpg"}' "http://localhost:8080/schedule/amap_qywx"
```

## 插件
插件主要用于注册用户私有的func，可以不编译到主程序中。

程序启动会导入环境变量`TGAH_PLUGIN_DEF_PATH`制定的目录下以及子目录下所有的.so插件，将插件提供的函数载入hub.FuncMap和hub.FuncMapForTemplate
插件不需要在与主程序相同的环境进行编译，但.so插件需要定义接口函数:
func Register() (map[string](interface{}), map[string](interface{}))，其中第一个map将指定载入hub.FuncMap，第二个map将载入hub.FuncMapForTemplate使用

