#!/bin/bash
# ############################################################
# 可配置文件路径及名称：
# 
# apihub_addr：                 apihub应用程序相对位置
# conf_addr：                   json文件相对位置
# 
# ############################################################

apihub_addr="./tms-go-apihub"
conf_addr="../example/"

# ############################################################
echo "-------auto run tms-go-apihub-------"
killnum=`ps -C tms-go-apihub -o pid=`
kill $killnum
runapihub="./tms-go-apihub"
if [ -f "$runapihub" ];then
    $apihub_addr --base $conf_addr
    echo "success: tms-go-apihub运行结束"
else
    echo "error: ./tms-go-apihub 可执行文件不存在"
    go build -o tms-go-apihub
    if [ -f "$runapihub" ]
    then
        $apihub_addr --base $conf_addr
        go get
    else
        echo "error: 重新编译源码失败, 请检查源码是否正确"
    fi
fi





