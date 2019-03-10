package qihu360

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"git.code.tencent.com/xinhuameiyu/common/util"
	"github.com/go-baa/log"
)

// VideoBlockInfo 视频上传时第一次请求返回的数据
type VideoBlockInfo struct {
	ErrorCode    int               `json:"errorCode"`
	ErrorMessage string            `json:"errorMessage"`
	BlockHeads   []VideoBlockHeads `json:"block_heads"`
}

// VideoBlockHeads 视频分片信息
type VideoBlockHeads struct {
	IsLastBlock  bool   `json:"is_last_block"`
	NextBlockURL string `json:"next_block_url"`
	TID          string `json:"tid"`
	BID          int    `json:"bid"`
	Begin        int    `json:"begin"`
	End          int    `json:"end"`
}

type VideoBlockParams struct {
	IsInit   int    `json:"is_init"`
	Name     string `json:"name"`
	Size     string `json:"size"`
	LastTime string `json:"lastModifiedDate"`
}

// GetVideoBlock 视频上传时第一次请求数据
func GetVideoBlock(token, filename, size, lastTime string) (*VideoBlockInfo, *ErrorResult) {
	if token == "" || filename == "" || size == "" || lastTime == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}
	params := "is_init=1&name=" + filename + "&size=" + size + "&lastModifiedDate=" + lastTime

	req, err := http.NewRequest("POST", APIURL+"asset/jwt/video/upload?from=test", strings.NewReader(params))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if APIHost != "" {
		req.Header.Set("Host", APIHost)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	client.Timeout = time.Second * time.Duration(APIRequestTimeout)
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("获取360 视频分片信息 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "获取360 视频分片信息 失败"}
	}
	var body []byte
	if resp.StatusCode == 200 {
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Errorf("获取360 视频分片信息 失败：%s", err.Error())
			return nil, &ErrorResult{Message: "获取360 视频分片信息 失败"}
		}
	}
	ret := new(VideoBlockInfo)
	if err := json.Unmarshal(body, ret); err != nil {
		log.Errorf("解析360 视频分片信息 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "解析360 视频分片信息 响应失败"}
	}
	if ret.ErrorCode != 0 {
		return nil, &ErrorResult{Message: ret.ErrorMessage}
	}

	return ret, nil
}

// UploadVideoBlock 上传分片信息
func UploadVideoBlock(filePath, url string, begin, end int) *ErrorResult {
	fileByte, err := util.ReadFileBlock(filePath, begin, end+1)
	if err != nil {
		return &ErrorResult{Message: "缺少参数"}
	}
	fileData := string(fileByte)
	mimeBoundary := "----" + util.MD5(util.IntToString(int(time.Now().Unix())))
	var data []string
	data = append(data, "--"+mimeBoundary)
	data = append(data, "Content-Disposition: form-data; name=\"bhash\"")
	data = append(data, "")
	data = append(data, util.SHA1(fileData))
	data = append(data, "--"+mimeBoundary)
	data = append(data, "Content-Disposition: form-data; name=\"file\"; filename=\"block\"")
	data = append(data, "Content-Type: application/octet-stream")
	data = append(data, "")
	data = append(data, fileData)
	data = append(data, "--"+mimeBoundary+"--")
	bodyData := strings.Join(data, "\r\n")

	req, err := http.NewRequest("POST", url, strings.NewReader(bodyData))
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+mimeBoundary)
	req.Header.Set("Content-Length", util.IntToString(len(bodyData)))
	client := &http.Client{}
	client.Timeout = time.Second * time.Duration(APIRequestTimeout)
	resp, err := client.Do(req)
	if err != nil {
		return &ErrorResult{Message: err.Error()}
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(err)
	var body []byte
	if resp.StatusCode == 200 {
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return &ErrorResult{Message: err.Error()}
		}
	}
	ret := new(VideoBlockHeads)
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return &ErrorResult{Message: err.Error()}
	}
	// fmt.Println(string(body))
	// fmt.Println(ret)
	// fmt.Println(ret.End)
	if ret.IsLastBlock == false {
		UploadVideoBlock(filePath, ret.NextBlockURL, ret.Begin, ret.End)
	}

	return nil
}

// CommitVideoBlockInfo 合并分片信息返回的数据
type CommitVideoBlockInfo struct {
	Status     int             `json:"status"`
	StatusMsg  string          `json:"status_msg"`
	UploadID   string          `json:"upload_id"`
	URL        string          `json:"url"`
	UploadInfo VideoUploadInfo `json:"upload_info"`
}

// VideoUploadInfo 视频上传后返回的信息
type VideoUploadInfo struct {
	CdnURL    string    `json:"cdn_url"`
	VideoInfo VideoInfo `JSON:"video_info"`
	FileName  string    `json:"file_name"`
	Size      string    `json:"size"`
}

// VideoInfo 视频本身数据
type VideoInfo struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Duration  string `json:"duration"`
	CodecName string `json:"codec_name"`
}

// CommitVideoBlock 合并视频分片信息
func CommitVideoBlock(token, filename, size, lastTime string) (*CommitVideoBlockInfo, *ErrorResult) {
	if token == "" || filename == "" || size == "" || lastTime == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}

	params := "is_init=1&name=" + filename + "&size=" + size + "&lastModifiedDate=" + lastTime
	req, err := http.NewRequest("POST", APIURL+"asset/jwt/video/upload/commit?from=test", strings.NewReader(params))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if APIHost != "" {
		req.Header.Set("Host", APIHost)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		log.Errorf("合并360 视频分片信息 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "合并360 视频分片信息 失败"}
	}

	client := &http.Client{}
	client.Timeout = time.Second * time.Duration(APIRequestTimeout)
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("获取360 视频分片信息 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "获取360 视频分片信息 失败"}
	}
	var body []byte
	if resp.StatusCode == 200 {
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Errorf("获取360 视频分片信息 失败：%s", err.Error())
			return nil, &ErrorResult{Message: "获取360 视频分片信息 失败"}
		}
	}

	// fmt.Println(string(body))

	ret := new(CommitVideoBlockInfo)
	if err := json.Unmarshal(body, ret); err != nil {
		log.Errorf("合并360 视频分片信息 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "合并360 视频分片信息 响应失败"}
	}
	if ret.Status != 0 {
		return nil, &ErrorResult{Message: ret.StatusMsg}
	}

	return ret, nil
}

// ImageFileInfo 上传图片文件返回的图片地址
type ImageFileInfo struct {
	Common string `json:"common"`
}

// UploadImageFile 上传图片文件
func UploadImageFile(token, filePath string) (*ImageFileInfo, *ErrorResult) {
	if token == "" || filePath == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}

	params := "common=" + filePath

	req, err := http.NewRequest("POST", APIURL+"asset/jwt/upload", strings.NewReader(params))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if APIHost != "" {
		req.Header.Set("Host", APIHost)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	client.Timeout = time.Second * time.Duration(APIRequestTimeout)
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("360 上传图片文件 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "360 上传图片文件 失败"}
	}
	var body []byte
	if resp.StatusCode == 200 {
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Errorf("360 上传图片文件 失败：%s", err.Error())
			return nil, &ErrorResult{Message: "360 上传图片文件 失败"}
		}
	}

	ret := new(ImageFileInfo)
	if err := json.Unmarshal(body, ret); err != nil {
		log.Errorf("解析360 上传图片文件 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "解析360 上传图片文件 响应失败"}
	}

	return ret, nil
}

// ImageURLInfo 使用图片地址上传图片 返回数据
type ImageURLInfo struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    ImageURL `json:"data"`
}

// ImageURL 上传完成后的图片地址
type ImageURL struct {
	URL string `json:"url"`
}

// UploadImageURL 上传图片地址
func UploadImageURL(token, url string) (*ImageURLInfo, *ErrorResult) {
	if token == "" || url == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}

	req, err := http.NewRequest("GET", APIURL+"asset/jwt/image/upload?url="+url, nil)
	if APIHost != "" {
		req.Header.Set("Host", APIHost)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	client.Timeout = time.Second * time.Duration(APIRequestTimeout)
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("360 上传图片地址 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "360 上传图片地址 失败"}
	}
	var body []byte
	if resp.StatusCode == 200 {
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Errorf("360 上传图片地址 失败：%s", err.Error())
			return nil, &ErrorResult{Message: "360 上传图片地址 失败"}
		}
	}

	ret := new(ImageURLInfo)
	if err := json.Unmarshal(body, ret); err != nil {
		log.Errorf("解析360 上传图片地址 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "解析360 上传图片地址 响应失败"}
	}

	return ret, nil
}

// Content 360 文章发布数据
type Content struct {
	Title      string `json:"title"`
	Brief      string `json:"brief"`
	Tag        string `json:"tag"`
	Cat        string `json:"cat"`
	Thumb      string `json:"thumb"`
	VideoSize  string `json:"video_size"`
	UploadID   string `json:"upload_id"`
	SearchWord string `json:"search_word"`
	Summary    string `json:"summary"`
}

// ContentStep 360 文章发布分步骤内容
type ContentStep struct {
	Text string `json:"text"`
	Img  string `json:"img"`
}

// ContentRespInfo 文章发布返回的数据
type ContentRespInfo struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    ContentRespData `json:"data"`
}

// ContentRespData 文章发布返回的ID
type ContentRespData struct {
	ID string `json:"id"`
}

// PushContentDraft 向360发布文章 保存为草稿
func PushContentDraft(token string, row Content) (*ContentRespInfo, *ErrorResult) {
	if token == "" {
		return nil, &ErrorResult{Message: "参数不正确"}
	}
	params := url.Values{}
	params.Add("title", row.Title)
	params.Add("brief", row.Brief)
	params.Add("tag", row.Tag)
	params.Add("cat", row.Cat)
	params.Add("thumb", row.Thumb)
	params.Add("video_size", row.VideoSize)
	params.Add("upload_id", row.UploadID)
	params.Add("search_word", row.SearchWord)
	if row.Summary != "" {
		params.Add("summary", row.Summary)
	}
	fmt.Println(params)

	req, err := http.NewRequest("POST", APIURL+"mgrvideo/jwt/save", strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if APIHost != "" {
		req.Header.Set("Host", APIHost)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	client.Timeout = time.Second * time.Duration(APIRequestTimeout)
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("360 发布文章 保存草稿 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "360 发布文章 保存草稿 失败"}
	}
	var body []byte
	if resp.StatusCode == 200 {
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Errorf("360 发布文章 保存草稿 失败：%s", err.Error())
			return nil, &ErrorResult{Message: "360 发布文章 保存草稿 失败"}
		}
	}
	fmt.Println(string(body))

	ret := new(ContentRespInfo)
	if err := json.Unmarshal(body, ret); err != nil {
		log.Errorf("360 发布文章 保存草稿 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "360 发布文章 保存草稿 响应失败"}
	}
	if ret.Status != 0 {
		return nil, &ErrorResult{Message: "发布文章 保存草稿失败"}
	}

	return ret, nil
}

// PublishContentResp 发布文章后返回的数据
type PublishContentResp struct {
	Status int `json:"status"`
}

// PublishContentDraft 发布文章草稿
func PublishContentDraft(token string, id string) (*PublishContentResp, *ErrorResult) {
	if token == "" {
		return nil, &ErrorResult{Message: "参数不正确"}
	}
	params := "id=" + id

	req, err := http.NewRequest("POST", APIURL+"mgrvideo/jwt/publish", strings.NewReader(params))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if APIHost != "" {
		req.Header.Set("Host", APIHost)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	client.Timeout = time.Second * time.Duration(APIRequestTimeout)
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("360 发布文章草稿 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "合并360 发布文章草稿 失败"}
	}
	var body []byte
	if resp.StatusCode == 200 {
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Errorf("360 发布文章草稿 失败：%s", err.Error())
			return nil, &ErrorResult{Message: "合并360 发布文章草稿 失败"}
		}
	}
	fmt.Println(string(body))

	ret := new(PublishContentResp)
	if err := json.Unmarshal(body, ret); err != nil {
		log.Errorf("360 发布文章草稿 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "360 发布文章草稿 响应失败"}
	}
	if ret.Status != 0 {
		return nil, &ErrorResult{Message: "发布文章草稿 失败"}
	}

	return ret, nil
}

// buildAPIRequestURL 构建请求参数
func buildAPIRequestURL(gateway string, params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	query := values.Encode()
	if len(query) > 0 {
		if strings.Contains(gateway, "?") {
			return gateway + "&" + query
		}
		return gateway + "?" + query
	}

	return gateway
}
