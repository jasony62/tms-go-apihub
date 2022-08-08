#!/bin/bash
##############################################################
# 可配置文件路径及名称：
# 
# apihub_app：                 apihub应用程序相对位置
# conf_path：                  json文件夹相对位置
#
# postman_collection_app：     postman_collection文件相对位置
# postman_environment_app：    postman_environment文件相对位置
#
# httpserver_app：             httpserver应用程序相对位置
# httpserver_ip：              httpserver应用程序默认监听IP和端口号
# 
##############################################################
####################自定义位置#################################

apihub_app="./tms-go-apihub"
conf_path="../example/"

postman_collection_app="./APIHUB_0623.postman_collection"
postman_environment_app="./34test_0623.postman_environment"

httpserver_app="../test/http-server/http-server"
httpserver_path="../test/http-server/"
httpserver_ip="127.0.0.1:6060"

##############################################################
##############################################################

# 检查api网关 http-server程序是否已经启动，确保不会因为重复启动导致端口占用问题
ps aux | grep tms-go-apihub
killnum=`ps -C tms-go-apihub -o pid=`
kill $killnum
if [ -f $? ];then
    echo "error: tms-go-apihub进程杀死失败"
else
    echo "success: tms-go-apihub进程杀死成功!"
fi

ps aux | grep http-server
killnum=`ps -C http-server -o pid=`
kill $killnum
if [ -f $? ];then
    echo "error: http-server进程杀死失败"
else
    echo "success: http-server进程杀死成功!"
fi

# 启动http-server服务程序
echo "--------------auto run http-server--------------"
# ../test/http-server/http-server --addr 127.0.0.1:6060 &
if [ -f "$httpserver_app" ];then
    $httpserver_app --addr $httpserver_ip &
    echo "success: http-server后台运行成功!"
else
    echo "error: ../test/http-server/http-server 可执行文件不存在"
    cd $httpserver_path
    echo "running: 正在重新编译..."
    go get
    go build -o http-server
    cd -
    if [ -f "$httpserver_app" ];then
        $httpserver_app --addr $httpserver_ip &
        echo "success: http-server重新编译且后台运行成功!"
    else
        echo "error: http-server重新编译源码失败, 请检查 .../test/http-server/ 源码是否正确"
        echo "http-server应用程序启动失败, 未找到源码或可执行文件, 请检查shell文件路径"
    fi
fi

# 启动tms-go-apihub服务程序
echo "--------------auto run tms-go-apihub--------------"
if [ -f "$apihub_app" ];then
    $apihub_app --base $conf_path &
    echo "success: tms-go-apihub后台运行成功!"
else
    echo "error: ./tms-go-apihub 可执行文件不存在"
    echo "running: 正在重新编译..."
    go get
    go build -o tms-go-apihub
    if [ -f "$apihub_app" ];then
        $apihub_app --base $conf_path &
        echo "success: tms-go-apihub重新编译且后台运行成功!"
    else
        echo "error: tms-go-apihub重新编译源码失败, 请检查 ./ 源码是否正确"
        echo "tms-go-apihub应用程序启动失败, 未找到源码或可执行文件, 请检查shell文件路径"
    fi
fi

sleep 2

echo "正在启动postman测试程序"
echo "--------------auto run postman--------------"
if [ -f $postman_collection_app -a $postman_environment_app ];then
    newman run $postman_collection_app -e $postman_environment_app
    echo "success: postman运行结束"
else
    echo "error: ./*.postman_collection 和 *.postman_environment 文件不存在"
fi
