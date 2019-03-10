package errors

var (
	// ErrContentVideoUploadFail 视频上传失败
	ErrContentVideoUploadFail = New(1200, "视频上传失败")
	// ErrContentVideoTranscodingFail 视频转码失败
	ErrContentVideoTranscodingFail = New(1201, "视频转码失败")
	// ErrContentVideoEmpty 视频没有文件上传
	ErrContentVideoEmpty = New(1203, "视频没有文件上传")
	// ErrContentVideoCheckFail 视频送审失败
	ErrContentVideoCheckFail = New(1204, "视频送审失败")
	// ErrContentVideoRevokeFail 视频撤回失败
	ErrContentVideoRevokeFail = New(1205, "视频撤回失败")
	// ErrContentVideoDeleteFail 视频删除失败
	ErrContentVideoDeleteFail = New(1206, "视频删除失败")
	// ErrContentVideoRejectFail 获取驳回数据失败
	ErrContentVideoRejectFail = New(1207, "获取驳回数据失败")
	// ErrContentVideoListFail 获取视频列表失败
	ErrContentVideoListFail = New(1208, "获取视频列表失败")
	// ErrContentVideoStatFail 获取视频统计数据失败
	ErrContentVideoStatFail = New(1209, "获取视频统计数据失败")
	// ErrContentVideoEmptyFail 没有视频数据
	ErrContentVideoEmptyFail = New(1210, "没有视频数据")
	// ErrContentVideoExpertFail 视频不属于此专家
	ErrContentVideoExpertFail = New(1211, "视频不属于此专家")
	// ErrContentVideoPowerFail 无权操作此视频
	ErrContentVideoPowerFail = New(1212, "无权操作此视频")
	// ErrContentPraiseFail 点赞失败
	ErrContentPraiseFail = New(1301, "点赞失败")
	// ErrContentTreadFail 点踩失败
	ErrContentTreadFail = New(1303, "点踩失败")
	// ErrContentPraiseTreadFail 已点过
	ErrContentPraiseTreadFail = New(1304, "已点过")
)
