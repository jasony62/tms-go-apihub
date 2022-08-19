# 调试命令   
参考test.sh

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

## 百度图片

`baidu_image_classify_key.json`

```json
{
  "privates": [
    {
      "name": "xappkey",
      "value": "iih1Zs1Vn0xCICICCfduxI0O"
    },
	{
      "name": "xsecret",
      "value": "r32evIieR4IumIlvUQnlwDyVP8jIeTvU"
    },
	{
      "name": "xgranttype",
      "value": "client_credentials"
    }
  ]
}
```

## 
