package errors

var (
	// ErrTopicSelectedFail 选题认领失败
	ErrTopicSelectedFail = New(1501, "选题认领失败")
	// ErrTopicListFail 选题列表为空
	ErrTopicListFail = New(1502, "选题列表为空")
)
