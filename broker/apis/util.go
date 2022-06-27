package apis

import (
	"net/http"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"

	klog "k8s.io/klog/v2"
)

func loadConf(stack *hub.Stack, params map[string]string) (interface{}, int) {
	basePath := util.GetBasePath()
	if len(basePath) == 0 {
		return nil, http.StatusInternalServerError
	}

	util.LoadConf(basePath)

	return nil, 200
}

//downloadConf 解压zip包
func downloadConf(stack *hub.Stack, params map[string]string) (interface{}, int) {
	basePath := util.GetBasePath()
	if len(basePath) == 0 {
		return nil, http.StatusInternalServerError
	}

	remoteUrl := params["url"]
	klog.Infoln("DownloadConf remoteUrl:", remoteUrl)
	if len(remoteUrl) != 0 {
		password := params["password"]
		klog.Infoln("DownloadConf password:", password)
		if util.DownloadConf(remoteUrl, basePath, password) {
			klog.Infoln("Download conf zip package from remote url OK")
		} else {
			return nil, http.StatusInternalServerError
		}
	}
	return nil, 200
}

//DecompressZip 解压zip包
func decompressZip(stack *hub.Stack, params map[string]string) (interface{}, int) {
	basePath := util.GetBasePath()
	path := params["path"]
	if len(basePath) == 0 && len(path) == 0 {
		return nil, http.StatusInternalServerError
	}

	if len(path) > 0 {
		basePath = path
	}

	filename := params["file"]
	klog.Infoln("DecompressZip filename:", filename, " path:", basePath)
	if len(filename) != 0 {
		password := params["password"]
		klog.Infoln("DecompressZip password:", password)
		err := util.DeCompressZip(filename, basePath, password, nil, 0)
		if err != nil {
			klog.Errorln(err)
			return nil, http.StatusInternalServerError
		}
	}
	return nil, 200
}
