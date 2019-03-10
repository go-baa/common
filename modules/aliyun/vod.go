package aliyun

// VodMessage 视频点播通知消息
type VodMessage struct {
	RunId                  string
	Name                   string
	Type                   string
	State                  string
	MediaWorkflowExecution *VodMediaWorkflowExecution
}

// VodMediaWorkflowExecution 视频点播工作流消息
type VodMediaWorkflowExecution struct {
	Name         string
	RunId        string
	Input        *VodInput
	State        string
	MediaId      string
	ActivityList []*VodActivity
	CreationTime string
}

// VodInput 视频点播输入
type VodInput struct {
	InputFile *VodInputFile
}

// VodInputFile 视频点播输入文件
type VodInputFile struct {
	Bucket   string
	Location string
	Object   string
}

// VodOutputFile 视频点播输出文件
type VodOutputFile struct {
	OutputObject string
	TemplateID   string `json:"TemplateId"`
}

// VodActivity 视频点播活动对象
type VodActivity struct {
	RunId     string
	Name      string
	Type      string
	State     string
	StartTime string
	EndTime   string
}
