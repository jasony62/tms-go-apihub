package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/unit"
	"github.com/jasony62/tms-go-apihub/util"
)

// 转发API调用
func Run(stack *hub.Stack) (interface{}, int) {
	var formBody *http.Request
	var outBody string
	var err error
	var hasBody bool
	apiDef, err := unit.FindApiDef(stack, stack.Name)

	if apiDef == nil {
		klog.Errorln("获得API定义失败：", err)
		panic(err)
	}

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
	paramLen := len(*outReqParamRules)
	if outReqParamRules != nil && paramLen > 0 {
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
	var jsonInRspBody interface{}
	json.Unmarshal(returnBody, &jsonInRspBody)
	var jsonOutRspBody interface{}

	// 构造发送的响应内容
	if apiDef.Response != nil && apiDef.Response.Json != nil {
		jsonOutRspBody = util.Json2Json(jsonInRspBody, apiDef.Response.Json)
	} else {
		// 直接转发返回的结果
		jsonOutRspBody = jsonInRspBody
	}

	klog.Infoln("处理", apiDef.Url, ":", http.StatusOK, "\r\n返回结果：", jsonOutRspBody)
	return jsonOutRspBody, http.StatusOK
}
