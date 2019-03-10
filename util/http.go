package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// HTTPGet 带超时设置的请求一个url，单位: 秒
func HTTPGet(uri string, timeout int) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	resp, err := client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// HTTPPost 使用POST Form方式请求数据，超时 单位：秒
func HTTPPost(uri string, params map[string]string, timeout int, header map[string]string) ([]byte, error) {
	udata := make(url.Values)
	for k, v := range params {
		udata[k] = []string{v}
	}
	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	req, err := http.NewRequest("POST", uri, strings.NewReader(udata.Encode()))
	if err != nil {
		return nil, err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// HTTPPostJSON 使用POST JSON方式请求数据，超时 单位：秒
func HTTPPostJSON(uri string, params interface{}, timeout int) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	j, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(j)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	resp, err := client.Post(uri, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// HTTPPostFile 使用POST上传文件，超时 单位：秒
// 在参数中，如果要上传文件，设 params[file]=@/data/upload/1.zip 即可，注意 @ 符号
func HTTPPostFile(uri string, params map[string]string, cookie string, timeout int) ([]byte, http.Header, error) {
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	for key, val := range params {
		if len(val) > 0 && val[0] == '@' {
			file, err := os.Open(val[1:])
			if err != nil {
				return nil, nil, err
			}
			part, err := writer.CreateFormFile(key, filepath.Base(val[1:]))
			if err != nil {
				return nil, nil, err
			}
			_, err = io.Copy(part, file)
		} else {
			_ = writer.WriteField(key, val)
		}
	}
	err := writer.Close()
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest("POST", uri, reqBody)
	fmt.Println("FormDataContentType", writer.FormDataContentType())
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "Baa http agent")
	if len(cookie) > 0 {
		req.Header.Set("Cookie", cookie)
	}

	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = time.Second * time.Duration(timeout)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	var body []byte
	var header http.Header
	if resp.StatusCode == 200 {
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, nil, err
		}
		header = resp.Header
	}
	return body, header, nil
}

// HTTPDownload HTTP下载文件
func HTTPDownload(uri, path string) (int, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	return WriteFile(path, body)
}

// HTTPRangeDownload HTTP断点续传下载文件，提供下载进度，通过 range process得到进度条
func HTTPRangeDownload(uri, path string, process chan<- float64) (int64, error) {
	var offset, limit int64        // 初始偏移量
	var piece int64 = 65536        // 分片下载大小 64k
	var timeout time.Duration = 30 // 下载64k 最大超时30s

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return 0, fmt.Errorf("HTTPRangeDownload create request error: %v", err)
	}

	// 检查文件路径是否可写
	if err = MkdirAll(filepath.Dir(path)); err != nil {
		return 0, fmt.Errorf("HTTPRangeDownload check save path error: %v", err)
	}

	// 检查是否存在未下载完成的文件
	pieceFile := path + ".piece"
	if IsExist(pieceFile) {
		pieceStat, err := os.Stat(pieceFile)
		if err != nil {
			return 0, fmt.Errorf("HTTPRangeDownload check piece file error: %v", err)
		}
		offset = pieceStat.Size()
	}
	limit = offset + piece

	// 发起第一个分片请求，同时检查是否支持断点续传
	client := http.Client{
		Timeout: time.Second * timeout,
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", offset, limit))

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("HTTPRangeDownload first http request error: %v", err)
	}

	if resp.StatusCode != http.StatusPartialContent && resp.StatusCode == http.StatusOK {
		// 不支持断点续传，保存到文件
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, fmt.Errorf("HTTPRangeDownload first http read response error: %v", err)
		}
		resp.Body.Close()
		n, err := WriteFile(path, body)
		// 进度 100%
		go func() {
			if process != nil {
				process <- 100
			}
		}()
		if err != nil {
			return 0, fmt.Errorf("HTTPRangeDownload normal http download error: %v", err)
		}
		return int64(n), nil
	}

	if resp.StatusCode != http.StatusPartialContent {
		return 0, fmt.Errorf("HTTPRangeDownload got invalid response status code [%d], should be 206", resp.StatusCode)
	}

	// 分片下载，打开临时文件
	fh, err := os.OpenFile(pieceFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		return 0, fmt.Errorf("HTTPRangeDownload open piece file error: %v", err)
	}
	defer fh.Close()

	// 写入首次请求
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("HTTPRangeDownload http read response error: %v", err)
	}
	resp.Body.Close()
	fh.Write(body)

	// 获取文件总大小
	contentRange := strings.Split(resp.Header.Get("Content-Range"), "/")
	if len(contentRange) != 2 {
		return 0, fmt.Errorf("HTTPRangeDownload got invalid range response")
	}
	contentLength, _ := strconv.ParseInt(contentRange[1], 10, 64)
	// 修正首次切片值
	if limit > contentLength {
		limit = contentLength
	}
	// 输出进度
	if process != nil {
		go func(limit int64) {
			defer func() {
				// 防止无人读取，长时间占用，关闭通道后写入panic
				recover()
			}()
			process <- float64(limit) * 100 / float64(contentLength)
		}(limit)
	}

	// 分片下载
	for offset = limit; contentLength > offset; offset = limit {
		limit = offset + piece
		if limit > contentLength {
			limit = contentLength
		}
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", offset+1, limit))
		resp, err := client.Do(req)
		if err != nil {
			return offset, fmt.Errorf("HTTPRangeDownload range http request error: %v", err)
		}
		if resp.StatusCode != http.StatusPartialContent {
			return offset, fmt.Errorf("HTTPRangeDownload got range invalid response status code [%d], should be 206", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return offset, fmt.Errorf("HTTPRangeDownload range http read response error: %v", err)
		}
		resp.Body.Close()
		fh.Write(body)
		// 输出进度
		if process != nil {
			go func(limit int64) {
				defer func() {
					// 防止无人读取，长时间占用，关闭通道后写入panic
					recover()
				}()
				process <- float64(limit) * 100 / float64(contentLength)
			}(limit)
		}
	}

	// 下载完成，重命名文件
	fh.Close()
	err = os.Rename(pieceFile, path)
	if err != nil {
		return offset, fmt.Errorf("HTTPRangeDownload rename file error: %v", err)
	}

	// 防止无人读取，长时间占用，等待3秒后 关闭进度通道
	if process != nil {
		go func() {
			defer func() {
				// 防止通道关闭后再次关闭
				recover()
			}()
			time.Sleep(time.Second * 3)
			close(process)
		}()
	}

	return contentLength, nil
}
