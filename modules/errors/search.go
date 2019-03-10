package errors

var (
	// ErrSearchNotExist 搜索数据不存在
	ErrSearchNotExist = New(1300, "搜索数据不存在")
	// ErrSearchListEmpty 搜索结果为空
	ErrSearchListEmpty = New(1301, "搜索结果为空")
)
