package baidu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"git.code.tencent.com/xinhuameiyu/common/util"
	"git.code.tencent.com/xinhuameiyu/common/util/uuid"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
)

// Text2AudioParams 语音合成参数
type Text2AudioParams struct {
	Token string // (必填) access_token
	Text  string // (必填) 合成文本：最大1024字节
	Lan   string // (必填) 语种选择：中文=zh、粤语=ct、英文=en，不区分大小写，默认中文
}

// Text2AudioResponseError 语音合成错误
type Text2AudioResponseError struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	SN     string `json:"sn"`
	Idx    int    `json:"idx"`
}

// Audio2TextParams 语音识别参数
// 注意："speech"和"len"参数绑定，"url"和"callback"参数绑定，这两组参数二选一填写
type Audio2TextParams struct {
	Format string `json:"format"` // (必填) 语音格式:pcm、wav、opus、speex、amr、x-flac
	Rate   int    `json:"rate"`   // (必填) 采样率：8000 / 16000
	Token  string `json:"token"`  // (必填) access_token
	Lan    string `json:"lan"`    // (选填) 语种选择：中文=zh、粤语=ct、英文=en，不区分大小写，默认中文
}

// Audio2TextResponse 语音识别响应
type Audio2TextResponse struct {
	ErrNo  int      `json:"err_no"`
	ErrMsg string   `json:"err_msg"`
	SN     string   `json:"sn"`
	Result []string `json:"result"`
}

// AudioResponseError 语音服务响应错误
var AudioResponseError = map[int]string{
	500:  "不支持输入",
	501:  "输入参数不正确",
	502:  "token验证失败",
	503:  "合成后端错误",
	3300: "输入参数不正确",
	3301: "识别错误",
	3302: "验证失败",
	3303: "语音服务器后端问题",
	3304: "请求 GPS 过大，超过限额",
	3305: "产品线当前日请求数超过限额",
}

// Text2Audio 文本转语音，返回下载地址
func (t *Baidu) Text2Audio(params *Text2AudioParams, filePath string) error {
	if params.Lan == "" {
		params.Lan = "zh"
	}
	gateway := "http://tsn.baidu.com/text2audio"
	queryParams := map[string]interface{}{
		"tex":  params.Text,
		"lan":  params.Lan,
		"tok":  params.Token,
		"ctp":  "1",
		"cuid": uuid.NewV1().String(),
		"spd":  "5",
	}
	query := util.HTTPBuildQuery(queryParams, 0)
	downLoadURL := gateway + "?" + query

	// 判断目录是否存在
	if err := util.MkdirAll(filepath.Dir(filePath)); err != nil {
		return fmt.Errorf("创建目录失败：%v", err)
	}

	resp, err := http.Get(downLoadURL)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取数据错误: %v", err)
	}
	resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType == "application/json" {
		response := new(Text2AudioResponseError)
		err = json.Unmarshal(body, response)
		if err != nil {
			return fmt.Errorf("JSON解码错误:%v", err)
		}

		if msg, ok := AudioResponseError[response.ErrNo]; ok {
			return fmt.Errorf(msg)
		}

		return fmt.Errorf("语音合成错误:code:%d, msg:%s", response.ErrNo, response.ErrMsg)
	}

	if contentType != "audio/mp3" {
		return fmt.Errorf("语音合成错误：content-type: %s", contentType)
	}

	_, err = util.WriteFile(filePath, body)
	if err != nil {
		return fmt.Errorf("保存文件错误: %v", err)
	}

	return nil
}

// DefaultAudio2TextParams 默认语音识别配置
func DefaultAudio2TextParams() *Audio2TextParams {
	return &Audio2TextParams{
		Format: "amr",
		Rate:   16000,
		Lan:    "zh",
	}
}

// Audio2Text 语音转文本
func (t *Baidu) Audio2Text(params *Audio2TextParams, filePath string) ([]string, error) {
	gateway := "http://vop.baidu.com/server_api"
	if params.Lan == "" {
		params.Lan = "zh"
	}

	f, err := util.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件数据错误: %v", err)
	}

	queryParams := map[string]interface{}{
		"format":  params.Format,
		"channel": "1",
		"cuid":    uuid.NewV1().String(),
		"token":   params.Token,
		"lan":     params.Lan,
	}
	query := util.HTTPBuildQuery(queryParams, 0)
	header := map[string]string{
		"Content-Type": fmt.Sprintf("audio/%s;rate=%d", params.Format, params.Rate),
	}

	res, err := t.request(gateway, query, f, header)
	if err != nil {
		return nil, fmt.Errorf("请求失败:%v", err)
	}

	if setting.Debug {
		log.Debugf("baidu res: %v\n", string(res))
	}

	response := new(Audio2TextResponse)
	err = json.Unmarshal(res, response)
	if err != nil {
		log.Errorf("baidu res parse error: %s\n", res)
		return nil, fmt.Errorf("解析响应结果错误:%v", err)
	}

	if response.ErrNo > 0 {
		log.Errorf("baidu res server error: %#v\n", response)
		if _, ok := AudioResponseError[response.ErrNo]; ok {
			return nil, fmt.Errorf("识别错误: %#v", response)
		}
		return nil, fmt.Errorf("未知错误: %#v", response)
	}

	return response.Result, nil
}
