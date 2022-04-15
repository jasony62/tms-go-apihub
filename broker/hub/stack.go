package hub

import (
	"bytes"
	"html/template"

	"github.com/gin-gonic/gin"
)

type Stack struct {
	BucketEnable   bool
	ApiDefPath     string
	PrivateDefPath string
	FlowDefPath    string
	ApiDef         *ApiDef
	FlowDef        *FlowDef
	GinContext     *gin.Context
	RequestBody    *interface{}
	StepResult     map[string]interface{}
	CurrentStep    *FlowStepDef
}

// 从执行结果中获取查询参数
func (stack Stack) queryValueFromResult(name string) string {
	tmpl, _ := template.New("key").Parse(name)
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, stack.StepResult)

	return buf.String()
}

// 获取查询参数
func (stack Stack) Query(name string) string {
	/*检查当前步骤是否指定了参数获取规则*/
	if stack.CurrentStep != nil {
		currentTask := stack.CurrentStep
		params := currentTask.Api.Parameters
		for _, param := range params {
			if param.Name == name {
				if param.From.In == `StepResult` {
					v := stack.queryValueFromResult(param.From.Name)
					return v
				}
				break
			}
		}
	}

	/*默认从请求的查询参数中获得*/
	return stack.GinContext.Query(name)
}
