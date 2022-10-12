# prometheus GUI 使用教程
## 1 登录

### 1.1 网址

```html
http://localhost:30001/
```
或通过docker desktop软件入口进入

![在这里插入图片描述](https://img-blog.csdnimg.cn/99fdeb65661a4955ba4577f3d9cdb916.png)

### 1.2 登录账号、密码

* 账号：admin
* 密码：admin
  
**注意：修改密码窗口skip跳过即可**

## 2 导入GUI

### 2.1 Import json
![在这里插入图片描述](https://img-blog.csdnimg.cn/27451f6160204993b06f15b355c79b92.png)

![在这里插入图片描述](https://img-blog.csdnimg.cn/7647535ecedf4e538db4ce38c33359ea.png)

## 3 配置prometheus

### 3.1 Configuration

![在这里插入图片描述](https://img-blog.csdnimg.cn/43ccbb5fec4f4db59676f96ef1537656.png)

URL输入


```c
http://gah-prometheus:9090
```

![在这里插入图片描述](https://img-blog.csdnimg.cn/4d958ac78ef04926aec2a4b28fdfb344.png)

拉到最下，save

![在这里插入图片描述](https://img-blog.csdnimg.cn/f234a18cf3db411e80b40f0902b68b91.png)

### 3.2 配置Dasboards

按顺序找到Dasboards文件，并打开

![在这里插入图片描述](https://img-blog.csdnimg.cn/546fd5ff3fd14c3ca723c2d7d7728b30.png)

进入GUI界面后点击右上角 **设置** 按钮

![在这里插入图片描述](https://img-blog.csdnimg.cn/1545888fdc9344209cdad5588161a617.png)

依次配置三个Variable

![在这里插入图片描述](https://img-blog.csdnimg.cn/4ccd2f8a245442d395cc2b1ba30a2dfb.png)
配置Variable的Data source为Prometheus，然后Update
![在这里插入图片描述](https://img-blog.csdnimg.cn/e24b484dd80d468cad39602ab83c5a5f.png)
三个Variable设置完成后，保存设置
![在这里插入图片描述](https://img-blog.csdnimg.cn/da6f7af9bd1e418d937888a6490a8fd1.png)
返回主界面等待结果的抓取即可


# 界面说明（v.0.20221010）

![在这里插入图片描述](https://img-blog.csdnimg.cn/b1f89502c1b141bd93cf88e82ed93a49.jpeg)
左边图，表示请求时间；右边图表示请求速率

右边请求速率图速率非0时，表示本时刻发送API请求，左侧相应时间显示该时刻请求的平均时间

	注：目前模式是健康度检测，即拨测一次停一分钟，所以每个API的速率最高是1p/s