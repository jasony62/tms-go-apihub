
## CONF目录结构
| 名称                     | 用途                                                         |
| ----------------------------- | ------------------------------------------------------------ |
| main.json | 启动文件 |
| privates                     | 文件夹，存放密码文件                                                 |
| httpapis                         | 文件夹，存放HTTPAPI定义文件|
| flows            | 文件夹，存放FLOW定义文件                               |
| schedules             | 文件夹，存放SCHEDULE定义文件                                         |
| plugins             | 文件夹，存放动态注册func的.so                                         |
| templates | 文件夹，存放html tmpl文件 |
| rights | 文件夹，存放httpapi，flow和schedule对应的权限列表 |

## 命令行
通过`--env`指定使用的环境变量文件，通过`--base`指定conf文件夹的路径，默认为./conf/。
启动时读取base路径下的main.json启动，其定义为通用的flow结构，文件里的变量可以写死，也可以从env里获取。

```
run go . --env envfile
```

```
run build -o tms-gah-broker
```

```
./tms-gah-broker --env envfile --base ./conf/
```

## docker

```
docker build -t tms/gah-broker .
```

```
docker run -it --rm --name tms-gah-broker -p 8080:8080 tms/gah-broker sh
```

```
docker compose build tms-gah-broker
```

```
docker compose up tms-gah-broker
```

## 安装插件
插件编译不依赖于本代码。

```
cd plugins
cd kdxfnlp
go build -buildmode=plugin -o kdxfnlp.so kdxfnlp.go
```
将生成的.so放到conf/plugins下，模块启动时候会自动加载