package om

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-baa/common/util"
	"github.com/go-baa/log"
)

// OpenAPITokenURL OPEN API token获取地址
const OpenAPITokenURL = "https://auth.om.qq.com/omoauth2/accesstoken"

// APIRequestTimeout 请求超时时间
var APIRequestTimeout = 120

// TencentOmAccessToken token返回参数
type TencentOmAccessToken struct {
	AccessToken string `json:"access_token"` // 获取到的网页授权接口调用凭证
	ExpiresIn   int    `json:"expires_in"`   // 凭证有效时间，单位：秒
	Openid      string `json:"openid"`       // 用户唯一标识
}

// accessTokenResult access_token 返回数据
type accessTokenResult struct {
	Code string                `json:"code"`
	Msg  string                `json:"msg"`
	Data *TencentOmAccessToken `json:"data"`
}

// ErrorResult 错误信息
type ErrorResult struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

// GetTencentOmAccessToken 获取token
func GetTencentOmAccessToken(appID, secret string) (*TencentOmAccessToken, *ErrorResult) {
	if appID == "" || secret == "" {
		return nil, &ErrorResult{Msg: "缺少参数"}
	}
	params := make(map[string]interface{})
	// params["grant_type"] = "clientcredentials"
	// params["client_id"] = appID
	// params["client_secret"] = secret
	data, err := util.HTTPPostJSON(OpenAPITokenURL+"?grant_type=clientcredentials&client_id="+appID+"&client_secret="+secret, params, APIRequestTimeout)

	if err != nil {
		log.Errorf("获取腾讯内容平台 access_token 失败：%s", err.Error())
		return nil, &ErrorResult{Msg: "获取腾讯内容平台 access_token 失败"}
	}

	ret := new(accessTokenResult)
	if err := json.Unmarshal(data, &ret); err != nil {
		log.Errorf("获取腾讯内容平台 access_token 接口错误：%s", err.Error())
		return nil, &ErrorResult{Msg: "获取腾讯内容平台 access_token 接口错误"}
	}

	if ret.Code != "0" {
		log.Errorf("获取腾讯内容平台 access_token 接口错误：%s", ret.Msg)
		return nil, &ErrorResult{Msg: ret.Msg}
	}

	return ret.Data, nil
}

// TencentVideoParam 向腾讯内容平台推送的参数
type TencentVideoParam struct {
	Title    string `json:"title"`
	Tags     string `json:"tags"`
	Cat      int    `json:"cat"`
	Desc     string `json:"desc"`
	VideoURL string `json:"videourl"`
	Apply    int    `json:"apply"`
}

// TencentVideo 返回的参数
type TencentVideo struct {
	TransactionID int `json:"transaction_id"`
}

// tencentVideoResult 视频信息返回数据
type tencentVideoResult struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data *TencentVideo `json:"data"`
}

// PushTencentVideo 向腾讯内容平台推送视频数据
func PushTencentVideo(accessToken string, param TencentVideoParam) (*TencentVideo, *ErrorResult) {
	url := "https://api.om.qq.com/article/clientpuburlvid?access_token=" + accessToken

	params := make(map[string]interface{})
	// params["title"] = param.Title
	// params["tags"] = param.Tags
	// params["cat"] = param.Cat
	// params["desc"] = param.Desc
	// params["videourl"] = param.VideoURL
	// params["apply"] = param.Apply
	param.Title = strings.Replace(param.Title, "【", "", -1)
	param.Title = strings.Replace(param.Title, "】", "：", -1)
	data, err := util.HTTPPostJSON(url+"&title="+param.Title+"&tags="+param.Tags+"&cat="+util.IntToString(param.Cat)+"&desc="+param.Desc+"&videourl="+param.VideoURL+"&apply="+util.IntToString(param.Apply), params, APIRequestTimeout)

	if err != nil {
		log.Errorf("腾讯内容平台 历史视频上传失败：%s", err.Error())
		return nil, &ErrorResult{Msg: "腾讯内容平台 历史视频上传失败"}
	}

	ret := new(tencentVideoResult)
	if err := json.Unmarshal(data, &ret); err != nil {
		log.Errorf("腾讯内容平台 历史视频上传失败：%s", err.Error())
		return nil, &ErrorResult{Msg: "腾讯内容平台 历史视频上传失败"}
	}

	if ret.Code != 0 {
		log.Errorf("腾讯内容平台 历史视频上传失败：%s", ret.Msg)
		return nil, &ErrorResult{Msg: ret.Msg}
	}

	return ret.Data, nil
}

// TransactionInfo 根据事物编号获取的详细信息
type TransactionInfo struct {
	TransactionID     string       `json:"transaction_id"`
	TransactionStatus string       `json:"transaction_status"`
	TransactionType   string       `json:"transaction_type"`
	TransactionCtime  string       `json:"transaction_ctime"`
	ArticleInfo       *ArticleInfo `json:"article_info"`
}

// ArticleInfo 文章信息
type ArticleInfo struct {
	ArticleAbstract  string            `json:"article_abstract"`
	ArticleImgurl    string            `json:"article_imgurl"`
	ArticlePubFlag   string            `json:"article_pub_flag"`
	ArticlePubTime   string            `json:"article_pub_time"`
	ArticleTitle     string            `json:"article_title"`
	ArticleType      string            `json:"article_type"`
	ArticleURL       string            `json:"article_url"`
	ArticleVideoInfo *ArticleVideoInfo `json:"article_video_info"`
}

// ArticleVideoInfo 文章中的视频数据
type ArticleVideoInfo struct {
	Desc      string `json:"desc"`
	Title     string `json:"title"`
	VideoType string `json:"type"`
	VID       string `json:"vid"`
}

// transactionInfoResult 视频信息返回数据
type transactionInfoResult struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data *TransactionInfo `json:"data"`
}

// GetVideoInfo 根据事物编号获取视频信息
func GetVideoInfo(accessToken, transactionID string) (*TransactionInfo, *ErrorResult) {
	if accessToken == "" || transactionID == "" {
		return nil, &ErrorResult{Msg: "腾讯内容平台 获取事物编号的数据错误：参数不正确"}
	}
	url := "https://api.om.qq.com/transaction/infoclient?access_token=" + accessToken + "&transaction_id=" + transactionID

	data, err := util.HTTPGet(url, APIRequestTimeout)
	if err != nil {
		log.Errorf("腾讯内容平台 获取事物编号的数据错误：%s", err.Error())
		return nil, &ErrorResult{Msg: "腾讯内容平台 获取事物编号的数据错误"}
	}

	ret := new(transactionInfoResult)
	if err := json.Unmarshal(data, &ret); err != nil {
		log.Errorf("腾讯内容平台 获取事物编号的数据错误：%s", err.Error())
		return nil, &ErrorResult{Msg: "腾讯内容平台 获取事物编号的数据错误"}
	}

	if ret.Code != 0 {
		log.Errorf("腾讯内容平台 获取事物编号的数据错误：%s", ret.Msg)
		return nil, &ErrorResult{Msg: ret.Msg}
	}

	return ret.Data, nil
}

// GetVideotTansactionID 根据视频大小等获取事物编号
func GetVideotTansactionID(accessToken, size, md5, sha string) (*tencentVideoResult, *ErrorResult) {
	if accessToken == "" || size == "" || md5 == "" || sha == "" {
		return nil, &ErrorResult{Msg: "腾讯内容平台 获取事物编号错误：参数不正确"}
	}
	url := "https://api.om.qq.com/video/clientuploadready?access_token=" + accessToken + "&size=" + size + "&md5=" + md5 + "&sha=" + sha
	params := make(map[string]interface{})
	data, err := util.HTTPPostJSON(url, params, APIRequestTimeout)
	fmt.Println(err)
	fmt.Println(string(data))
	if err != nil {
		log.Errorf("腾讯内容平台 获取事物编号错误：%s", err.Error())
		return nil, &ErrorResult{Msg: "腾讯内容平台 获取事物编号错误"}
	}

	ret := new(tencentVideoResult)
	if err := json.Unmarshal(data, &ret); err != nil {
		log.Errorf("腾讯内容平台 获取事物编号错误：%s", err.Error())
		return nil, &ErrorResult{Msg: "腾讯内容平台 获取事物编号错误"}
	}

	if ret.Code != 0 {
		log.Errorf("腾讯内容平台 获取事物编号错误：%s", ret.Msg)
		return nil, &ErrorResult{Msg: ret.Msg}
	}

	return ret, nil
}

// UploadVideoResult 视频上传后返回数据
type UploadVideoResult struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data *UploadVideoData `json:"data"`
}

// UploadVideoData 返回的数据
type UploadVideoData struct {
	EndOffset     int    `json:"end_offset"`
	StartOffset   int    `json:"start_offset"`
	TransactionID string `json:"transaction_id"`
}

// UploadVideo 上传视频
func UploadVideo(accessToken, transactionID, mediatrunk string, startOffset int) (*UploadVideoResult, *ErrorResult) {
	if accessToken == "" || transactionID == "" || mediatrunk == "" || startOffset <= -1 {
		return nil, &ErrorResult{Msg: "腾讯内容平台 上传视频错误：参数不正确"}
	}
	url := "http://api.om.qq.com/video/clientuploadtrunk?access_token=" + accessToken + "&transaction_id=" + transactionID + "&start_offset=" + util.IntToString(startOffset)
	params := make(map[string]string)
	if strings.Contains(mediatrunk, "@") {
		params["mediatrunk"] = mediatrunk
	} else {
		params["mediatrunk"] = "@" + mediatrunk
	}
	data, _, err := util.HTTPPostFile(url, params, "", 7200)
	fmt.Println(321)
	fmt.Println(err)
	fmt.Println(data)
	fmt.Println(123)
	return nil, &ErrorResult{Msg: "腾讯内容平台 上传视频错误"}
	if err != nil {
		log.Errorf("腾讯内容平台 上传视频错误：%s", err.Error())
		return nil, &ErrorResult{Msg: "腾讯内容平台 上传视频错误"}
	}

	ret := new(UploadVideoResult)
	if err := json.Unmarshal(data, &ret); err != nil {
		log.Errorf("腾讯内容平台 上传视频错误：%s", err.Error())
		return nil, &ErrorResult{Msg: "腾讯内容平台 上传视频错误"}
	}

	if ret.Code != 0 {
		log.Errorf("腾讯内容平台 上传视频错误：%s", ret.Msg)
		return nil, &ErrorResult{Msg: ret.Msg}
	}

	return ret, nil
}
