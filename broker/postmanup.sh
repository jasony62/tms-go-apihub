#!/bin/bash
# ############################################################
# 可配置文件路径及名称：
# 
# postman_collection_addr：     postman_collection文件相对位置
# postman_environment_addr：    postman_environment文件相对位置
# 
# ############################################################

postman_collection_addr="./APIHUB_0623.postman_collection"
postman_environment_addr="./34test_0623.postman_environment"

# ############################################################


echo "-------auto run postman-------"

run1="./APIHUB_0623.postman_collection"
run2="./34test_0623.postman_environment"
if [ -f $run1 -a $run2 ];then
    newman run $postman_collection_addr -e $postman_environment_addr
    echo "success: postman运行结束"
else
    echo "error: ./*.postman_collection 和 *.postman_environment 文件不存在"
fi

# newman run APIHUB_0623.postman_collection -e 34test_0623.postman_environment

# echo "success: "postman测试启动成功"
