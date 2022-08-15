# 前言

目前API描述或测试中最常使用的是Swagger和Postman脚本，团队自主研发的API网关程序低代码的方式是编写JSON文件，实现对API单个调用、流程调用、计划调用。

为提高API网关程序的兼容性，降低API网关接入难度，充分利用现有资源，对Swagger和Postman脚本进行自动化转换，即实现对OpenAPI和Postman脚本的兼容变得具有现实意义。

目前功能验证阶段，Swagger转换程序源码目录`./swagger`，`./swagger/from`子目录存放json或yaml格式的Swagger脚本，转换后的json文件目录`./swagger/to`。

Postman脚本转换程序源码目录`./postman`，`./postman/postman_collection`子目录存放postman collection导出的collection2.1规范的脚本文件，