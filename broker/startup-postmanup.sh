#!/bin/bash
# ############################################################
# 可配置文件路径及名称：
# 
# apihub_addr：                 apihub应用程序相对位置
# conf_addr：                   json文件相对位置
# postman_collection_addr：     postman_collection文件相对位置
# postman_environment_addr：    postman_environment文件相对位置
# 
# ############################################################

apihub_addr="./tms-go-apihub"
conf_addr="../example/"
postman_collection_addr="./APIHUB_0623.postman_collection"
postman_environment_addr="./34test_0623.postman_environment"

# ############################################################

echo "-------auto run tms-go-apihub-------"

killnum=`ps -C tms-go-apihub -o pid=`
kill $killnum

killnum=`ps -C http-server -o pid=`
kill $killnum

runapihub="./tms-go-apihub"
../test/http-server/http-server --addr 127.0.0.1:6060 &
if [ -f "$runapihub" ];then
    $apihub_addr --base $conf_addr &
    echo "success: tms-go-apihub运行结束"
else
    echo "error: ./tms-go-apihub 可执行文件不存在"
    go build -o tms-go-apihub
    if [ -f "$runapihub" ];then
        $apihub_addr --base $conf_addr &
    else
        echo "error: 重新编译源码失败, 请检查源码是否正确"
    fi
fi

sleep 2
echo "正在启动postman测试程序"


echo "-------auto run postman-------"

if [ -f $postman_collection_addr -a $postman_environment_addr ];then
    newman run $postman_collection_addr -e $postman_environment_addr
    echo "success: postman运行结束"
else
    echo "error: ./*.postman_collection 和 *.postman_environment 文件不存在"
fi
