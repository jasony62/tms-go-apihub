package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/unit"
	"github.com/jasony62/tms-go-apihub/util"
)

// 转发API调用
func Relay(stack *hub.Stack, resultKey string) (interface{}, int) {
	apiDef := stack.ApiDef
	var formBody *http.Request
	var outBody string
	var err error
	var hasBody bool

	// 要发送的请求
	outReq, _ := http.NewRequest(apiDef.Method, "", nil)
	hasBody = len(apiDef.RequestContentType) > 0 && apiDef.RequestContentType != "none"
	if hasBody {
		if apiDef.RequestContentType == "form" {
			outReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			formBody = new(http.Request)
			formBody.ParseForm()
		} else if apiDef.RequestContentType == "json" {
			outReq.Header.Set("Content-Type", "application/json")
		} else if apiDef.RequestContentType == "origin" {
			contentType := stack.GinContext.Request.Header.Get("Content-Type")
			outReq.Header.Set("Content-Type", contentType)
			// 收到的请求中的数据
			inData, _ := json.Marshal(stack.StepResult["origin"])
			outBody = string(inData)
		} else {
			outReq.Header.Set("Content-Type", apiDef.RequestContentType)
		}
	}

	// 发出请求的URL
	outReqURL, _ := url.Parse(apiDef.Url)

	// 设置请求参数
	outReqParamRules := apiDef.Parameters
	if outReqParamRules != nil && len(*outReqParamRules) > 0 {
		q := outReqURL.Query()
		for _, param := range *outReqParamRules {
			if len(param.Name) > 0 {
				if len(param.Value) == 0 {
					if param.From != nil {
						param.Value = unit.GetParameterValue(stack, param.From.From, param.From.Name, param.From.Template)
					}
				}

				if param.In == "query" {
					q.Set(param.Name, param.Value)
				} else if param.In == "header" {
					outReq.Header.Set(param.Name, param.Value)
				} else if param.In == "body" {
					if hasBody && apiDef.RequestContentType != "origin" {
						if apiDef.RequestContentType == "form" {
							formBody.Form.Add(param.Name, param.Value)
						} else {
							if len(outBody) == 0 {
								outBody = param.Value
							} else {
								log.Println("Double content body :\r\n", outBody, "\r\nVS\r\n", param.Value)
							}
						}
					} else {
						log.Println("Refuse to set body :", apiDef.RequestContentType, "VS\r\n", param.Value)
					}
				}
				log.Println("设置入参，位置", param.In, "名字", param.Name, "值", param.Value)
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
		fmt.Println("err", err)
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
		outRspBodyRules := apiDef.Response.Json
		jsonOutRspBody = util.Json2Json(jsonInRspBody, outRspBodyRules)
	} else {
		// 直接转发返回的结果
		jsonOutRspBody = jsonInRspBody
	}

	// 在上下文中保存结果
	if len(resultKey) > 0 {
		stack.StepResult[resultKey] = jsonOutRspBody
	}
	log.Println("处理", apiDef.Url, ":", http.StatusOK, "\r\n返回结果(原始)：", jsonInRspBody,
		"\r\n返回结果(修改后)：", jsonOutRspBody)
	return jsonOutRspBody, http.StatusOK
}
