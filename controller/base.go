package controller

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.code.tencent.com/xinhuameiyu/common/modules/assets"
	"git.code.tencent.com/xinhuameiyu/common/modules/errors"
	"git.code.tencent.com/xinhuameiyu/common/util"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
	"gopkg.in/baa.v1"
)

// NormalReturn 标准返回格式
type NormalReturn struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// NewReturn 实例化一个新的返回结构
func NewReturn() *NormalReturn {
	re := new(NormalReturn)
	re.Data = make(map[string]interface{})
	return re
}

// Errorf 对fmt.Errorf()的一个包装
func Errorf(format string, a ...interface{}) error {
	if len(a) > 0 {
		return fmt.Errorf(format, a...)
	}
	return fmt.Errorf(format)
}

// ParseError 解析错误代码和错误消息
func ParseError(err error) (int, string) {
	switch e := err.(type) {
	case *errors.APIError:
		return e.Code, e.Message
	default:
		return 1, err.Error()
	}
}

// Request 获取请求的参数字典，含 url和form中的数据，不含上传的文件
func Request(c *baa.Context) map[string]interface{} {
	params := make(map[string]interface{})
	hasFile := false
	c.ParseForm(0)
	if c.Req.MultipartForm != nil && len(c.Req.MultipartForm.File) > 0 {
		hasFile = true
	}
	for k, v := range c.Req.Form {
		if hasFile && c.Req.MultipartForm.File[k] != nil {
			continue
		}
		if len(v) > 1 {
			params[k] = v
		} else {
			params[k] = v[0]
		}
	}
	return params
}

// RequestQuery 获取请求的参数字典，仅含 url中的数据
func RequestQuery(c *baa.Context) map[string]interface{} {
	params := make(map[string]interface{})
	var newValues url.Values
	if c.Req.URL != nil {
		newValues, _ = url.ParseQuery(c.Req.URL.RawQuery)
	}
	for k, v := range newValues {
		if len(v) > 1 {
			params[k] = v
		} else {
			params[k] = v[0]
		}
	}
	return params
}

// RequestForm 获取请求的参数字典，不含url中的数据
func RequestForm(c *baa.Context) map[string]interface{} {
	params := make(map[string]interface{})
	c.Req.ParseForm()
	for k, v := range c.Req.PostForm {
		if len(v) > 1 {
			params[k] = v
		} else {
			params[k] = v[0]
		}
	}
	return params
}

// RequestTree 格式化request对象，如果是data[][option]这样的格式，进行多级map转换
func RequestTree(c *baa.Context, name ...string) interface{} {
	params := make(map[string]interface{})
	c.Req.ParseForm()
	for k, v := range c.Req.Form {
		if len(k) > 2 && k[len(k)-2] == '[' && k[len(k)-1] == ']' {
			k = k[:len(k)-2]
		}
		if _, ok := params[k]; ok {
			params[k] = append(params[k].([]string), v...)
		} else {
			params[k] = v
		}
	}
	newParams := make(map[string]interface{})
	for k, v := range params {
		if _, ok := v.([]string); ok {
			if len(v.([]string)) == 1 {
				v = v.([]string)[0]
			}
		}
		formatTree(newParams, k, v)
	}
	if len(name) == 0 {
		return newParams
	}
	return newParams[name[0]]
}

func formatTree(m map[string]interface{}, k string, v interface{}) {
	var keys []string
	var key []byte
	for i := 0; i < len(k); i++ {
		if k[i] == '[' {
			if len(key) > 0 {
				keys = append(keys, string(key))
				key = key[:0]
			}
			for i = i + 1; i < len(k); i++ {
				if k[i] == '\'' || k[i] == '"' {
					continue
				}
				if k[i] == ']' {
					break
				}
				key = append(key, k[i])
			}
			if len(key) > 0 {
				keys = append(keys, string(key))
				key = key[:0]
			} else {
				keys = append(keys, "")
			}
			continue
		}
		key = append(key, k[i])
	}

	if len(key) > 0 {
		keys = append(keys, string(key))
		key = key[:0]
	}

	// 只有一级
	if len(keys) == 1 {
		m[keys[0]] = v
		return
	}

	// 二级含以上
	nm := m
	for i := 0; i < len(keys); i++ {
		if _, ok := nm[keys[i]]; !ok {
			nm[keys[i]] = make(map[string]interface{})
		}
		if i+1 == len(keys) {
			nm[keys[i]] = v
			break
		}
		nm = nm[keys[i]].(map[string]interface{})
	}
}

var (
	globalUploadBasePath  string // 全局上传根路径
	globalUploadBaseURI   string // 全局上传根URI
	globalUploadExtension string // 全局文件后缀限制
	globalUploadMaxsize   int64  // 全局文件大小限制
)

// UploadFile 从Http流中上传一个文件, 返回上传后的文件地址
// 默认上传获取到的第一个文件，如果指定 fieldName 仅上传指定的文件
// 允许指定一个附件路径，默认会上传到upload目录，如果有addonPath则附加
func UploadFile(uploadType, fieldName string, c *baa.Context, addonPath string) (string, string, error) {
	maxSize := setting.Config.MustInt64("upload."+uploadType+".maxsize", globalUploadMaxsize)
	err := c.Req.ParseMultipartForm(maxSize)
	if err != nil {
		return "", "", Errorf("超过上传限制，最大允许 %d m, %s", maxSize, err)
	}

	// 如果没有指定 文件字段，取第一个获取到的文件
	if fieldName == "" {
		for k := range c.Req.MultipartForm.File {
			fieldName = k
			break
		}
	}
	files := c.Req.MultipartForm.File[fieldName]
	if len(files) == 0 {
		return "", "", Errorf("没有文件被上传")
	}

	ext := strings.ToLower(filepath.Ext(files[0].Filename))
	allowExt := strings.Split(setting.Config.MustString("upload."+uploadType+".extension", globalUploadExtension), ";")
	if ext == "" || util.InSlice(allowExt, ext, "string") == false {
		return "", "", Errorf("只允许上传指定的格式: %s", strings.Join(allowExt, ";"))
	}

	file, err := files[0].Open()
	if err != nil {
		return "", "", Errorf("文件读取失败: %s", err)
	}
	defer file.Close()

	uploadPath := strings.Trim(setting.Config.MustString("upload."+uploadType+".path", ""), "/")
	if uploadPath == "" {
		uploadPath = uploadType
	}
	uploadPath = globalUploadBasePath + "/" + uploadPath
	uploadPath, err = filepath.Abs(uploadPath)
	if err != nil {
		return "", "", Errorf("上传目录转化失败: %s", err)
	}
	if addonPath != "" {
		addonPath = "/" + strings.Trim(addonPath, "/")
	}
	err = util.MkdirAll(uploadPath + addonPath)
	if err != nil {
		return "", "", Errorf("上传目录创建失败: %s", err)
	}
	dstPath := uploadPath + addonPath + "/" + util.RandFileName() + ext
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", "", Errorf("文件创建失败: %s", err)
	}
	defer dst.Close()
	size, err := io.Copy(dst, file)
	if err != nil {
		return "", "", Errorf("文件写入失败: %s", err)
	}

	if setting.Debug {
		log.Infof("upload a file: %s fileSzie: %d saved to %s\n", files[0].Filename, size, dstPath)
	}

	uploadURI := strings.Trim(setting.Config.MustString("upload."+uploadType+".uri", ""), "/")
	if uploadURI == "" {
		uploadURI = uploadType
	}
	uploadURI = globalUploadBaseURI + "/" + uploadURI
	if len(uploadURI) == 0 {
		return dstPath, "", nil
	}
	uploadRelativeURI := dstPath[len(uploadPath):]
	return dstPath, uploadURI + uploadRelativeURI, nil
}

// UploadBase64 ...
func UploadBase64(uploadType string, alias string, data string, addonPath string) (string, string, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return "", "", Errorf("文件内容格式不正确")
	}
	content, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", Errorf("文件内容解析失败")
	}

	maxSize := setting.Config.MustInt64("upload."+uploadType+".maxsize", globalUploadMaxsize)
	if len(content) > int(maxSize) {
		return "", "", Errorf("超过上传限制，最大允许 %d m, %s", maxSize, err)
	}

	ext := strings.ToLower(filepath.Ext(alias))
	allowExt := strings.Split(setting.Config.MustString("upload."+uploadType+".extension", globalUploadExtension), ";")
	if ext == "" || util.InSlice(allowExt, ext, "string") == false {
		return "", "", Errorf("只允许上传指定的格式: %s", strings.Join(allowExt, ";"))
	}

	uploadPath := strings.Trim(setting.Config.MustString("upload."+uploadType+".path", ""), "/")
	if uploadPath == "" {
		uploadPath = uploadType
	}
	uploadPath = globalUploadBasePath + "/" + uploadPath
	uploadPath, err = filepath.Abs(uploadPath)
	if err != nil {
		return "", "", Errorf("上传目录转化失败: %s", err)
	}
	if addonPath != "" {
		addonPath = "/" + strings.Trim(addonPath, "/")
	}
	err = util.MkdirAll(uploadPath + addonPath)
	if err != nil {
		return "", "", Errorf("上传目录创建失败: %s", err)
	}
	dstPath := uploadPath + addonPath + "/" + util.RandFileName() + ext
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", "", Errorf("文件创建失败: %s", err)
	}
	defer dst.Close()
	size, err := dst.Write(content)
	if err != nil {
		return "", "", Errorf("文件写入失败: %s", err)
	}

	if setting.Debug {
		log.Infof("upload a file: %s fileSzie: %d saved to %s\n", alias, size, dstPath)
	}

	uploadURI := strings.Trim(setting.Config.MustString("upload."+uploadType+".uri", ""), "/")
	if uploadURI == "" {
		uploadURI = uploadType
	}
	uploadURI = globalUploadBaseURI + "/" + uploadURI
	if len(uploadURI) == 0 {
		return dstPath, "", nil
	}
	uploadURI = strings.TrimRight(uploadURI, "/")
	uploadRelativeURI := dstPath[len(uploadPath):]
	if len(uploadRelativeURI) == 0 {
		return dstPath, "", nil
	}

	return dstPath, assets.AbsoluteUploadURL(uploadURI + uploadRelativeURI), nil
}

// GetStartEndByType 获取一段时间的开始和结束
func GetStartEndByType(name string) (time.Time, time.Time) {
	var start, end time.Time

	now := time.Now()

	switch name {
	case "today":
		start = now
		end = now
		break
	case "yesterday":
		start = now.AddDate(0, 0, -1)
		end = start
		break
	case "week":
		start = now
		for start.Weekday() != time.Monday {
			start = start.AddDate(0, 0, -1)
		}
		end = now
		break
	case "month":
		start = now.AddDate(0, 0, -1*now.Day()+1)
		end = now
		break
	}

	return start, end
}

// ConvertPageToOffset 把页码转化为偏移量
func ConvertPageToOffset(page, pagesize, pagesizeDefault, pagesizeMax int) (offset int, limit int) {
	if page == 0 {
		page = 1
	}
	if pagesizeDefault == 0 {
		pagesizeDefault = 10
	}
	if pagesizeMax == 0 {
		pagesizeMax = 10
	}

	if pagesizeMax > 0 && pagesize > pagesizeMax {
		pagesize = pagesizeMax
	}
	if pagesize == 0 {
		pagesize = pagesizeDefault
	}

	offset = (page - 1) * pagesize
	limit = pagesize

	return
}

// ConvertOffsetToPage 把偏移量转化为页码
func ConvertOffsetToPage(offset, limit int) (page int, pagesize int) {
	if limit == 0 {
		return 0, 0
	}

	if offset < limit {
		return 1, limit
	}

	page = int(math.Ceil(float64(offset) / float64(limit)))
	pagesize = limit

	if page*pagesize == offset {
		page++
	}

	return
}

// Watermark 添加水印
func Watermark(sourceImg, newImg, watermarkImg string) (string, error) {
	if watermarkImg == "" {
		staticPath := setting.Config.MustString("static.basePath", "")
		staticPath, err := filepath.Abs(staticPath)
		if err != nil {
			return "", err
		}
		watermarkImg = staticPath + "/images/watermark/watermark.png"
	}
	if newImg == "" {
		newImg = sourceImg
	}
	// 原始图片
	imgb, _ := os.Open(sourceImg)
	img, _ := jpeg.Decode(imgb)
	defer imgb.Close()

	fmt.Println(sourceImg)
	fmt.Println(watermarkImg)
	wmb, _ := os.Open(watermarkImg)
	watermark, _ := png.Decode(wmb)
	defer wmb.Close()
	//把水印写到右下角，并向0坐标各偏移20个像素
	offset := image.Pt(img.Bounds().Dx()-watermark.Bounds().Dx()-20, img.Bounds().Dy()-watermark.Bounds().Dy()-10)
	b := img.Bounds()
	m := image.NewNRGBA(b)

	draw.Draw(m, b, img, image.ZP, draw.Src)
	draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)

	// 生成新图片，并设置图片质量..
	imgw, _ := os.Create(newImg)
	err := jpeg.Encode(imgw, m, &jpeg.Options{100})
	if err != nil {
		return "", err
	}

	defer imgw.Close()
	return newImg, nil
}

func init() {
	// 处理全局的上传配置
	globalUploadBasePath = setting.Config.MustString("upload.basePath", "")
	if len(globalUploadBasePath) > 1 {
		globalUploadBasePath = strings.TrimRight(globalUploadBasePath, "/")
	}
	globalUploadBaseURI = setting.Config.MustString("upload.baseUri", "")
	if len(globalUploadBaseURI) > 1 {
		globalUploadBaseURI = strings.TrimRight(globalUploadBaseURI, "/")
	}
	globalUploadExtension = setting.Config.MustString("upload.extension", "")
	globalUploadMaxsize = setting.Config.MustInt64("upload.maxsize", 1048576) // 1m
}
