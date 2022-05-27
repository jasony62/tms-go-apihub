package apis

import (
	"encoding/json"
	"net/url"
	"strconv"

	"strings"
	"time"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
	"github.com/valyala/fasthttp"

	jsoniter "github.com/json-iterator/go"
)

//json反序列化是造成整数的精度丢失，所以使用一个扩展的json工具做反序列化
var jsonEx = jsoniter.Config{
	UseNumber: true,
}.Froze()

func handleReq(stack *hub.Stack, HttpApi *hub.HttpApiDef, privateDef *hub.PrivateArray) (interface{}, int) {
	var jsonInRspBody interface{}

	outReq, err := newRequest(stack, HttpApi, privateDef)
	if err != nil {
		return nil, fasthttp.StatusInternalServerError
	}
	// 发出请求
	client := &fasthttp.Client{}
	resp := fasthttp.AcquireResponse()
	err = client.Do(outReq, resp)
	defer fasthttp.ReleaseRequest(outReq)
	if err != nil {
		klog.Errorln("ERR Connection error: ", err)
		return nil, fasthttp.StatusInternalServerError
	}
	defer fasthttp.ReleaseResponse(resp)
	returnBody := resp.Body()
	// 将收到的结果转为JSON对象
	jsonEx.Unmarshal(returnBody, &jsonInRspBody)
	stack.StepResult[hub.ResultName] = jsonInRspBody

	klog.Errorln("消息体: ", string(returnBody))

	if HttpApi.Cache != nil {
		//解析过期时间，如果存在则记录下来
		expires, ok := handleExpireTime(stack, HttpApi, resp)
		if !ok {
			klog.Warningln("没有查询到过期时间")
		} else {
			klog.Infof("更新Cache信息，过期时间为: %v", expires)
			HttpApi.Cache.Expires = expires
			HttpApi.Cache.Resp = jsonInRspBody
		}
	}

	return jsonInRspBody, fasthttp.StatusOK
}

func newRequest(stack *hub.Stack, HttpApi *hub.HttpApiDef, privateDef *hub.PrivateArray) (*fasthttp.Request, error) {
	var outBody string
	var hasBody bool
	var err error
	// 要发送的请求
	outReq := fasthttp.AcquireRequest()
	outReq.Header.SetMethod(HttpApi.Method)
	hasBody = len(HttpApi.RequestContentType) > 0 && HttpApi.RequestContentType != "none"
	if hasBody {
		switch HttpApi.RequestContentType {
		case "form":
			outReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		case "json":
			outReq.Header.Set("Content-Type", "application/json")
		case hub.OriginName:
			contentType := stack.GinContext.Request.Header.Get("Content-Type")
			outReq.Header.Set("Content-Type", contentType)
			// 收到的请求中的数据
			inData, _ := json.Marshal(stack.StepResult[hub.OriginName])
			outBody = string(inData)
		default:
			outReq.Header.Set("Content-Type", HttpApi.RequestContentType)
		}
	}

	// 发出请求的URL
	var finalUrl string
	var args fasthttp.Args
	if len(HttpApi.Url) == 0 {
		if HttpApi.DynamicUrl != nil {
			finalUrl, err = util.GetParameterStringValue(stack, privateDef, HttpApi.DynamicUrl)
			if err != nil {
				return nil, err
			}
		} else {
			klog.Errorln("无有效url")
			panic("无有效url")
		}
	} else {
		finalUrl = HttpApi.Url
	}
	outReqURL, _ := url.Parse(finalUrl)
	// 设置请求参数
	outReqParamRules := HttpApi.Parameters
	if outReqParamRules != nil {
		paramLen := len(*outReqParamRules)
		if paramLen > 0 {
			var value string
			q := outReqURL.Query()
			vars := make(map[string]string, paramLen)
			stack.StepResult[hub.VarsName] = vars
			defer delete(stack.StepResult, hub.VarsName)

			for _, param := range *outReqParamRules {
				if len(param.Name) > 0 {
					value, err = util.GetParameterStringValue(stack, privateDef, &param.Value)
					if err != nil {
						return nil, err
					}

					switch param.In {
					case "query":
						q.Set(param.Name, value)
					case "header":
						outReq.Header.Set(param.Name, value)
					case "body":
						if hasBody && HttpApi.RequestContentType != hub.OriginName {
							if HttpApi.RequestContentType == "form" {
								args.Set(param.Name, value)
							} else {
								if len(outBody) == 0 {
									if value == "null" {
										klog.Errorln("获得body失败：")
										panic("获得body失败：")
									} else {
										outBody = value
										klog.Infoln("Set body :\r\n", outBody, "\r\n", len(outBody))
									}
								} else {
									klog.Infoln("Double content body :\r\n", outBody, "\r\nVS\r\n", value)
								}
							}
						} else {
							klog.Infoln("Refuse to set body :", HttpApi.RequestContentType, "VS\r\n", value)
						}
					case hub.VarsName:
					default:
						klog.Infoln("Invalid in:", param.In, "名字", param.Name, "值", value)
					}
					vars[param.Name] = value
					klog.Infoln("设置入参，位置", param.In, "名字", param.Name, "值", value)
				}
			}
			outReqURL.RawQuery = q.Encode()
		}
	}
	outReq.SetRequestURI(outReqURL.String())

	// 处理要发送的消息体
	if HttpApi.Method == "POST" {
		if HttpApi.RequestContentType != "none" {
			if HttpApi.RequestContentType == "form" {
				args.WriteTo(outReq.BodyWriter())
			} else {
				outReq.SetBodyString(outBody)
			}
		}
	}

	return outReq, nil
}

func handleExpireTime(stack *hub.Stack, HttpApi *hub.HttpApiDef, resp *fasthttp.Response) (time.Time, bool) {
	klog.Infoln("获得参数，[src]:", HttpApi.Cache.Expire.From, "; [key]:", HttpApi.Cache.Expire.Content, "; [format]:", HttpApi.Cache.Format)
	if strings.EqualFold(HttpApi.Cache.Expire.From, "header") {
		return handleHeaderExpireTime(HttpApi, resp)
	} else {
		return handleBodyExpireTime(stack, HttpApi)
	}
}

func handleHeaderExpireTime(HttpApi *hub.HttpApiDef, resp *fasthttp.Response) (time.Time, bool) {
	//首先在api 的json文件中配置参数 cache
	// "cache": {
	// 	"value": {
	// 		"from": "header",
	// 		"name": "Set-Cookie.expires"
	// 	},
	// 	"format": "Mon, 02-Jan-06 15:04:05 MST"
	//   }
	//from 为从header还是从body中获取过期时间
	//name 为获取过期时间的关键字串
	//format：如果是date格式，则配置具体格式串，如果是second数，则按照秒数解析
	//	baidu_image_classify_token: Mon, 02-Jan-06 15:04:05 MST
	//	body中一个例子："expireTime":"20220510153521",格式为：20060102150405

	//format = "20060102150405"
	key := HttpApi.Cache.Expire.Content
	format := HttpApi.Cache.Format

	if strings.Contains(key, "Set-Cookie.") {
		key = strings.TrimPrefix(key, "Set-Cookie.")
		//判断Set-Cookie中是否含有Expires 的header
		cookie := resp.Header.Peek("Set-Cookie")
		klog.Infoln("Header中Set-Cookie: ", cookie)
		if len(cookie) > 0 {
			expiresIndex := strings.Index(string(cookie), key) //"expires="
			if expiresIndex >= 0 {
				semicolonIndex := strings.Index(string(cookie[expiresIndex:]), ";")
				if semicolonIndex < 0 {
					semicolonIndex = 0
				}

				expires, err := parseExpireTime(string(cookie[expiresIndex+len(key)+1:expiresIndex+semicolonIndex]), format)
				if err == nil {
					return expires, true
				}
			}
		}
	} else {
		//判断是否含有Expires 的header
		expires, err := parseExpireTime(string(resp.Header.Peek(key)), format)
		if err == nil {
			return expires, true
		}
	}

	return time.Time{}, false
}

func handleBodyExpireTime(stack *hub.Stack, HttpApi *hub.HttpApiDef) (time.Time, bool) {
	//首先在api 的json文件中配置参数 cache
	// "cache": {
	// 	"value": {
	// 		"from": "json",
	// 		"name": "{{.result.expires_in}}"
	// 	},
	// 	"format": "second"
	//   }
	//name 为获取过期时间的关键字串
	//format：如果是date格式，则配置具体格式串，如果是second数，则按照秒数解析
	//	baidu_image_classify_token: Mon, 02-Jan-06 15:04:05 MST
	//	body中一个例子："expireTime":"20220510153521",格式为：20060102150405

	format := HttpApi.Cache.Format
	result, err := util.GetParameterStringValue(stack, nil, HttpApi.Cache.Expire)
	if err != nil {
		return time.Time{}, false
	}

	klog.Infof("handleBodyExpireTime:", result)

	formatTime, err := parseExpireTime(result, format)
	if err == nil {
		return formatTime, true
	}

	return time.Time{}, false
}

func parseExpireTime(value string, format string) (time.Time, error) {
	var exptime time.Time
	var err error

	if strings.EqualFold(format, "second") {
		seconds, err := strconv.Atoi(value)
		if err != nil {
			klog.Errorln("解析过期时间失败, err: ", err)
			return time.Time{}, err
		}
		klog.Infoln("解析后过期秒数: ", seconds)
		exptime = time.Now().Add(time.Second * time.Duration(seconds))
	} else {
		exptime, err = time.Parse(format, value)
		if err != nil {
			klog.Errorln("解析过期时间失败, err: ", err)
			return time.Time{}, err
		}
	}
	klog.Infoln("解析后过期时间: ", exptime)
	return exptime.Local(), nil
}

func getCacheContent(HttpApi *hub.HttpApiDef) interface{} {
	//如果支持缓存，判断过期时间
	if time.Now().Local().After(HttpApi.Cache.Expires) {
		return nil
	}
	return HttpApi.Cache.Resp
}

func getCacheContentWithLock(HttpApi *hub.HttpApiDef) interface{} {
	//如果支持缓存，判断过期时间
	HttpApi.Cache.Locker.RLock()
	defer HttpApi.Cache.Locker.RUnlock()
	if time.Now().Local().After(HttpApi.Cache.Expires) {
		return nil
	}
	return HttpApi.Cache.Resp
}

// 转发API调用
func run(stack *hub.Stack, name string, private string) (jsonOutRspBody interface{}, ret int) {
	var err error
	var privateDef *hub.PrivateArray
	HttpApi, err := util.FindApiDef(stack, name)

	if HttpApi == nil {
		klog.Errorln("获得API定义失败：", err)
		panic(err)
	}
	if len(private) == 0 {
		private = HttpApi.PrivateName
	}

	if len(private) != 0 {
		privateDef, err = util.FindPrivateDef(stack, private)
		if err != nil {
			klog.Errorln("获得private定义失败：", err)
			panic(err)
		}
	}

	if HttpApi.Cache != nil { //如果Json文件中配置了cache，表示支持缓存
		if jsonOutRspBody = getCacheContentWithLock(HttpApi); jsonOutRspBody == nil {
			defer HttpApi.Cache.Locker.Unlock()
			HttpApi.Cache.Locker.Lock()

			if jsonOutRspBody = getCacheContent(HttpApi); jsonOutRspBody == nil {
				klog.Infoln("获取缓存Cache ... ...")
				jsonOutRspBody, _ = handleReq(stack, HttpApi, privateDef)
			} else {
				klog.Infoln("Cache缓存有效，直接回应")
			}
		} else {
			klog.Infoln("Cache缓存有效，直接回应")
		}
	} else { //不支持缓存，直接请求
		jsonOutRspBody, _ = handleReq(stack, HttpApi, privateDef)
	}

	klog.Infoln("处理", HttpApi.Url, ":", fasthttp.StatusOK, "\r\n返回结果：", jsonOutRspBody)
	if jsonOutRspBody == nil {
		return nil, fasthttp.StatusInternalServerError
	}
	return jsonOutRspBody, fasthttp.StatusOK
}

func runHttpApi(stack *hub.Stack, params map[string]string) (interface{}, int) {
	name, OK := params["name"]
	if !OK {
		str := "缺少api名称"
		klog.Errorln(str)
		panic(str)
	}

	/*private may doesn't exist*/
	private := params["private"]
	return run(stack, name, private)
}

func httpResponse(stack *hub.Stack, params map[string]string) (interface{}, int) {
	name, OK := params["type"]
	if !OK {
		str := "缺少api名称"
		klog.Errorln(str)
		panic(str)
	}

	key, OK := params["key"]
	if !OK {
		str := "缺少api名称"
		klog.Errorln(str)
		panic(str)
	}
	result := stack.StepResult[key]
	switch name {
	case "html":
		stack.GinContext.Header("Content-Type", "text/html; charset=utf-8")
		stack.GinContext.String(fasthttp.StatusOK, "%s", result)
	case "json":
		stack.GinContext.IndentedJSON(fasthttp.StatusOK, result)
	default:
		stack.GinContext.Header("Content-Type", name)
		stack.GinContext.String(fasthttp.StatusOK, "%s", result)
	}
	return nil, fasthttp.StatusOK
}
