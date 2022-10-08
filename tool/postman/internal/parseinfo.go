package postmaninternal

import (
	"github.com/rbretecher/go-postman-collection"
	"k8s.io/klog/v2"
)

func getHttpapiInfo(postmanItem *postman.Items) {
	if postmanItem == nil {
		return
	}

	apiHubHttpConf.Description = postmanItem.Name
	klog.Infoln("__request Description : ", apiHubHttpConf.Description)
	apiHubHttpConf.URL = getPostmanURL(postmanItem.Request.URL)
	klog.Infoln("__request URL : ", apiHubHttpConf.URL)
	apiHubHttpConf.Method = string(postmanItem.Request.Method)
	klog.Infoln("__request Method : ", apiHubHttpConf.Method)

	apiHubHttpConf.Private = ""                // default private content
	apiHubHttpConf.Requestcontenttype = "json" //default content typeStr
}

// 获取Request URL
func getPostmanURL(postmanUrl *postman.URL) string {
	if postmanUrl == nil {
		return ""
	}

	httpapiUrl := postmanUrl.Protocol + "://"
	// Host IP
	if postmanUrl.Host != nil {
		for i := range postmanUrl.Host {
			if i > 0 {
				httpapiUrl = httpapiUrl + "." + postmanUrl.Host[i]
			} else {
				httpapiUrl = httpapiUrl + postmanUrl.Host[i]
			}
		}
	} else {
		klog.Infoln("__getPostmanURL Error, url不符合规范")
	}

	// Port number
	if postmanUrl.Port != "" {
		httpapiUrl = httpapiUrl + ":" + postmanUrl.Port + "/"
	} else {
		httpapiUrl = httpapiUrl + "/"
	}
	// Path
	for i := range postmanUrl.Path {
		if postmanUrl.Path[i] != "" {
			if i != (len(postmanUrl.Path) - 1) {
				httpapiUrl = httpapiUrl + postmanUrl.Path[i] + "/"
			} else {
				httpapiUrl = httpapiUrl + postmanUrl.Path[i]
			}
		}
	}
	return httpapiUrl
}
