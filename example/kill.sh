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
echo "****************************************************" 
echo "提示: 如果需要重新编译应用程序, 直接删除可执行文件, 重新运行本脚本即可!"
echo "删除命令提示: rm ./tms-go-apihub"
echo "删除命令提示: rm ../test/http-server/http-server"
echo "****************************************************"
