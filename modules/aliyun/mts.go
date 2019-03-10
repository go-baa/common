package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/go-baa/common/util"
)

type MTSQueryMediaListResponse struct {
	Message          string
	Code             string
	NonExistMediaIds *MTSMediaId
	MediaList        *MTSMediaList
	RequestId        string
}

type MTSMediaId struct {
	MediaId []string
}

type MTSMediaList struct {
	Media []*MTSMedia
}

type MTSMedia struct {
	CoverURL     string
	Format       string
	PublishState string
	File         *MTSMediaFile
	MediaId      string
	Width        string
	Height       string
	Fps          string
	Bitrate      string
	Size         string
	Duration     string
	Title        string
	CateId       int
	CreationTime string
	MediaInfo    *MTSMediaInfo
	PlayList     *MTSPlayList
}

type MTSPlayList struct {
	Play []*MTSPlay
}

type MTSPlay struct {
	ActivityName string
	Size         string
	File         MTSMediaFile
}

type MTSMediaFile struct {
	State string
	URL   string
}

type MTSMediaInfo struct {
	Format  map[string]interface{}
	Streams map[string]interface{}
}

type MTS struct {
	accessKeyId     string
	accessKeySecret string
}

func (t *MTS) getCommonParams() map[string]string {
	return map[string]string{
		"Format":           "JSON",
		"Version":          "2014-06-18",
		"AccessKeyId":      t.accessKeyId,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   string(util.RandStr(16, util.KC_RAND_KIND_ALL)),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	}
}

func (t *MTS) QueryMediaList(location string, mediaIds []string, includePlayList bool, includeSnapshotList bool, includeMediaInfo bool) (*MTSQueryMediaListResponse, error) {
	params := t.getCommonParams()
	params["Action"] = "QueryMediaList"
	params["MediaIds"] = strings.Join(mediaIds, ",")

	if includePlayList {
		params["IncludePlayList"] = "true"
	}
	if includeSnapshotList {
		params["IncludeSnapshotList"] = "true"
	}
	if includeMediaInfo {
		params["IncludeMediaInfo"] = "true"
	}

	// 签名
	sign := t.sign(http.MethodGet, params)
	params["Signature"] = sign

	// 获取响应
	query := url.Values{}
	for k := range params {
		query.Add(k, params[k])
	}
	res, err := t.request(fmt.Sprintf("http://mts.%s.aliyuncs.com/", location), query)
	if err != nil {
		return nil, err
	}

	// 解析响应
	ret := new(MTSQueryMediaListResponse)
	err = json.Unmarshal(res, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// SubmitJobs 提交转码作业
func (t *MTS) SubmitJobs(location string, input VodInputFile, outputBucket string, outputLocation string, output []VodOutputFile, pipelineID string) error {
	params := t.getCommonParams()
	params["Action"] = "SubmitJobs"
	inputJSON, err := json.Marshal(input)
	if err == nil {
		params["Input"] = string(inputJSON)
	}
	params["OutputBucket"] = outputBucket
	params["OutputLocation"] = outputLocation
	outputJSON, err := json.Marshal(output)
	if err == nil {
		params["Outputs"] = string(outputJSON)
	}
	params["PipelineId"] = pipelineID

	// 签名
	sign := t.sign(http.MethodGet, params)
	params["Signature"] = sign

	// 获取响应
	query := url.Values{}
	for k := range params {
		query.Add(k, params[k])
	}
	res, err := t.request(fmt.Sprintf("http://mts.%s.aliyuncs.com/", location), query)
	if err != nil {
		return err
	}
	fmt.Println(string(res))

	// 解析响应
	ret := new(MTSQueryMediaListResponse)
	err = json.Unmarshal(res, ret)
	if err != nil {
		return err
	}

	return nil
}

func (t *MTS) QueryMediaListByURL(location string, fileURLs []string, includePlayList bool, includeSnapshotList bool, includeMediaInfo bool) (*MTSQueryMediaListResponse, error) {
	params := t.getCommonParams()
	params["Action"] = "QueryMediaListByURL"

	for i := range fileURLs {
		item, err := url.Parse(fileURLs[i])
		if err == nil {
			fileURLs[i] = item.String()
		}
	}
	params["FileURLs"] = strings.Join(fileURLs, ",")

	if includePlayList {
		params["IncludePlayList"] = "true"
	}
	if includeSnapshotList {
		params["IncludeSnapshotList"] = "true"
	}
	if includeMediaInfo {
		params["IncludeMediaInfo"] = "true"
	}

	// 签名
	sign := t.sign(http.MethodGet, params)
	params["Signature"] = sign

	// 获取响应
	query := url.Values{}
	for k := range params {
		query.Add(k, params[k])
	}
	res, err := t.request(fmt.Sprintf("http://mts.%s.aliyuncs.com/", location), query)
	if err != nil {
		return nil, err
	}

	// 解析响应
	ret := new(MTSQueryMediaListResponse)
	err = json.Unmarshal(res, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *MTS) request(url string, query url.Values) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	// 超时设置
	client := new(http.Client)
	client.Timeout = time.Second * 10

	// https 支持
	if strings.HasPrefix(url, "https") {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	// 执行请求
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 处理响应
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (t *MTS) sign(method string, params map[string]string) string {
	// 排序
	keys := make([]string, 0)
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接
	pairs := make([]string, 0)
	for i := range keys {
		key := keys[i]
		pairs = append(pairs, t.encode(key)+"="+t.encode(params[key]))
	}
	str := strings.ToUpper(method) + "&" + t.encodePercent("/") + "&" + t.encode(strings.Join(pairs, "&"))

	// 加密
	mac := hmac.New(sha1.New, []byte(t.accessKeySecret+"&"))
	mac.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (t *MTS) encode(str string) string {
	enc := url.QueryEscape(str)
	enc = strings.Replace(enc, "+", "%20", -1)
	enc = strings.Replace(enc, "*", "%2A", -1)
	enc = strings.Replace(enc, "%7E", "~", -1)
	return enc
}

func (t *MTS) encodePercent(str string) string {
	return strings.Replace(str, "/", "%2F", -1)
}

func NewMTS(accessKeyId, accessKeySecret string) (*MTS, error) {
	ins := new(MTS)

	if accessKeyId == "" {
		return nil, fmt.Errorf("Invalid MTS accessKeyId: %s", accessKeyId)
	}
	ins.accessKeyId = accessKeyId

	if accessKeySecret == "" {
		return nil, fmt.Errorf("Invalid MTS accessKeySecret: %s", accessKeySecret)
	}
	ins.accessKeySecret = accessKeySecret

	return ins, nil
}
