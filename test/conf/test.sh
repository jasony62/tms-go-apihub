#!/bin/sh

echo "\r\n测试:"
echo "\r\n查询城市区域编码:\r\n"
curl -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/api/amap_district_test"

echo "\r\n测试:"
echo "\r\n\r\n根据区域编码获得天气数据:\r\n"
curl -H "Content-Type: application/json" -d '{"city": "110100"}' "http://localhost:8080/api/amap_weather_test"

echo "\r\n测试:"
echo "\r\n通过编排实现直接指定城市名称获得天气数据:\r\n"
curl -H "Content-Type: application/json" -d '{"city": "北京"}' "http://localhost:8080/flow/amap_city_weather_test"
