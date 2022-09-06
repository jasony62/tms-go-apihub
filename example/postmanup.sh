#!/bin/bash
# ############################################################
# 可配置文件路径及名称：
# 
# postman_collection_app：     postman_collection文件相对位置
# postman_environment_app：    postman_environment文件相对位置
# 
# ############################################################

postman_collection_app="./postman/APIHUB_0623.postman_collection"
postman_environment_app="./postman/34test_0623.postman_environment"

# ############################################################


echo "正在启动postman测试程序"
echo "--------------auto run postman--------------"
if [ -f $postman_collection_app -a $postman_environment_app ];then
    newman run $postman_collection_app -e $postman_environment_app
    echo "success: postman运行结束"
else
    echo "error: ./*.postman_collection 和 *.postman_environment 文件不存在"
fi



