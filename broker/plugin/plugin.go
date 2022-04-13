package plugin

import (
	"log"
	"plugin"
)

func RewriteApiDef(pluginPath string) (func(interface{}), error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	// 导出函数变量
	RewriteApiDef, err := p.Lookup("RewriteApiDef")
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	return RewriteApiDef.(func(interface{})), nil
}
