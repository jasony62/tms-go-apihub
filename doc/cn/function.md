## 插件
插件主要用于注册用户私有的func，可以不编译到主程序中。

程序启动会导入环境变量`TGAH_PLUGIN_DEF_PATH`制定的目录下以及子目录下所有的.so插件，将插件提供的函数载入hub.FuncMap和hub.FuncMapForTemplate
插件不需要在与主程序相同的环境进行编译，但.so插件需要定义接口函数:
func Register() (map[string](interface{}), map[string](interface{}))，其中第一个map将指定载入hub.FuncMap，第二个map将载入hub.FuncMapForTemplate使用

