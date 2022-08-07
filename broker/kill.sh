# 检擦api网关 http-server程序是否已经启动，确保不会因为重复启动导致端口占用问题
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