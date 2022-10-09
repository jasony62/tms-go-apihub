package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"postmaninternal"
	"strconv"
	"strings"

	"k8s.io/klog/v2"
)

var curl_path string
var httpapis_path string
var apiHubHttpConf postmaninternal.ApiHubHttpConf

func init() {
	flag.StringVar(&curl_path, "curlfrom", "./curl/", "指定curl文件路径")
	flag.StringVar(&httpapis_path, "curlto", "./httpapis/", "指定httpaps文件路径")
}

/*
 * main function
 */
func main() {
	convertCurl(curl_path)
}

func convertCurl(path string) error {
	fileInfoList, err := os.ReadDir(path)
	if err != nil {
		klog.Errorln(err)
		return err
	}
	var prefix string
	// 遍历curl文件
	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())
		// 是否是一个子目录 若是子目录，进入进入子目录遍历curl文件
		if fileInfoList[i].IsDir() {
			klog.Infoln("__curl 子目录名: ", fileName)
			prefix = fileInfoList[i].Name()
			convertCurl(path + "/" + prefix)
		} else {
			// 判断文件是否curl类型
			if !strings.HasSuffix(fileName, ".txt") {
				continue
			}
			klog.Infoln("__Load curl 文件: ", fileName)

			lineStringArray, err := lineByLine(fileName)
			if err != nil {
				klog.Errorln(err)
				return err
			}

			for i := 0; i < len(lineStringArray); i++ {
				apiHubHttpConf.Args = append(apiHubHttpConf.Args[:0], apiHubHttpConf.Args[len(apiHubHttpConf.Args):]...)
				if strings.Contains(lineStringArray[i], "curl") && (strings.Index(lineStringArray[i], "#") != 0) {
					parseCurlCmd(lineStringArray[i])
					apiHubHttpConf.ID = "Curl_Line_Number" + strconv.Itoa(i)
					apiHubHttpConf.Description = "Curl_Line_Number" + strconv.Itoa(i)
					generateApiHubJson(httpapis_path, "Curl_Line_Number"+strconv.Itoa(i))
				}
			}
		}
	}
	return nil
}

func generateApiHubJson(apiHubJsonPath string, multipleName string) {
	var byteHttpApi []byte
	var err error
	fileName := apiHubJsonPath + multipleName + ".json"

	byteHttpApi, err = json.Marshal(apiHubHttpConf)
	if err != nil {
		klog.Errorln("json.Marshal失败!", fileName)
		return
	}

	// ！！！os.Create无法自动创建文件路径中不存在的文件夹
	f, err := os.Create(fileName)
	if err != nil {
		klog.Errorln("创建文件失败!", fileName)
	} else {
		defer f.Close()
		_, err = f.Write(byteHttpApi)
		if err != nil {
			klog.Errorln("写入文件失败!", fileName)
		}
	}
}

func parseCurlCmd(line string) {

	apiHubHttpConf.Method = "POST"

	if strings.Contains(line, "-d") {
		hContent, hIndex := postmaninternal.GetStringBetweenSpecifySymbols(line[strings.Index(line, "-d"):], "{", "}")
		if hIndex != -1 {
			hContent = "{" + hContent + "}"
			dataMap := make(map[string]string)
			err := json.Unmarshal([]byte(hContent), &dataMap)
			if err != nil {
				fmt.Println("Umarshal failed:", err)
				return
			}
			args := postmaninternal.Args{In: "body", Name: "body", Value: postmaninternal.Value{From: "json", Json: dataMap}}
			apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
		}
	}

	if strings.Contains(line, "-H") {
		hContent, hIndex := postmaninternal.GetStringBetweenSpecifySymbols(line[strings.Index(line, "-H"):], "\"", "\"")
		if (hIndex != -1) && (strings.Contains(hContent, "application/json")) {
			apiHubHttpConf.Requestcontenttype = "json"
		} else {
			apiHubHttpConf.Requestcontenttype = "json" // default to application/json
		}
	}

	if strings.Contains(line, "\"http") {
		apiHubHttpConf.URL = line[strings.Index(line, "\"http"):]
		hContent, hIndex := postmaninternal.GetStringBetweenSpecifySymbols(apiHubHttpConf.URL, "\"", "\"")
		if hIndex != -1 {
			// apiHubHttpConf.URL = strings.Replace(hContent, "\u0026", "&", -1)
			apiHubHttpConf.URL = hContent
		}
		/*
			if hIndex != -1 {
				if strings.Contains(hContent, "?") {
					apiHubHttpConf.URL = hContent[:strings.Index(hContent, "?")]
					if strings.Contains(hContent, "&") {
						tempflag := 0
						for strings.Contains(hContent, "&") {
							if tempflag == 0 {
								tempQuery, tempIndex := postmaninternal.GetStringBetweenSpecifySymbols(hContent, "?", "&")
								tempflag = 1
								tempkey, tempvalue := parseGetContent(tempQuery)
								args := postmaninternal.Args{In: "query", Name: tempkey, Value: postmaninternal.Value{From: "lietral", Content: tempvalue}}
								apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
								hContent = hContent[tempIndex:]
							} else {
								tempQuery, tempIndex := postmaninternal.GetStringBetweenSpecifySymbols(hContent, "&", "&")
								tempkey, tempvalue := parseGetContent(tempQuery)
								args := postmaninternal.Args{In: "query", Name: tempkey, Value: postmaninternal.Value{From: "lietral", Content: tempvalue}}
								apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
								if tempIndex == -1 {
									hContent = hContent[1:]
									break
								}
								hContent = hContent[tempIndex:]
							}
						}
						_ = hContent
					} else {
						tempQuery, tempIndex := postmaninternal.GetStringBetweenSpecifySymbols(hContent, "?", "\"")
						if tempIndex != -1 {
							tempkey, tempvalue := parseGetContent(tempQuery)
							args := postmaninternal.Args{In: "query", Name: tempkey, Value: postmaninternal.Value{From: "lietral", Content: tempvalue}}
							apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
						}
					}
				} else {
					apiHubHttpConf.URL = hContent
				}
			}
		*/
	}

}

func parseGetContent(content string) (string, string) {
	if strings.Contains(content, "=") {
		key := content[:strings.Index(content, "=")]
		value := content[strings.Index(content, "=")+1:]
		return key, value
	}
	return "", ""
}

func lineByLine(file string) ([]string, error) {

	var err error

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	var lineStringArray []string
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error reading file %s", err)
			break
		}
		lineStringArray = append(lineStringArray, line)
		// fmt.Print(line)
	}
	return lineStringArray, nil
}
