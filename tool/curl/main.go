package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"postmaninternal"
	"strings"

	"k8s.io/klog/v2"
)

var curlPath string
var httpapis string
var apiHubHttpConf postmaninternal.ApiHubHttpConf

func init() {
	flag.StringVar(&curlPath, "curlfrom", "./curl/", "指定curl文件路径")
	flag.StringVar(&httpapis, "curlto", "./httpapis/", "指定httpaps文件路径")
}

/*
 * main function
 */
func main() {
	convertCurl(curlPath)
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
				if strings.Contains(lineStringArray[i], "curl") && (strings.Index(lineStringArray[i], "#") != 0) {
					pareCurlLine(lineStringArray[i])
				}
			}
		}
	}
	return nil
}

func pareCurlLine(line string) {
	if strings.Contains(line, "POST") { //POST
		apiHubHttpConf.Method = "POST"
	} else { // GET
		apiHubHttpConf.Method = "GET"
	}

	if strings.Contains(line, "-H") {
		hContent, hIndex := postmaninternal.GetStringBetweenSpecifySymbols(line[strings.Index(line, "-H"):], "\"", "\"")
		if (hIndex != -1) && (strings.Contains(hContent, "application/json")) {
			apiHubHttpConf.Requestcontenttype = "json"
		}
	}

	if strings.Contains(line, "-d") {
		hContent, hIndex := postmaninternal.GetStringBetweenSpecifySymbols(line[strings.Index(line, "-d"):], "{", "}")
		hContent = "{" + hContent + "}"
		if hIndex != -1 {

		}
	}
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
