#!/bin/sh

echo "\r\n高德地图:"
echo "\r\n查询城市区域编码(进行多租户管理):\r\n"
curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/httpapi/amap_district?appID=001"

echo "\r\n查询城市区域编码(带API版本号):\r\n"
curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/httpapi/amap_district/v1?appID=001"

echo "\r\n\r\n根据区域编码获得天气数据:\r\n"
curl -i -H "Content-Type: application/json" -d '{"city": "110100"}' "http://localhost:8080/httpapi/amap_weather"

echo "\r\n通过编排实现直接指定城市名称获得天气数据Json格式:\r\n"
curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/flow/amap_city_weather"

echo "\r\n通过编排实现直接指定城市名称获得天气数据HTML格式:\r\n"
curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/flow/amap_city_weather_html"

echo "\r\n科大讯飞 NLP:"
echo "\r\n对输入内容进行分词:\r\n"
curl -i -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/httpapi/kdxf_nlp_cws"

echo "\r\n对输入内容标注词性:\r\n"
curl -i -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/httpapi/kdxf_nlp_pos"

echo "\r\n对输入内容提取关键词:\r\n"
curl -i -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/httpapi/kdxf_nlp_ke"

echo "\r\n组合文本处理结果:\r\n"
curl -i -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/flow/kdxf_nlp"

echo "\r\n企业微信:"
#echo "\r\n获得 access_token:\r\n"
#curl -i "http://localhost:8080/httpapi/qywx_gettoken"

echo "\r\n发送消息\r\n"
curl -i -X POST -d '{"touser": "YangYue","msgtype": "text","agentid": "1000002", "content": "试试企业微信" }' "http://localhost:8080/flow/qywx_message_send"


#echo "\r\n查询百度图片分类token"
#curl  -i "http://localhost:8080/httpapi/baidu_image_classify_token"

echo "\r\n通过编排从百度获得相关图片分类"
curl -i -X POST -d '{"content": "https://img.zcool.cn/community/01ff2059770a25a8012193a37c7695.jpg"}'  "http://localhost:8080/flow/baidu_image_classify"

echo "\r\n原子能力:对话情绪识别\r\n"
curl -i -X POST -H "Content-Type: application/json" -d '{"text": "hello!", "seqid": "c7574913-5a4f-4622-989c-455f9bd20640"}' "http://localhost:8080/httpapi/yznl_nlp_motion"

#echo "\r\n发送短信\r\n"
#curl -H "Content-Type: application/json" -d '{"number": "138104xxx69", "msg":"my test !!!"}' "http://localhost:8080/httpapi/sm_send"

echo "\r\nSCHEDULE:\r\n"
echo "\r\n并行查询天气和天气并发送企业微信"
curl -i -H "Content-Type: application/json" -d '{"cities":["sh", "bj", "ls", "sz"], "image":"https://img.zcool.cn/community/01ff2059770a25a8012193a37c7695.jpg"}' "http://localhost:8080/schedule/amap_qywx"


echo "\r\n地图服务查询"
curl -i -H "Content-Type: application/json" "http://localhost:8080/flow/gis_base_map?lat=39.915599&lng=116.406568"

#需要配置rights目录下的json文件，目前配置支持user为001,002,user是在query中
echo "\r\nquery带用户appID的地区查询"   
curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/httpapi/amap_district?appID=001"

#不支持
#echo "\r\nheader带用户appID的地区查询"   
#curl -i -H "Content-Type: application/json" -H "appID: 001" -d '{"city": "北京"}' "http://localhost:8080/httpapi/amap_district"

echo "\r\nquery带用户appID的城市天气查询"   
curl -i -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/flow/amap_city_weather?appID=001"

echo "\r\nheader带用户appID的城市天气查询"   
curl -i -H "Content-Type: application/json" -H "appID: 001" -d '{"city": "北京"}' "http://localhost:8080/flow/amap_city_weather"

echo "\r\nquery带用户appID的企业微信schedule查询"  
curl -i -H "Content-Type: application/json" -d '{"cities":["sh", "bj", "sh", "sh"], "image":"https://img.zcool.cn/community/01ff2059770a25a8012193a37c7695.jpg"}' "http://localhost:8080/schedule/amap_qywx?appID=001"


