package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"strings"
	"time"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/unit"
	"github.com/jasony62/tms-go-apihub/util"
	jsoniter "github.com/json-iterator/go"
)

//json反序列化是造成整数的精度丢失，所以使用一个扩展的json工具做反序列化
var jsonEx = jsoniter.Config{
	UseNumber: true,
}.Froze()

func handleReq(stack *hub.Stack, apiDef *hub.ApiDef, privateDef *hub.PrivateArray) (interface{}, int) {
	var jsonInRspBody interface{}

	outReq := newRequest(stack, apiDef, privateDef)
	// 发出请求
	client := &http.Client{}
	resp, err := client.Do(outReq)
	if err != nil {
		klog.Errorln("err", err)
		return nil, 500
	}
	defer resp.Body.Close()
	returnBody, _ := io.ReadAll(resp.Body)

	// 将收到的结果转为JSON对象
	jsonEx.Unmarshal(returnBody, &jsonInRspBody)
	stack.StepResult[hub.ResultName] = jsonInRspBody

	klog.Errorln("消息体: ", string(returnBody))

	if !handleRespStatus(stack, apiDef) {
		klog.Errorln("消息体中返回码显示不成功，回应错误")
		return nil, 500
	}

	out := newOutRspBody(apiDef, jsonInRspBody)

	if apiDef.Cache != nil {
		//解析过期时间，如果存在则记录下来
		expires, ok := handleExpireTime(stack, apiDef, resp)
		if !ok {
			klog.Warningln("没有查询到过期时间")
		} else {
			klog.Infof("更新Cache信息，过期时间为: %v", expires)
			apiDef.Cache.Expires = expires
			apiDef.Cache.Resp = out
		}
	}

	return out, http.StatusOK
}

// 构造发送的响应内容
func newOutRspBody(apiDef *hub.ApiDef, in interface{}) interface{} {
	var out interface{}
	if apiDef.Response != nil && apiDef.Response.Json != nil {
		out = util.Json2Json(in, apiDef.Response.Json)
	} else {
		// 直接转发返回的结果
		out = in
	}
	return out
}

func newRequest(stack *hub.Stack, apiDef *hub.ApiDef, privateDef *hub.PrivateArray) *http.Request {
	var formBody *http.Request
	var outBody string
	var hasBody bool
	// 要发送的请求
	outReq, _ := http.NewRequest(apiDef.Method, "", nil)
	hasBody = len(apiDef.RequestContentType) > 0 && apiDef.RequestContentType != "none"
	if hasBody {
		switch apiDef.RequestContentType {
		case "form":
			outReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			formBody = new(http.Request)
			formBody.ParseForm()
		case "json":
			outReq.Header.Set("Content-Type", "application/json")
		case hub.OriginName:
			contentType := stack.GinContext.Request.Header.Get("Content-Type")
			outReq.Header.Set("Content-Type", contentType)
			// 收到的请求中的数据
			inData, _ := json.Marshal(stack.StepResult[hub.OriginName])
			outBody = string(inData)
		default:
			outReq.Header.Set("Content-Type", apiDef.RequestContentType)
		}
	}

	// 发出请求的URL
	var finalUrl string
	if len(apiDef.Url) == 0 {
		if apiDef.DynamicUrl != nil {
			finalUrl = unit.GetParameterValue(stack, privateDef, apiDef.DynamicUrl)
		} else {
			klog.Errorln("无有效url")
			panic("无有效url")
		}
	} else {
		finalUrl = apiDef.Url
	}
	outReqURL, _ := url.Parse(finalUrl)
	// 设置请求参数
	outReqParamRules := apiDef.Parameters
	if outReqParamRules != nil {
		paramLen := len(*outReqParamRules)
		if paramLen > 0 {
			var value string
			q := outReqURL.Query()
			vars := make(map[string]string, paramLen)
			stack.StepResult[hub.VarsName] = vars
			defer func() { stack.StepResult[hub.VarsName] = nil }()

			for _, param := range *outReqParamRules {
				if len(param.Name) > 0 {
					if len(param.Value) == 0 {
						if param.From != nil {
							value = unit.GetParameterValue(stack, privateDef, param.From)
						}
					} else {
						value = param.Value
					}

					switch param.In {
					case "query":
						q.Set(param.Name, value)
					case "header":
						outReq.Header.Set(param.Name, value)
					case "body":
						if hasBody && apiDef.RequestContentType != hub.OriginName {
							if apiDef.RequestContentType == "form" {
								formBody.Form.Add(param.Name, value)
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
							klog.Infoln("Refuse to set body :", apiDef.RequestContentType, "VS\r\n", value)
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

	outReq.URL = outReqURL

	// 处理要发送的消息体
	if apiDef.Method == "POST" {
		if apiDef.RequestContentType != "none" {
			if apiDef.RequestContentType == "form" {
				outBody = formBody.Form.Encode()
			}
			outReq.Body = ioutil.NopCloser(strings.NewReader(outBody))
		}
	}

	return outReq
}

func handleExpireTime(stack *hub.Stack, apiDef *hub.ApiDef, resp *http.Response) (time.Time, bool) {
	klog.Infoln("获得参数，[src]:", apiDef.Cache.From.From, "; [key]:", apiDef.Cache.From.Content, "; [format]:", apiDef.Cache.Format)
	if strings.EqualFold(apiDef.Cache.From.From, "header") {
		return handleHeaderExpireTime(apiDef, resp)
	} else {
		return handleBodyExpireTime(stack, apiDef)
	}
}

func handleHeaderExpireTime(apiDef *hub.ApiDef, resp *http.Response) (time.Time, bool) {
	//首先在api 的json文件中配置参数 cache
	// "cache": {
	// 	"from": {
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
	key := apiDef.Cache.From.Content
	format := apiDef.Cache.Format

	if strings.Contains(key, "Set-Cookie.") {
		key = strings.TrimPrefix(key, "Set-Cookie.")
		//判断Set-Cookie中是否含有Expires 的header
		cookie := resp.Header.Get("Set-Cookie")
		klog.Infoln("Header中Set-Cookie: ", cookie)
		if len(cookie) > 0 {
			expiresIndex := strings.Index(cookie, key) //"expires="
			if expiresIndex >= 0 {
				semicolonIndex := strings.Index(cookie[expiresIndex:], ";")
				if semicolonIndex < 0 {
					semicolonIndex = 0
				}

				expires, err := parseExpireTime(cookie[expiresIndex+len(key)+1:expiresIndex+semicolonIndex], format)
				if err == nil {
					return expires, true
				}
			}
		}
	} else {
		//判断是否含有Expires 的header
		expires, err := parseExpireTime(resp.Header.Get(key), format)
		if err == nil {
			return expires, true
		}
	}

	return time.Time{}, false
}

func handleBodyExpireTime(stack *hub.Stack, apiDef *hub.ApiDef) (time.Time, bool) {
	//首先在api 的json文件中配置参数 cache
	// "cache": {
	// 	"from": {
	// 		"from": "json",
	// 		"name": "{{.result.expires_in}}"
	// 	},
	// 	"format": "second"
	//   }
	//name 为获取过期时间的关键字串
	//format：如果是date格式，则配置具体格式串，如果是second数，则按照秒数解析
	//	baidu_image_classify_token: Mon, 02-Jan-06 15:04:05 MST
	//	body中一个例子："expireTime":"20220510153521",格式为：20060102150405

	format := apiDef.Cache.Format
	result := unit.GetParameterValue(stack, nil, apiDef.Cache.From)

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
		seconds := util.GetInterfaceToInt(value)
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

func getCacheContent(apiDef *hub.ApiDef) interface{} {
	//如果支持缓存，判断过期时间
	if time.Now().Local().After(apiDef.Cache.Expires) {
		return nil
	}
	return apiDef.Cache.Resp
}

func getCacheContentWithLock(apiDef *hub.ApiDef) interface{} {
	//如果支持缓存，判断过期时间
	apiDef.Cache.Locker.RLock()
	defer apiDef.Cache.Locker.RUnlock()
	if time.Now().Local().After(apiDef.Cache.Expires) {
		return nil
	}
	return apiDef.Cache.Resp
}

func handleRespStatus(stack *hub.Stack, apiDef *hub.ApiDef) bool {
	if apiDef.RespStatus == nil { //如果没有定义，则直接返回正确
		return true
	}

	result := unit.GetParameterValue(stack, nil, apiDef.RespStatus.From)
	klog.Infoln("handleRespStatus 结果", result)
	return apiDef.RespStatus.Expected == result
}

// 转发API调用
func Run(stack *hub.Stack, private string) (jsonOutRspBody interface{}, ret int) {
	var err error
	apiDef, err := unit.FindApiDef(stack, stack.ChildName)

	if apiDef == nil {
		klog.Errorln("获得API定义失败：", err)
		panic(err)
	}
	privateDef, err := unit.FindPrivateDef(stack, private, apiDef.PrivateName)
	if err != nil {
		klog.Errorln("获得API定义失败：", err)
		panic(err)
	}

	if apiDef.Cache != nil { //如果Json文件中配置了cache，表示支持缓存
		if jsonOutRspBody = getCacheContentWithLock(apiDef); jsonOutRspBody == nil {
			defer apiDef.Cache.Locker.Unlock()
			apiDef.Cache.Locker.Lock()

			if jsonOutRspBody = getCacheContent(apiDef); jsonOutRspBody == nil {
				klog.Infoln("获取缓存Cache ... ...")
				jsonOutRspBody, _ = handleReq(stack, apiDef, privateDef)
			} else {
				klog.Infoln("Cache缓存有效，直接回应")
			}
		} else {
			klog.Infoln("Cache缓存有效，直接回应")
		}
	} else { //不支持缓存，直接请求
		//klog.Infoln("不支持Cache缓存 ... ...")
		jsonOutRspBody, _ = handleReq(stack, apiDef, privateDef)
	}

	klog.Infoln("处理", apiDef.Url, ":", http.StatusOK, "\r\n返回结果：", jsonOutRspBody)
	if jsonOutRspBody == nil {
		return nil, 500
	}
	return jsonOutRspBody, http.StatusOK
}
