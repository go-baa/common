package im

const (
	// MsgTypeText 文本消息类型
	MsgTypeText = "TIMTextElem"
	// MsgTypeLocation 地理位置消息类型
	MsgTypeLocation = "TIMLocationElem"
	// MsgTypeFace 表情消息类型
	MsgTypeFace = "TIMFaceElem"
	// MsgTypeCustom 自定义消息类型
	MsgTypeCustom = "TIMCustomElem"
	// MsgTypeSound 语音消息类型
	MsgTypeSound = "TIMSoundElem"
	// MsgTypeImage 图像消息类型
	MsgTypeImage = "TIMImageElem"
	// MsgTypeFile 文件消息类型
	MsgTypeFile = "TIMFileElem"
)

const (
	// ImageFormatOther 其它
	ImageFormatOther = iota
	// ImageFormatBMP BMP类型
	ImageFormatBMP
	// ImageFormatJPG JPG类型
	ImageFormatJPG
	// ImageFormatGIF GIF类型
	ImageFormatGIF
)

// TextMsgContent 文本消息内容
type TextMsgContent struct {
	Text string // 消息内容
}

// NewTextMsgBody 实例化文本消息体
func (t *IM) NewTextMsgBody(text string) *MsgBodyItem {
	body := new(MsgBodyItem)
	body.MsgType = MsgTypeText
	body.MsgContent = TextMsgContent{Text: text}

	return body
}

// LocationMsgContent 地理位置消息内容
type LocationMsgContent struct {
	Desc      string  // 地理位置描述信息
	Latitude  float64 // 纬度
	Longitude float64 // 经度
}

// NewLocationMsgBody 实例化地理位置消息体
func (t *IM) NewLocationMsgBody(desc string, latitude, longitude float64) *MsgBodyItem {
	body := new(MsgBodyItem)
	body.MsgType = MsgTypeLocation
	body.MsgContent = LocationMsgContent{Desc: desc, Latitude: latitude, Longitude: longitude}

	return body
}

// FaceMsgContent 表情消息内容
type FaceMsgContent struct {
	Index int    // 表情索引，用户自定义
	Data  string // 额外数据
}

// NewFaceMsgBody 实例化表情消息体
func (t *IM) NewFaceMsgBody(index int, data string) *MsgBodyItem {
	body := new(MsgBodyItem)
	body.MsgType = MsgTypeFace
	body.MsgContent = FaceMsgContent{Index: index, Data: data}

	return body
}

// CustomMsgContent 自定义消息内容
type CustomMsgContent struct {
	Data  string // 自定义消息数据
	Desc  string // 自定义消息描述信息
	Ext   string // 扩展字段
	Sound string // 自定义APNS推送铃音
}

// NewCustomMsgBody 实例化自定义消息内容
func (t *IM) NewCustomMsgBody(data, desc, ext, sound string) *MsgBodyItem {
	body := new(MsgBodyItem)
	body.MsgType = MsgTypeCustom
	body.MsgContent = CustomMsgContent{Data: data, Desc: desc, Ext: ext, Sound: sound}

	return body
}

// SoundMsgContent 语音消息内容，服务端不支持发送本类型消息
type SoundMsgContent struct {
	UUID   string // 语音序列号。后台用于索引语音的键值
	Size   int64  // 语音数据大小
	Second int64  // 语音时长，单位:秒
}

// ImageInfo 图像信息
type ImageInfo struct {
	Type   int    // 图片类型： 1-原图，2-大图 3-缩略图
	Size   int64  // 图片数据大小
	Width  int    // 图片宽度
	Height int    // 图片高度
	URL    string // 图片下载地址
}

// ImageMsgContent 图像消息内容，服务端不支持发送本类型消息
type ImageMsgContent struct {
	UUID           string       // 图片序列号。后台用于索引图片的键值
	ImageFormat    int          // 图片格式。BMP=1,JPG=2,GIF=3,其他=0
	ImageInfoArray []*ImageInfo // 原图、缩略图或者大图下载信息
}

// FileMsgContent 文件消息内容，服务端不支持发送本类型消息
type FileMsgContent struct {
	UUID     string // 文件序列号。后台用于索引语音的键值
	FileSize int64  // 文件数据大小
	FileName string // 文件名称
}

// MsgBodyItem 消息体
type MsgBodyItem struct {
	MsgType    string      // 消息对象类型
	MsgContent interface{} // 消息内容
}

// OfflinePushInfo 离线推送消息
type OfflinePushInfo struct {
	PushFlag int    // 0表示推送，1表示不离线推送
	Desc     string // 离线推送内容
	Ext      string // 离线推送透传内容
	Sound    string // 离线推送声音文件路径
}
