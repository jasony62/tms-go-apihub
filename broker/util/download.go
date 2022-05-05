package util

import (
	"bufio"
	"errors"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	klog "k8s.io/klog/v2"
)

const (
	FILE_READ_MAX_SIZE int = 1024 * 32
)

func DownloadFile(fileUrl string) (string, error) {

	klog.Infoln("DownloadFile url: ", fileUrl)
	url, err := url.ParseRequestURI(fileUrl)
	if err != nil {
		klog.Errorln("DownloadFile url invalid: ", err)
		return "", err
	}

	filename := path.Base(url.Path)
	klog.Infoln("DownloadFile name: ", filename)

	client := http.DefaultClient
	client.Timeout = time.Second * 600

	resp, err := client.Get(fileUrl)
	if err != nil {
		klog.Errorln(err)
		return "", err
	}

	if resp.ContentLength <= 0 {
		klog.Errorln("DownloadFile: server response length error")
		return "", errors.New("DownloadFile: server response length error")
	}

	content := resp.Body
	defer content.Close()

	written := 0
	go copyFileContent(content, filename, &written)

	spaceTime := time.Second * 1
	ticker := time.NewTicker(spaceTime)
	lastWtn := 0
	stop := false

	for {
		select {
		case <-ticker.C:
			speed := written - lastWtn
			klog.Infof("[DownloadFile] Speed %s / %s \n", bytesToSize(speed), spaceTime.String())
			if written-lastWtn == 0 {
				ticker.Stop()
				stop = true
				break
			}
			lastWtn = written
		}
		if stop {
			break
		}
	}

	klog.Infoln("DownloadFile OK: ", filename)
	return filename, nil
}

//下载远端url文件
func copyFileContent(raw io.ReadCloser, filename string, written *int) error {
	klog.Infoln("Download url file starting!")
	reader := bufio.NewReaderSize(raw, FILE_READ_MAX_SIZE)

	file, err := os.Create(filename)
	if err != nil {
		klog.Errorln("copyFileContent create file error:", err)
		return err
	}
	writer := bufio.NewWriter(file)
	buff := make([]byte, FILE_READ_MAX_SIZE)

	for {
		nReader, errReader := reader.Read(buff)
		if nReader > 0 {
			nWriter, errWriter := writer.Write(buff[0:nReader])
			if nWriter > 0 {
				*written += nWriter
			}
			if errWriter != nil {
				err = errWriter
				break
			}
			if nReader != nWriter {
				err = io.ErrShortWrite
				break
			}
		}

		if errReader != nil {
			if errReader != io.EOF {
				err = errReader
			}
			break
		}
	}

	if err != nil {
		klog.Errorln("copyFileContent read or write error:", err)
		return err
	}

	return nil
}

func bytesToSize(length int) string {
	var k = 1024
	var sizes = []string{"Bytes", "KB", "MB", "GB", "TB"}
	if length == 0 {
		return "0 Bytes"
	}
	i := math.Floor(math.Log(float64(length)) / math.Log(float64(k)))
	r := float64(length) / math.Pow(float64(k), i)
	return strconv.FormatFloat(r, 'f', 3, 64) + " " + sizes[int(i)]
}
