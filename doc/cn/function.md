## 函数列表
| 调用方式     | 名称           | 入参  | 返回类型 | 用途     |
| -------------  | ----- | --------  | -------- | -------- |
| FuncMap | utc | 空 | string | 返回UTC时间,十进制秒数的字符串(UTC时间：距离1970.1.1的秒数)
| FuncMap | md5 | 任意个字符串 | string | 将输入的多个入参，按顺序连成一个新的字符串，返回其md5哈希后的字符串
| FuncMapForTemplate | utc | 空 | string | 返回UTC时间,十进制秒数的字符串(UTC时间：距离1970.1.1的秒数)
| FuncMapForTemplate | md5 | 任意个字符串 | string | 将输入的多个入参，按顺序连成一个新的字符串，返回其md5哈希后的字符串

## 函数的调用
函数都需要存入FuncMap或FuncMapForTemplate中，二者使用方法不同
### FuncMap说明
#### 定义
var FuncMap map\[string](interface{})
#### 注意
FuncMap中的函数指针必须是func ([]string) string类型。
#### 使用
##### 关键字：  "value": \{"from": "func"...}
value.from为"funcs",value.content填入函数名字;
若有入参，在args域中列出，每个参数必须是.vars.域中可以访问到的参数，且以空格间隔， 例如："args": "apikey X-CurTime X-Param"
程序会自动将args中的多个参数，转换成一个[]string类型的参数，做为入参
#### 举例
##### md5函数调用
//如下是调用md5函数，入参为.vars.apikey, .vars.X-CurTime, .vars.X-Param
  "parameters": [
    {
      "in": "header",
      "name": "X-CheckSum",
      "value": {
        "from": "func",
        "content": "md5",
        "args": "apikey X-CurTime X-Param"
      }
    },
]
##### utc函数调用
//如下是调用utc函数
    {
      "in": "header",
      "name": "X-CurTime",
      "value": {
        "from": "func",
        "content": "utc"
      }
    },
### FuncMapForTemplate说明
关于Template用法详细说明见[Template语法说明](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/template.md)
#### 定义
var FuncMapForTemplate map\[string](interface{})
#### 使用
##### 关键字： "value": \{"from": "template"...}
当value.from 为"template"时，才生效。
#### 注意
FuncMapForTemplate中的函数，入参个数不限，入参和返回值必须为string类型。
当调用FuncMapForTemplate中函数时，入参不能含有'-'字符，若含有，可以利用.vars域转换一下，见下面例子。
#### 举例
##### md5函数调用
// apis/kdxf_mlp_ke.json文件
//首先转换参数 .vars.XParam = .vars.X-Param
//然后调用md5函数，入参为三个.vars.apikey .vars.XCurTime .vars.XParam
    {
      "in": "vars",
      "name": "XParam",
      "value": {
        "from": "json",
        "content": "{{index .vars \"X-Param\"}}"
      }
    },
    {
      "in": "header",
      "name": "X-CheckSum",
      "value": {
        "from": "template",
        "content": "{{md5 .vars.apikey .vars.XCurTime .vars.XParam}}"
      }
    },
## 插件
插件主要用于注册用户私有的func，可以不编译到主程序中,仅提供.so插件。
### 位置和要求
环境变量`TGAH_PLUGIN_DEF_PATH`下的任意目录或子目录。
### 插件接口函数的编写
插件不需要在与主程序相同的环境进行编译，但必须定义接口函数:
func Register() (map\[string](interface{}), map\[string](interface{}))
无需入参;返回值为两个map, map的键为调用函数的名字，键值为函数指针。
返回值的第一个map将载入FuncMap，第二个map将载入FuncMapForTemplate供使用，使用方法见上。
#### 注意
FuncMap中的函数，函数指针必须是func ([]string) string类型。
FuncMapForTemplate中的函数，入参个数不限，入参和返回值必须为string类型。
可根据需要，将函数写入对应的map。若函数入参名字中含有'-'字符，建议存入FuncMap中
