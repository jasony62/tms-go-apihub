# 前言

目前API描述或测试中最常使用的是Swagger和Postman脚本，团队自主研发的API网关程序低代码的方式是编写JSON文件，实现对API单个调用、流程调用、计划调用。

为提高API网关程序的兼容性，降低API网关接入难度，充分利用现有资源，对Swagger和Postman脚本进行自动化转换，即实现对Swagger和Postman脚本的兼容变得具有现实意义。

# postman脚本转换程序使用说明

目前功能验证阶段，Postman脚本转换程序源码目录`./postman`，`./postman/postman_collection`子目录默认存放postman collection导出的collection2.1规范的脚本文件，`./postman/jsonFiles`子目录默认存放转换后的json文件。

使用方法

源码目录`./postman`下，执行`go get`获取相关依赖包

```
go get
```

源码目录`./postman`下，执行build命令，生成名叫coverPostmanApp的可执行文件

```
go build -o coverPostmanApp
```

运行coverPostmanApp可执行文件

```
./coverPostmanApp
```
或指定导入文件目录和导出文件目录


```
./coverPostmanApp  --from ./postman_collection --to ./jsonFiles
```

`--from` 读取`./postman_collection`子目录下postman脚本文件

`--to `导出至`./jsonFiles`子目录下，导出文件名命名规则，postman collection名称_request名称.json。

# swagger脚本转换程序使用说明
目前功能验证阶段，Swagger转换程序源码目录`./swagger`，`./swagger/from`子目录默认存放json或yaml格式的Swagger脚本，`./swagger/to`子目录默认存放转换后的json文件。

使用方法

源码目录`./swagger`下，执行`go get`获取相关依赖包

```
go get
```

源码目录`./swagger`下，执行build命令，生成名叫coverSwaggerApp的可执行文件

```
go build -o coverSwaggerApp
```

运行coverSwaggerApp可执行文件

```
./coverSwaggerApp --from ./from --to ./to
```