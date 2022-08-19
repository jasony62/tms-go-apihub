APIHUB是一款基于 Golang 开发的API调度平台，能够实现基于JSON定义的灵活的编排能力。

第一阶段主要是提供微服务网关功能，充分利用了Go协程的高并发性能。   

未来可以对接低代码平台，大大简化对API调用的管理。

# [基本概念](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/base.md)
介绍本平台HTTPAPI、FLOW、SCHEDULE的基本概念以及它们之间的关系。
# [快速启动和编译](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/start.md)
介绍源码程序编译和启动流程，快速实现程序黑盒运行的验证工作。
# [API定义和执行](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/apigateway.md)
介绍apihub的主功能apigateway，以及如何快速定义和执行HTTPAPI、FLOW、SCHEDULE方式json文件。
# [JSON SCHEMA定义](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/json.md)
介绍JSON SCHEMA定义，字段名称与描述，及相关json格式定义。
# [Template语法说明](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/template.md)
介绍json文件中使用到的template模板。
# [API相关接口](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/apis.md)
介绍json文件中涉及内部API接口的输入参数配置方法。
# [函数](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/function.md)
介绍json配置文件中涉及函数调用。
# [流程说明](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/flow.md)
介绍flow、schedule的调用流程。
# [测试](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/test.md)
介绍相关测试命令以及返回值。
# [导入](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/postmanswagger.md)
通过swagger，postman生成httpapi json。
# [需求](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/work.md)
介绍需求，开发计划以及开发进度。
# [promtheus](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/promtheus.md)
介绍promtheus中的指标。
# 隔离
使用`bucket`进行数据隔离。
# [参考](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/reference.md)
相关项目和文档。
