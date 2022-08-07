#!/bin/bash
# ############################################################
# 可配置文件路径及名称：
# 
# apihub_app：                 apihub应用程序相对位置
# conf_path：                  json文件夹相对位置
# 
##############################################################

apihub_app="./tms-go-apihub"
conf_path="../example/"

##############################################################

ps aux | grep tms-go-apihub
killnum=`ps -C tms-go-apihub -o pid=`
kill $killnum
if [ -f $? ];then
    echo "error: tms-go-apihub进程杀死失败"
else
    echo "success: tms-go-apihub进程杀死成功!"
fi

# 启动tms-go-apihub服务程序
echo "--------------auto run tms-go-apihub--------------"
if [ -f "$apihub_app" ];then
    $apihub_app --base $conf_path &
    echo "success: tms-go-apihub后台运行成功!"
else
    echo "error: ./tms-go-apihub 可执行文件不存在"
    go build -o tms-go-apihub
    if [ -f "$apihub_app" ];then
        $apihub_app --base $conf_path &
        echo "success: tms-go-apihub重新编译且后台运行成功!"
        go get
    else
        echo "error: tms-go-apihub重新编译源码失败, 请检查 ./ 源码是否正确"
        echo "tms-go-apihub应用程序启动失败, 未找到源码或可执行文件, 请检查shell文件路径"
    fi
fi

echo "提示: 如果需要重新编译应用程序, 直接删除可执行文件, 重新运行本脚本即可!"
echo "删除命令提示: rm ./tms-go-apihub"