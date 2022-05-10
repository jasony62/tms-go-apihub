package api

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/unit"
	"github.com/jasony62/tms-go-apihub/util"
)

// 转发API调用
func Run(stack *hub.Stack) (interface{}, int) {
	var err error
	apiDef, err := unit.FindApiDef(stack, stack.Name)

	if apiDef == nil {
		klog.Errorln("获得API定义失败：", err)
		panic(err)
	}

	var jsonInRspBody interface{}
	var jsonOutRspBody interface{}
	var jsonInRspCacheBody interface{}

	apiDef.Token.Lock.RLock()
	TokenExp := TokenExpired(apiDef)
	apiDef.Token.Lock.RUnlock()

	if !TokenExp {
		klog.Infoln("Token缓存过期或还未获取缓存，获取Token")
		defer apiDef.Token.Lock.Unlock()

		apiDef.Token.Lock.Lock()
		if !TokenExpired(apiDef) {
			outReq := NewRequest(stack, apiDef)

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
			json.Unmarshal(returnBody, &jsonInRspBody)

			//解析过期时间，如果存在则记录下来
			//str := `{"msg":"鉴权成功","expireTime":"600","ak":"MTY1MjEMTAwMU1UWTFNakUyT0RFeU1UUTNNeU14TURBd01USTJNQT09","resultcode":"1"}`
			//str := `{"msg":"鉴权成功","expireTime":"20220510153521","ak":"MTY1MjEMTAwMU1UWTFNakUyT0RFeU1UUTNNeU14TURBd01USTJNQT09","resultcode":"1"}`
			//expires, ok := HandleExpires(stack, resp, str, apiDef)
			expires, ok := HandleExpires(stack, resp, string(returnBody), apiDef)
			if !ok {
				klog.Warningln("没有查询到过期时间")
				jsonInRspCacheBody = jsonInRspBody
				apiDef.Token.Cache = false
			} else {
				klog.Infof("更新Token信息，过期时间为: %v", expires)
				apiDef.Token.Cache = true
				apiDef.Token.Expires = expires
				apiDef.Token.Resp = jsonInRspBody
				jsonInRspCacheBody = apiDef.Token.Resp
			}
		} else {
			klog.Infoln("Token缓存有效，直接回应")
			jsonInRspCacheBody = apiDef.Token.Resp
		}
	} else {
		klog.Infoln("Token缓存有效，直接回应")
		apiDef.Token.Lock.RLock()
		jsonInRspCacheBody = apiDef.Token.Resp
		apiDef.Token.Lock.RUnlock()
	}

	// 构造发送的响应内容
	if apiDef.Response != nil && apiDef.Response.Json != nil {
		jsonOutRspBody = util.Json2Json(jsonInRspCacheBody, apiDef.Response.Json)
	} else {
		// 直接转发返回的结果
		jsonOutRspBody = jsonInRspCacheBody
	}

	klog.Infoln("处理", apiDef.Url, ":", http.StatusOK, "\r\n返回结果：", jsonOutRspBody)
	return jsonOutRspBody, http.StatusOK
}

func NewRequest(stack *hub.Stack, apiDef *hub.ApiDef) *http.Request {
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
	outReqURL, _ := url.Parse(apiDef.Url)

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
							value = unit.GetParameterValue(stack, apiDef.Privates, param.From)
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
									outBody = value
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

func HandleExpires(stack *hub.Stack, resp *http.Response, body string, apiDef *hub.ApiDef) (time.Time, bool) {
	//首先在api 的json文件中配置参数
	//expire_key：from 为从header还是从body中获取过期时间
	//            name 为获取过期时间的关键字串
	//expire_unit：from 目前有两种类型，date格式和second数
	//            name 如果是date格式，则配置具体格式串，如果是second数，则name会忽略
	var src, key, unit, format string
	outReqParamRules := apiDef.Parameters
	if outReqParamRules != nil {
		paramLen := len(*outReqParamRules)
		if paramLen > 0 {
			for _, param := range *outReqParamRules {
				if len(param.Name) > 0 {
					if param.Name == "expire_key" {
						src = param.From.From
						key = param.From.Name
					} else if param.Name == "expire_unit" {
						unit = param.From.From
						format = param.From.Name
					}
				}
			}
		}
	}

	klog.Infoln("获得参数，[src]:", src, "; [key]:", key, "; [unit]:", unit, "; [format]:", format)

	if src == "header" {
		if key == "expires" {
			//判断Set-Cookie中是否含有Expires 的header
			cookie := resp.Header.Get("Set-Cookie")
			klog.Infoln("Header中Set-Cookie: ", cookie)
			if len(cookie) > 0 {
				expiresIndex := strings.Index(cookie, key) //"expires="
				if expiresIndex >= 0 {
					semicolonIndex := strings.Index(cookie[expiresIndex:], ";")
					expiresStr := cookie[expiresIndex+8 : expiresIndex+semicolonIndex]
					expires, err := ParseExpires(expiresStr, unit, format)
					if err == nil {
						klog.Infoln("消息头中Set-Cookie的过期时间: ", expires)
						// klog.Infoln("Mock now: ", time.Now())
						// klog.Infoln("Mock Expires: ", time.Now().Add(60000000000))
						// return time.Now().Add(60000000000), true
						return expires, true
					}
				}
			}
		} else if key == "Expires" {
			//判断是否含有Expires 的header
			expireHeader := resp.Header.Get(key)
			klog.Infoln("Header中Expires: ", expireHeader)
			if len(expireHeader) > 0 {
				expires, err := ParseExpires(expireHeader, unit, format)
				if err == nil {
					klog.Infoln("消息头中Expires的过期时间: ", expires)
					return expires, true
				}
			}
		}

	} else if src == "body" {
		index := strings.Index(body, key)
		if index >= 0 {
			comma := strings.Index(body[index:], ",")
			if comma >= 0 {
				timestr := body[index+12 : index+comma]

				timestr = strings.TrimPrefix(timestr, `"`)
				timestr = strings.TrimSuffix(timestr, `"`)
				klog.Infoln("消息体中过期时间:", timestr)
				formatTime, err := ParseExpires(timestr, unit, format)
				if err == nil {
					klog.Infoln("消息体解析后过期时间:", formatTime)
					return formatTime, true
				}
			}
		}
	} else {
		klog.Warningln("Json文件中未配置过期时间的来源: ", src)
	}

	return time.Time{}, false
}

func ParseExpires(str string, unit string, format string) (time.Time, error) {
	var exptime time.Time
	var err error
	if unit == "second" {
		var s int
		s, err = strconv.Atoi(str)
		if err != nil {
			klog.Errorln("ParseExpires秒数失败")
			return time.Time{}, errors.New("Parse expires failed")
		}

		exptime = time.Now()
		exptime = exptime.Add(time.Second * time.Duration(s))
		klog.Infoln("按照second解析Expires: ", exptime)
	} else if unit == "date" {
		exptime, err = time.Parse(format, str)
		if err != nil {
			klog.Infoln("解析Expires失败:")
			return time.Time{}, errors.New("Parse expires failed")
		}
		klog.Infoln("按照date解析Expires: ", exptime)
	}

	return exptime.Local(), nil
}

func TokenExpired(apiDef *hub.ApiDef) bool {
	//判断是否支持缓存，如果支持，并且没有过期，则直接返回结果，
	//如果不支持缓存，则请求，并且判断过期时间，并修改缓存标志
	//如果失效了（过期了），则加锁然后获取新的过期时间
	return apiDef.Token.Cache && time.Now().Local().Before(apiDef.Token.Expires)
}
