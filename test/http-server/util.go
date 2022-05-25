package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"time"

	cmAuth "github.com/chartmuseum/auth"
	klog "k8s.io/klog/v2"
)

func md5Template(args ...interface{}) string {

	if len(args) == 0 {
		return ""
	}
	str := fmt.Sprint(args...)
	w := md5.New()
	io.WriteString(w, str)
	checksum := fmt.Sprintf("%x", w.Sum(nil))
	return checksum
}

func checkSumValid(key string, reqBody *RegistEntity) bool {
	checkSum := md5Template(reqBody.Nonce, reqBody.Utc, key)
	klog.Infoln("calculated checkSum is: ", checkSum)
	if checkSum != reqBody.Checksum {
		klog.Errorln("CheckSum failed!")
		return false
	}
	klog.Errorln("CheckSum succeed!")
	return true
}
func checkUtcTimeValid(curTime int64, utc string) bool {
	reqUtc, err := strconv.ParseInt(utc, 10, 64)
	if err != nil {
		klog.Errorln("Request UTC time convert int64 error!")
		return false
	} else if reqUtc < curTime-300 {
		klog.Warningln("Request UTC time is outOfDate")
		return false
	} else if reqUtc > curTime {
		klog.Warningln("Request UTC time is more than current time!")
		return false
	}
	klog.Infoln("UTC time is valid!")
	return true
}

func readConfig(fileName string, v interface{}) {
	klog.Infoln("Load file: ", fileName)
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		klog.Errorln("readConfigFile Failed!")
	}
	klog.Infoln("String is ", string(data))
	err = json.Unmarshal(data, v)
	if err != nil {
		klog.Infoln(err)
		return
	}
}

func generateToken(urlPath []byte, action []byte, expires int64) string {
	cmTokenGenerator, err := cmAuth.NewTokenGenerator(&cmAuth.TokenGeneratorOptions{
		PrivateKeyPath: "./server.key",
	})
	if err != nil {
		panic(err)
	}
	access := []cmAuth.AccessEntry{
		{
			Name:    string(urlPath),
			Type:    cmAuth.AccessEntryType,
			Actions: []string{string(action)},
		},
	}
	signedString, err := cmTokenGenerator.GenerateToken(access, time.Duration(expires))
	if err != nil {
		panic(err)
	}
	// Prints a JWT token which you can use to make requests to ChartMuseum
	return signedString
}
