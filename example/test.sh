#!/bin/sh
echo "\r\n高德地图:"
echo "\r\n查询城市区域编码:\r\n"
curl "http://localhost:8080/api/amap_district?city=北京"

echo "\r\n\r\n根据区域编码获得天气数据:\r\n"
curl "http://localhost:8080/api/amap_weather?city=110100"

echo "\r\n通过编排实现直接指定城市名称获得天气数据:\r\n"
curl "http://localhost:8080/flow/amap_city_weather?city=北京"

echo "\r\njson版本天气编排:\r\n"
curl -H "Content-Type: application/json" -d '{"city":"北京"}' "http://localhost:8080/flow/amap_city_weather_json"

echo "\r\n科大讯飞 NLP:"
echo "\r\n对输入内容进行分词:\r\n"
curl -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/api/kdxf_nlp_cws"

echo "\r\n对输入内容标注词性:\r\n"
curl -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/api/kdxf_nlp_pos"

echo "\r\n对输入内容提取关键词:\r\n"
curl -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/api/kdxf_nlp_ke"

echo "\r\n组合文本处理结果:\r\n"
curl -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/flow/kdxf_nlp"

echo "\r\n企业微信:"
echo "\r\n获得 access_token:\r\n"
curl "http://localhost:8080/api/qywx_gettoken"

echo "\r\n发送文本消息\r\n"
curl -X POST -d '{"touser": "YangYue","msgtype": "text","agentid": "1000002","text": { "content": "试试企业微信" }}' "http://localhost:8080/flow/qywx_message_send"


