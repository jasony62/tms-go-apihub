
## 环境变量

| 环境变量                      | 用途                                                         | 默认值  |
| ----------------------------- | ------------------------------------------------------------ | ------- |
| TGAH_HOST                     | 服务的主机名                                                 | 0.0.0.0 |
| TGAH_PORT                     | 服务的端口号                                                 | 8080    |
| TGAH_BUCKET_ENABLE            | API 和 FLOW 是否按 bucket 隔离                               | no      |
| TGAH_CONF_BASE_PATH             | API 定义文件存放位置                                         | ./conf       |
| TGAH_REMOTE_CONF_URL          | 从远端http服务器下载conf压缩包的路径，未配置则不下载，直接使用本地文件                         | -       |
| TGAH_REMOTE_CONF_UNZIP_PWD    | 下载文件解压密码，如果有密码则写解压密码，如果没有则不填     | -       |
## CONF目录结构
| 环境变量                      | 用途                                                         |
| ----------------------------- | ------------------------------------------------------------ |
| privates                     | 存放密码文件                                                 |
| apis                         | 存放API定义文件|
| flows            | 存放FLOW定义文件                               |
| schedules             | 存放SCHEDULE定义文件                                         |
| plugins             | 存放动态注册func的.so                                         |

## 命令行

通过`--env`指定使用的环境变量文件。

```
run go . --env envfile
```

```
run build -o tms-gah-broker
```

```
./tms-gah-broker --env envfile
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