package baidu

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"git.code.tencent.com/xinhuameiyu/common/util"
	"github.com/go-baa/log"
)

// OpenAPIUserURL 获取用户信息地址
const OpenAPIUserURL = "https://openapi.baidu.com/rest/2.0/cambrian/sns/userinfo"

// UserInfo 结构 用户信息
type UserInfo struct {
	Openid     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	Headimgurl string `json:"headimgurl"`
}

// userErrorResult 获取用户信息时的报错信息
type userErrorResult struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

// GetUserInfo 获取用户信息
func GetUserInfo(accessToken, openID string) (*UserInfo, *ErrorResult) {
	if accessToken == "" || openID == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}

	data, err := util.HTTPGet(buildAPIRequestURL(OpenAPIUserURL, map[string]string{
		"access_token": accessToken,
		"openid":       openID,
	}), APIRequestTimeout)

	if err != nil {
		log.Errorf("获取百度 用户信息 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "获取百度 用户信息 失败"}
	}

	if err := checkUserResultError(data); err != nil {
		log.Errorf("获取百度 用户信息 接口错误：%s", string(data))
		return nil, err
	}

	ret := new(UserInfo)
	if err := json.Unmarshal(data, ret); err != nil {
		log.Errorf("解析百度 用户信息 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "解析百度 用户信息 响应失败"}
	}

	return ret, nil
}

// TicketInfo js 返回结构的结构体
type TicketInfo struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

// GetTicket 获取js的jsapi_ticket
func GetTicket(accessToken string) (*TicketInfo, *ErrorResult) {
	if accessToken == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}

	data, err := util.HTTPGet(buildAPIRequestURL("https://openapi.baidu.com/rest/2.0/cambrian/jssdk/getticket", map[string]string{
		"access_token": accessToken,
	}), APIRequestTimeout)

	if err != nil {
		log.Errorf("获取百度 JS接口验证码 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "获取百度 JS接口验证码 失败"}
	}

	if err := checkUserResultError(data); err != nil {
		log.Errorf("获取百度 JS接口验证码 接口错误：%s", string(data))
		return nil, err
	}

	ret := new(TicketInfo)
	if err := json.Unmarshal(data, ret); err != nil {
		log.Errorf("解析百度 JS接口验证码 响应失败：%s", err.Error())
		return nil, &ErrorResult{Message: "解析百度 JS接口验证码 响应失败"}
	}

	return ret, nil
}

// TempalteValue 模板值
type TempalteValue struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

// Template 模板结果
type Template struct {
	First    TempalteValue `json:"first"`
	Keyword1 TempalteValue `json:"keyword1"`
	Keyword2 TempalteValue `json:"keyword2"`
	Remark   TempalteValue `json:"remark"`
}

// TemplateSend 发送模板消息
type TemplateSend struct {
	Touser     string   `json:"touser"`
	TemplateID string   `json:"template_id"`
	URL        string   `json:"url"`
	Data       Template `json:"data"`
}

// PushTemplate 发送模板数据
func PushTemplate(accessToken string, params map[string]interface{}) *ErrorResult {
	if accessToken == "" {
		return &ErrorResult{Message: "缺少参数"}
	}
	url := "https://openapi.baidu.com/rest/2.0/cambrian/template/send?access_token=" + accessToken
	data, err := util.HTTPPostJSON(url, params, 3000)
	if err != nil {
		log.Errorf("推送模板消息 失败：%s", err.Error())
		return &ErrorResult{Message: "获取百度 用户信息 失败"}
	}

	if err := checkUserResultError(data); err != nil {
		log.Errorf("推送模板消息 接口错误：%s", string(data))
		return err
	}

	return nil
}

// SendAllMessage 群发消息
func SendAllMessage(accessToken string, openID []string, mediaID int, msgtype string) *ErrorResult {
	if accessToken == "" {
		return &ErrorResult{Message: "缺少参数"}
	}
	params := make(map[string]interface{})
	mpnews := make(map[string]interface{})
	mpnews["media_id"] = mediaID
	params["touser"] = openID
	params["mpnews"] = mpnews
	params["msgtype"] = msgtype
	url := "https://openapi.baidu.com/rest/2.0/cambrian/message/sendall?access_token=" + accessToken
	data, err := util.HTTPPostJSON(url, params, 3000)
	if err != nil {
		log.Errorf("群发消息 失败：%s", err.Error())
		return &ErrorResult{Message: "获取百度 用户信息 失败"}
	}

	if err := checkUserResultError(data); err != nil {
		log.Errorf("群发消息 接口错误：%s", string(data))
		return err
	}

	return nil
}

// Material 百度素材 结构
type Material struct {
	ContentSourceURL string `json:"content_source_url"`
	Title            string `json:"title"`
	ThumbMediaID     int    `json:"thumb_media_id"`
	Author           string `json:"author"`
	Digest           string `json:"digest"`
	Content          string `json:"content"`
}

// MaterialData 素材上传成功后返回的值
type MaterialData struct {
	MediaID int `json:"media_id"`
}

// UploadMaterial 上传素材
func UploadMaterial(accessToken string, articles []Material) (*MaterialData, *ErrorResult) {
	if accessToken == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}
	params := make(map[string]interface{})
	params["articles"] = articles
	url := "https://openapi.baidu.com/rest/2.0/cambrian/material/add_news?access_token=" + accessToken
	data, err := util.HTTPPostJSON(url, params, 3000)
	if err != nil {
		log.Errorf("上传素材 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "获取百度 用户信息 失败"}
	}

	if err := checkUserResultError(data); err != nil {
		log.Errorf("上传素材 接口错误：%s", string(data))
		return nil, err
	}

	row := new(MaterialData)
	if err := json.Unmarshal(data, row); err != nil {
		return nil, &ErrorResult{Message: "数据解析错误"}
	}

	return row, nil
}

// DelMaterialParams 删除参数
type DelMaterialParams struct {
	MediaID int `json:"media_id"`
}

// DelMaterial 删除已上传的素材
func DelMaterial(accessToken string, mediaID int) (*MaterialData, *ErrorResult) {
	if accessToken == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}
	url := "https://openapi.baidu.com/rest/2.0/cambrian/material/del_material?access_token=" + accessToken
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader("media_id="+util.IntToString(mediaID)))
	if err != nil {
		log.Errorf("删除素材 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "删除百度 素材 失败"}
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("删除素材 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "删除百度 素材 失败"}
	}

	if err := checkUserResultError(data); err != nil {
		log.Errorf("删除素材 接口错误：%s", string(data))
		return nil, err
	}

	row := new(MaterialData)
	if err := json.Unmarshal(data, row); err != nil {
		return nil, &ErrorResult{Message: "数据解析错误"}
	}

	return row, nil
}

// MaterialDatas 获取素材信息的结构体
type MaterialDatas struct {
	Title            string `json:"title"`
	ThumbMediaID     int    `json:"thumb_media_id"`
	Author           string `json:"author"`
	Digest           string `json:"digest"`
	Content          string `json:"content"`
	URL              string `json:"url"`
	ContentSourceURL string `json:"content_source_url"`
}

// GetMaterial 获取素材信息
func GetMaterial(accessToken string, mediaID int) (*MaterialDatas, *ErrorResult) {
	if accessToken == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}
	url := "https://openapi.baidu.com/rest/2.0/cambrian/material/get_material?access_token=" + accessToken + "&media_id=" + util.IntToString(mediaID)
	data, err := util.HTTPGet(url, 3000)
	if err != nil {
		log.Errorf("获取百度素材 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "获取百度 素材 失败"}
	}

	if err := checkUserResultError(data); err != nil {
		log.Errorf("获取百度素材 接口错误：%s", string(data))
		return nil, err
	}

	row := new(MaterialDatas)
	if err := json.Unmarshal(data, row); err != nil {
		return nil, &ErrorResult{Message: "数据解析错误"}
	}

	return row, nil
}

// MediaData 图片上传后返回的值
type MediaData struct {
	URL     string `json:"url"`
	MediaID int    `json:"media_id"`
}

// UploadImage 上传图片
func UploadImage(accessToken string, mediaURL string) (*MediaData, *ErrorResult) {
	if accessToken == "" {
		return nil, &ErrorResult{Message: "缺少参数"}
	}
	params := make(map[string]string)
	params["media"] = mediaURL
	url := "https://openapi.baidu.com/rest/2.0/cambrian/media/uploadimg?access_token=" + accessToken
	data, _, err := util.HTTPPostFile(url, params, "baidu-media", 3000)
	if err != nil {
		log.Errorf("图片上传 失败：%s", err.Error())
		return nil, &ErrorResult{Message: "获取百度 用户信息 失败"}
	}

	if err := checkUserResultError(data); err != nil {
		log.Errorf("图片上传 接口错误：%s", string(data))
		return nil, err
	}

	row := new(MediaData)
	if err := json.Unmarshal(data, row); err != nil {
		return nil, &ErrorResult{Message: "数据解析错误"}
	}

	return row, nil
}

// checkUserResultError 检查用户接口是否返回了错误结果
func checkUserResultError(data []byte) *ErrorResult {
	ret := new(userErrorResult)

	if err := json.Unmarshal(data, ret); err != nil {
		return &ErrorResult{Message: ret.ErrorMsg}
	}

	if ret.ErrorCode > 0 {
		return &ErrorResult{Message: ret.ErrorMsg}
	}

	return nil
}
