
echo "\r\npostman转json文件,信元_5G新增手机终端画像_首注册模型（数据自注册-Android）接口POST:\r\n"
curl -i -X POST  "http://localhost:8080/httpapi/信元_5G新增手机终端画像_首注册模型（数据自注册-Android）接口"

echo "\r\npostman转json文件,天津_发送接口POST:\r\n"
curl -i -H "Content-Type: application/json" "http://localhost:8080/httpapi/天津_发送接口"

echo "\r\npostman转json文件,天津_查询余额POST:\r\n"
curl -i -X POST  "http://localhost:8080/httpapi/天津_查询余额"

# echo "\r\npostman转json文件,schedule GET:\r\n"
# curl -i -H "Content-Type: application/json" "http://localhost:8080/schedule/postmanPressureGet?access_token=AT_APP_C3F6CFB0D2A6498BA98D88AED1EA8BEE&corporationText=FiberHome&modelText=MR820-LK&fwVersion=V4.0.3&macAddr=1CDE57DCA2AD&pinCode=00008023"

echo "\r\npostman转json文件,浙江信产_集团aWiFi漫游能力:\r\n"
curl -i -H "Content-Type: application/json" "http://localhost:8080/httpapi/浙江信产_集团aWiFi漫游能力"

curl -i -H "Content-Type: application/json" "http://localhost:8080/schedule/healthCheck"

curl -i -H "Content-Type: application/json" "http://localhost:8080/schedule/pressureTest"

curl -i -H "Content-Type: application/json" "http://localhost:8080/schedule/pressureArrange"
