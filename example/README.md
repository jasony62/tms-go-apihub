# 高德地图

查询城市区域编码

```
curl "http://localhost:8080/api/amap_district?city=北京"
```

根据区域编码获得天气数据

```
curl "http://localhost:8080/api/amap_weather?city=110100"
```

通过编排实现直接指定城市名称获得天气数据

```
curl "http://localhost:8080/flow/amap_city_weather?city=北京"
```

入参为json的版本
curl -H "Content-Type: application/json" -d '{"city":"北京"}' "http://localhost:8080/api/amap_district"
curl -H "Content-Type: application/json" -d '{"city":"北京"}' "http://localhost:8080/flow/amap_city_weather_json"

# 科大讯飞 NLP

对输入内容进行分词

```
curl -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/api/kdxf_nlp_cws"
```

对输入内容标注词性

```
curl -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/api/kdxf_nlp_pos"
```

对输入内容提取关键词

```
curl -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/api/kdxf_nlp_ke"
```

组合文本处理结果

```
curl -X POST -d '{"content": "北京的天气"}' "http://localhost:8080/flow/kdxf_nlp"
```

# 企业微信

## 获得 access_token

```
curl "http://localhost:8080/api/qywx_gettoken"
```

## 发送文本消息

先获取 access_token 再发送消息

```
curl -X POST -d '{"touser": "YangYue","msgtype": "text","agentid": "1000002","text": { "content": "试试企业微信" }}' "http://localhost:8080/flow/qywx_message_send"
```

# 私有数据

API 调用设计使用与鉴权相关的私有数据，需要将这些数据放置在单独的文件中。

## 高德地图

`amap_keys.json`

```json
{
  "privates": [
    {
      "name": "key1",
      "value": "替换为实际值"
    }
  ]
}
```

## 科大讯飞

`kdxf_keys.json`

```json
{
  "privates": [
    {
      "name": "appid",
      "value": "替换为实际值"
    }
  ]
}
```

## 企业微信

`qywx_keys.json`

```json
{
  "privates": [
    {
      "name": "corpid",
      "value": "替换为实际值"
    },
    {
      "name": "corpsecret",
      "value": "替换为实际值"
    }
  ]
}
```
