| 修订版本 | 修订者 | 说明 |
| -- | -- | -- |
| v0.20220807 | Sheng-ZM | 新增shell脚本快速启动说明 |
# 1 前言
本章节主要介绍程序在`Linux`开发环境下的初次启动流程

本章节程序启动后，建立**API网关**进程，API网关进程监听具体调用命令。

|||
| -- | -- |
|操作系统|Linux / Windows|
|软件环境|Golang1.17及以上、GCC9及以上|
|涉及目录|broker、schema、example|

    注：源码也可在Windows环境进行编译运行，具体步骤见3.1章节注。

如果已经配置过Go运行环境，且环境已经有相关依赖包，Linux环境下可直接运行`./broker`目录下的`startup.sh`shell脚本文件。脚本将自动**检查、编译、启动**API网关程序。
    
    注：脚本提供了一个快速启动的方式，若出现错误，可参见详细启动与测试介绍，跳转到第3章节内容。

本篇目录结构如下：

* 第一章节介绍基本信息与环境配置。

* 第二章介绍源码涉及目录详情。

* 第三章介绍详细的启动过程以及编译测试过程。


## 1.1 Linux环境配置
### 1.1.1 Golang下载
进入Golang下载官网，需选择Linux最新Go版本下载即可
```
https://studygolang.com/dl
```
### 1.1.2 安装
将`go1.18.3.linux-amd64.tar.gz`压缩包下载到Linux指定位置，例如：/home/

执行解压指令
```
tar zxvf go1.18.3.linux-amd64.tar.gz
```
此时解压后的go文件夹路径为：/home/go
### 1.1.3 配置环境
建议在同一目录下建立go语言工作环境文件夹，新建一个gopath文件夹，路径为：/home/gopath

执行命令，打开profile文件
```
sudo vim /etc/profile
```
然后在打开的文件末尾添加：
```
export GOROOT=/home/go
export GOPATH=/home/gopath
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
export GOPROXY="https://goproxy.io"
```
然后刷新文档
```
source /etc/profile
```
以上配置好之后，我们打开终端，属于如下命令，就可以看到go的版本等信息了。
```
go version
```
## 1.2 Windows环境配置
### 1.2.1 Golang下载
进入Golang下载官网，需选择Windows最新Go版本下载即可
```
https://studygolang.com/dl
```
### 1.2.2 安装
执行可执行文件exe安装即可。
### 1.2.3 配置环境
右击此电脑->属性->高级系统设置->环境变量，打开环境变量设置窗口。

需要新建两个环境变量配置，一个是 GOROOT ，这个就是 Go 环境所在目录的配置。另一个是 GOPATH ，这个是 Go 项目的工作目录

# 2 目录介绍
在介绍如何编译和运行源码前，首先介绍源码所涉及目录，以及各目录主要功能。若不感兴趣，想要快速执行API网关程序，可直接跳转第3章节。
## 2.1 broker文件目录介绍

apihub源码均在`broker`目录下，即针对源码的`build`编译在broker目录下运行。


|名称|用途|
| -- | -- |
|apis| API信息 |
|hub|API信息结构定义|
|util|工具包|
|main.go|主程序|
|go.mod|相关Go包的集合，替换旧的基于`GOPATH`的方法，来指定使用哪些源文件，执行`go get`可自动下载相关依赖包到本地gopath默认位置|

## 2.2 schema文件目录介绍
`schema`目录中存放`json schema`校验文件，`schema`检查`example`文件夹中json文件是否合法，即是否符合API网关格式规范。

schema目录下:

| 名称| 用途|
| -- | -- |
| httpapi.json | 定义httpapi的json schema |
| flow.json | 定义flow的json schema |
| right.json | 定义right的json schema |
| schedule.json | 定义schedule的json schema |

## 2.3 example文件目录介绍

`example`目录中存放了规范后的API接口信息，以及API编排后的`flow`，均为`json`格式。

    注：本文中低代码指的是，用户通过配置json文件，编写一个服务的flow，即可低成本，低代码的方式实现API编排功能，而无需理解API网关程序具体工作流程。

API网关与API服务配置文件相互分离。一方面，增加了程序部署的灵活，用户仅需要使用方按照规范提供API接口信息的`json`文件（后续通过mongoDB Web可视化生成和添加API接口信息的json文件），即可针对需求通过低代码的方式设计flow，实现低代码的API编排和调用；另一方面，apihub程序由开发人员维护升级，用户无需关心，降低了用户学习成本。

`example`目录说明：
|名称| 用途|
| -- | -- |
| main.json | 启动文件 |
| privates| 文件夹，存放密码文件|
| httpapis| 文件夹，存放HTTPAPI定义文件|
| flows| 文件夹，存放FLOW定义文件|
| schedules| 文件夹，存放SCHEDULE定义文件|
| plugins| 文件夹，存放动态注册func的.so|
| templates | 文件夹，存放html tmpl文件 |
| rights | 文件夹，存放httpapi，flow和schedule对应的权限列表 |

json文件定义参考[JSON定义](https://github.com/jasony62/tms-go-apihub/blob/main/doc/cn/json.md)

### 2.3.1 flows、httpapis、schedules
`broker`目录中`flows`、`httpapis`、`schedules`三个文件夹主要存放API相关文件。

### 2.3.2 privates
因为`privates`文件夹存放密码文件，所以没有暴露在git，即git中查找不到为正常现象。

## 2.4 文件关联
上述主要的三个文件夹`example`、`broker`、`schema`具体工作关系如下:

`broker`目录下编译并执行`tms-go-apihub`文件，`tms-go-apihub`文件通过`schema`文件夹中的校验格式，检查`example`文件夹中json文件的合法性。最后启动API网关的监听服务。

# 3 apihub启动
## 3.1 build
源码中已经初始化完成`go.mod`文件，即已经生成依赖包地址。

在`broker`程序源码文件下，执行命令，下载依赖包到主机gopath默认位置。
```
go get
```

在`broker`程序源码文件下，执行命令
```
go build -o tms-go-apihub
```
命令生成名称为`tms-go-apihub`的可执行文件。
    
    注：若build编译Windows版本，则需要在PowerShell终端窗口执行 
    go build -o tms-go-apihub.exe 
    生成exe可执行文件，后续步骤与Linux环境下一致
## 3.2 run
由于程序源码的代码预设，需要手动对某些文件进行指定或者软连接，也可根据实际文件位置与运行状态修改代码灵活调整，这里按照源代码预设路径进行操作。

* 方法一：`--env`指定环境变量文件
通过`--env`指定使用的环境变量文件(非必须，后续可以通过args里的from env访问)，本方法可跳过，直接使用方法二。

* 方法二：`--base`指定conf文件夹的路径
通过`--base`命令指定`tms-go-apihub`读取的`example`文件夹的路径。

默认`tms-go-apihub`程序需要调用`example`下的API信息，因此需要将`example`文件夹软链接或者拷贝到`broker/conf`目录下(前文中方法二便是通过`--base`命令直接指定文件夹路径，避免了软连接或者拷贝的操作，如果不能理解，直接按照后续命令执行亦可)

在`broker`文件下，具体执行命令如下,通过`--base`命令指定`tms-go-apihub`程序读取的`example`文件夹的路径：

```
./tms-go-apihub --base ../example/
```

程序运行后，正常应处于监听状态
```
[GIN-debug] Listening and serving HTTP on 0.0.0.0:8080
```

若出现异常现象，需要根据打印提示进行微调，例如：
* 1.根据提示删除某些无效json文件；
* 2.若提示端口号被占用，需要修改`/example/main.json`更改端口号

apihub程序启动后，打开新的终端窗口，执行curl命令发送请求，进而获取信息
```
curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/httpapi/amap_district"
```
新的终端窗口打印输出如下内容
```
ocalhost:8080/httpapi/amap_district"
HTTP/1.1 403 Forbidden
Content-Type: application/json; charset=utf-8
Date: Tue, 19 Jul 2022 09:30:44 GMT
Content-Length: 4
```
`tms-go-apihub`运行窗口打印输出如下内容
```
[GIN] 2022/07/19 - 17:30:44 | 403 |      3.2014ms |       127.0.0.1 | POST     "/httpapi/amap_district"
```
至此，apihub程序正常启动并工作。
# 4 补充
若有响应需求，可进行补充操作，或者功能查找。
## 4.1 docker
若在docker环境下运行，执行如下命令。


```
docker build -t tms/gah-broker .
```

```
docker run -it --rm --name tms-gah-broker -p 8080:8080 tms/gah-broker sh
cd broker/
./tms-gah-broker --base ../example/
```

```
docker compose build tms-gah-broker
```

```
docker compose up tms-gah-broker
```

## 4.2 安装插件
插件编译不依赖于本代码。

```
cd plugins
cd kdxfnlp
go build -buildmode=plugin -o kdxfnlp.so kdxfnlp.go
```
将生成的.so放到conf/plugins下，模块启动时候会自动加载