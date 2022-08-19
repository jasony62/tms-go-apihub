package apis

import (
	"net/http"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
	"go.uber.org/zap"
)

func loadConf(stack *hub.Stack, params map[string]string) (interface{}, int) {
	basePath := util.GetBasePath()
	if len(basePath) == 0 {
		str := "basePath is empty"
		zap.S().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	util.LoadConf(basePath)

	return nil, http.StatusOK
}

//downloadConf 解压zip包
func downloadConf(stack *hub.Stack, params map[string]string) (interface{}, int) {
	basePath := util.GetBasePath()
	if len(basePath) == 0 {
		str := "downloadConf base path is empty"
		zap.S().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	remoteUrl := params["url"]
	zap.S().Infoln("DownloadConf: ", stack.BaseString, " remoteUrl:", remoteUrl)
	if len(remoteUrl) != 0 {
		password := params["password"]
		if util.DownloadConf(remoteUrl, basePath, password) {
			zap.S().Infoln("Download conf OK:", remoteUrl)
		} else {
			str := "downloadConf conf failed"
			zap.S().Errorln(stack.BaseString, str)
			return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
		}
	}
	return nil, http.StatusOK
}

//DecompressZip 解压zip包
func decompressZip(stack *hub.Stack, params map[string]string) (interface{}, int) {
	basePath := util.GetBasePath()
	path := params["path"]
	if len(basePath) == 0 && len(path) == 0 {
		str := "decompressZip path is empty"
		zap.S().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	if len(path) > 0 {
		basePath = path
	}

	filename := params["file"]
	zap.S().Infoln("DecompressZip ", stack.BaseString, " filename:", filename, " path:", basePath)
	if len(filename) != 0 {
		password := params["password"]
		//		zap.S().Infoln("DecompressZip password:", password)
		err := util.DeCompressZip(filename, basePath, password, nil, 0)
		if err != nil {
			zap.S().Errorln(stack.BaseString, err)
			return util.CreateTmsError(hub.TmsErrorApisId, err.Error(), nil), http.StatusInternalServerError
		}
	}
	return nil, http.StatusOK
}
