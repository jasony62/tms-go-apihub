## 功能
该应用将swagger yaml/json文件，转换为apihub使用的json文件

## 注释
主要使用第三方库"github.com/neotoolkit/openapi"，该库是基于Openapi 3.0.3版本的swagger文件结构进行转换。


## 命令行
通过`--from`指定读取swagger文件的目录(默认为./from/)，通过`--to`指定生成apihub conf json文件的文件夹路径，默认为./to/。

## 暂不支持：
* 涉及swagger2.0的文件转换
* server URL使用的path，暂时不支持变量替换，也就是RESTful接口
* response格式定制、转换
* oauth2鉴权
·