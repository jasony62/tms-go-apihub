package plugin

import (
	klog "k8s.io/klog/v2"
	"plugin"
)

func RewriteApiDef(pluginPath string) (func(interface{}), error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		klog.Errorln(err)
		panic(err)
	}

	// 导出函数变量
	RewriteApiDef, err := p.Lookup("RewriteApiDef")
	if err != nil {
		klog.Errorln(err)
		panic(err)
	}

	return RewriteApiDef.(func(interface{})), nil
}
