package errors

var (
	// ErrChannelNotExist 频道不存在
	ErrChannelNotExist = New(1300, "频道不存在")
	// ErrChannelListEmpty 频道列表为空
	ErrChannelListEmpty = New(1301, "频道列表为空")
)
