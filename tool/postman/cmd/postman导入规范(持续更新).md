# 前言
postman软件支持多种多样输入内容，在接入APIHUB过程中，无法兼容全部场景，特此对postman输入内容进行些许限制，提高接入APIHUB的成功率。

# 输入要求（持续更新）

## API接口能力相关

* API名称中尽量减少 “ / ” 的使用，可更改为“ or ”、“ 或 ”
* API名称需明确接口能力，避免使用“ New Request ”名称
* 同一collection下的API名称不可以重复

## API接口能力Body内容
* 不支持Body内容中传入文件，例如：`KEY=file，VALUE=8K.wav `
* Body内容建议使用json或者json字符串格式

## Collection相关
* 导入程序目前最多支持4级collection目录，超出部分暂时无法解析
  * 原子能力平台
    * |__武汉电信实业
        * |__刘亦帆
            * |__行业三网短信能力
                * |__无忧短信接口
                    * |__短信发送接口（具体API能力）
                    * |__短信接收接口（具体API能力）