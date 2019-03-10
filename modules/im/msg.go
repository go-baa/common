package im

// 消息模型定义，除文本外都是自定义消息类型
const (
	// MsgModelText 文本
	MsgModelText = 0
	// MsgModelCommon 通用
	MsgModelCommon = 1
	// MsgModelArticle 文章
	MsgModelArticle = 2
	// MsgModelImage 图片
	MsgModelImage = 3
	// MsgModelAnnounce 公告
	MsgModelAnnounce = 4
	// MsgModelApproval 审批
	MsgModelApproval = 5
	// MsgModelCourseTask 课程任务
	MsgModelCourseTask = 6
	// MsgModelMedicineTask 拿药练习
	MsgModelMedicineTask = 7
	// MsgModelAttendance 考勤
	MsgModelAttendance = 8
	// MsgModelRecall 撤回消息 - 客户端占用
	MsgModelRecall = 9
	// MsgInviteToGroup 邀请好友加入群聊 - 客户端占用
	MsgInviteToGroup = 10
	// MsgModelGroupNotify 群内通知
	MsgModelGroupNotify = 11
	// MsgFirstAcceptFriend 初次统一添加好友 - 客户端占用
	MsgFirstAcceptFriend = 12
	// MsgShareToYST 文章分享到药视通 - 客户端占用
	MsgShareToYST = 13
	// MsgModelSpeaking 药我说
	MsgModelSpeaking = 14
	// MsgModelTask 发布任务
	MsgModelTask = 15
	// MsgModelLevelUp 等级提升
	MsgModelLevelUp = 16
	// MsgModelMedal 获得勋章
	MsgModelMedal = 17
	// MsgModelRedpack 红包通知
	MsgModelRedpack = 18
	// MsgModelDraw 提现通知
	MsgModelDraw = 19
	// MsgModelDrawFailed 提现失败通知
	MsgModelDrawFailed = 20
	// MsgModelRefund 退款通知
	MsgModelRefund = 21
	// MsgModelLeaveCompany 离开企业通知
	MsgModelLeaveCompany = 22
	// MsgModelCompanyExam 企业试卷
	MsgModelCompanyExam = 23
	// MsgModelCompanyCourse 企业课程
	MsgModelCompanyCourse = 24
	// MsgModelIndustryLuckyMoney 工业红包
	MsgModelIndustryLuckyMoney = 25
	// MsgModelUserCompanyApply 用户申请加入企业审核
	MsgModelUserCompanyApply = 26
	// MsgModelUserCompanyApplyResult 申请加入企业审核结果
	MsgModelUserCompanyApplyResult = 27
	// MsgModelLearningMaps 学习地图发布通知
	MsgModelLearningMaps = 28
	// MsgModelFlowing 进店顾客通知
	MsgModelFlowing = 29
	// MsgModelCouponUser 优惠券通知
	MsgModelCouponUser = 30
	// MsgModelCompanyArticle 企业资讯
	MsgModelCompanyArticle = 31
	// MsgModelTeachingOperation 带教操作相关消息
	MsgModelTeachingOperation = 32
	// MsgModelTeachingProgress 带教进行中消息
	MsgModelTeachingProgress = 33
	// MsgModelMedicineTaskStart 拿药练习开始通知
	MsgModelMedicineTaskStart = 34
	// MsgModelMedicineTaskProgress 拿药练习进行中消息
	MsgModelMedicineTaskProgress = 35

	// MsgModelTeachingWait 带教待认领，向客户端输出时转换为32
	MsgModelTeachingWait = 1001
	// MsgModelTeachingStart 带教开始，带教待认领，向客户端输出时转换为32
	MsgModelTeachingStart = 1002
	// MsgModelTeachingUnfinish 带教未完成通知，带教待认领，向客户端输出时转换为33
	MsgModelTeachingUnfinish = 1003
	// MsgModelTeachingEnd 带教结束通知，带教待认领，向客户端输出时转换为33
	MsgModelTeachingFinish = 1004
	// MsgModelTeachingEvaluation 带教评价，带教待认领，向客户端输出时转换为33
	MsgModelTeachingEvaluation = 1005
	// MsgModelMedicineTaskFinish 拿药练习结束通知，仅限带教，向客户端输出时转换为35
	MsgModelMedicineTaskFinish = 1006
	// MsgModelMedicineTaskEvaluation 拿药练习评价，向客户端输出时转换为35
	MsgModelMedicineTaskEvaluation = 1007
)

// CustomMsgTips 自定义消息通知提示信息
var CustomMsgTips = map[int]string{
	MsgModelCommon:                 "您有一条新消息",
	MsgModelArticle:                "您有一条新消息",
	MsgModelCompanyArticle:         "您有一条新的企业资讯",
	MsgModelImage:                  "您有一条新消息",
	MsgModelAnnounce:               "您有一条新公告",
	MsgModelApproval:               "您有一条新的审批通知",
	MsgModelCourseTask:             "您有一项新的课程任务",
	MsgModelMedicineTask:           "您有一项新的拿药练习任务",
	MsgModelAttendance:             "您有一条新的考勤消息",
	MsgModelSpeaking:               "您有一条新的药我说练习",
	MsgModelTask:                   "您有一个新的积分任务待完成",
	MsgModelLevelUp:                "恭喜您的等级提升啦",
	MsgModelMedal:                  "恭喜您获得一枚新的勋章",
	MsgModelRedpack:                "恭喜您获得了一个红包",
	MsgModelDraw:                   "您的提现申请已受理",
	MsgModelDrawFailed:             "提现失败提醒",
	MsgModelRefund:                 "退款到账提醒",
	MsgModelLeaveCompany:           "有员工退出了本企业",
	MsgModelCompanyExam:            "您有一个新的试卷待完成",
	MsgModelCompanyCourse:          "您有一个新的课程待完成",
	MsgModelIndustryLuckyMoney:     "您有一个红包待领取",
	MsgModelUserCompanyApply:       "您收到一个新的店员入驻申请",
	MsgModelUserCompanyApplyResult: "您有一条新的审核通知",
	MsgModelCouponUser:             "您有一条新的优惠券信息",
	MsgModelTeachingWait:           "您的企业为您发布了一个新的带教任务，待认领",
	MsgModelTeachingStart:          "带教开始通知",
	MsgModelMedicineTaskStart:      "拿药练习开始通知",
}

// MsgCommonElem 通用消息
type MsgCommonElem struct {
	Model       int
	ModelName   string // 模型名称
	Title       string // 内容标题
	Description string // 描述信息
	URL         string // 详情地址
}

// MsgArticleItem 文章条目
type MsgArticleItem struct {
	Title     string // 标题
	Digest    string // 描述
	Thumb     string // 缩略图
	Image     string // 大图
	SourceURL string // 文章地址
}

// MsgImageElem 图片消息
type MsgImageElem struct {
	Model       int
	URL         string
	ThumbURL    string
	ThumbWidth  int
	ThumbHeight int
}

// MsgArticleElem 图文消息
type MsgArticleElem struct {
	Model    int
	Articles []*MsgArticleItem
}

// MsgAnnounceElem 公告消息
type MsgAnnounceElem struct {
	Model          int
	Title          string
	Description    string
	Content        string
	Date           string
	URL            string
	AppID          int
	ModuleName     string
	AppVersion     string
	AppDownloadURL string
}

// MsgApprovalElem 审批消息，
// 通过SubModel区分不同审批模型，新增模型后直接扩展对应需要展示的字段
type MsgApprovalElem struct {
	Model          int
	SubModel       int // 子模型
	Title          string
	StartTime      string
	EndTime        string
	LeaveType      string // 请假类型
	Reason         string // 加班、外出原因
	CommonTitle    string // 通用审批标题
	CommonContent  string // 通用审批详情
	Date           string // 提交时间
	Duration       string // 时长
	URL            string // 详情页地址
	AppID          int
	ModuleName     string
	AppVersion     string
	AppDownloadURL string
}

// MsgCourseTaskElem 课程任务
type MsgCourseTaskElem struct {
	Model       int
	Title       string
	Description string
	TaskID      int
	Date        string
	EndDate     string
}

// MsgMedicineTaskElem 拿药练习
type MsgMedicineTaskElem struct {
	Model  int
	Title  string // 内容标题
	Task   string // 任务信息
	Date   string
	TaskID int
}

// MsgAttendanceElem 考勤消息
type MsgAttendanceElem struct {
	Model          int
	Title          string
	Description    string
	Date           string
	URL            string // 详情页地址
	AppID          int
	ModuleName     string
	AppVersion     string
	AppDownloadURL string
}

// MsgSpeakingElem 药我说
type MsgSpeakingElem struct {
	Model       int
	Title       string
	Description string
	Content     string
	Date        string
	EndDate     string
	SpeakingID  int
}

// MsgTaskElem 任务
type MsgTaskElem struct {
	Model     int
	Title     string
	Category  string
	ModelID   int
	ContentID int
}

// MsgLevelUpElem 等级提升
type MsgLevelUpElem struct {
	Model int
	Title string // 内容标题
}

// MsgMedalElem 勋章
type MsgMedalElem struct {
	Model int
	Title string // 内容标题
}

// MsgGroupNotifyElem 群内通知
type MsgGroupNotifyElem struct {
	Model       int
	PeerID      string
	Description string
}

// MsgRedpackElem 红包通知
type MsgRedpackElem struct {
	Model int
	Title string
	Money string
	Date  string
}

// MsgDrawElem 提现记录
type MsgDrawElem struct {
	Model       int
	Money       string
	Tax         string
	ActualMoney string
	Date        string
	Status      string
	Note        string
}

// MsgDrawFailedElem 提现失败
type MsgDrawFailedElem struct {
	Model       int
	Money       string
	Tax         string
	ActualMoney string
	Date        string
	Status      string
	Note        string
}

// MsgRefundElem 退款通知
type MsgRefundElem struct {
	Model int
	Money string
	Date  string
	Note  string
}

// MsgLeaveCompanyElem 离开企业
type MsgLeaveCompanyElem struct {
	Model      int
	Title      string
	Department string
	Job        string
	Date       string
}

// MsgCompanyExamElem 企业试卷
type MsgCompanyExamElem struct {
	Model       int
	Title       string
	ID          int
	Name        string
	SubjectDesc string
	ScoreDesc   string
	LimitTime   string
	Date        string
}

// MsgCompanyCourseElem 企业课程
type MsgCompanyCourseElem struct {
	Model int
	Title string
	ID    int
	Name  string
	Date  string
}

// MsgIndustryLuckyMoneyElem 工业红包通知
type MsgIndustryLuckyMoneyElem struct {
	Model       int
	Title       string
	Description string
	ID          int
	EndDate     string
}

// MsgUserCompanyApplyElem 加入企业申请审核通知
type MsgUserCompanyApplyElem struct {
	Model           int
	Title           string
	DepartmentID    int    // 门店编号
	DepartmentName  string // 门店名称
	ApplyUser       string
	ApplyUserMobile string
	Date            string
}

// MsgUserCompanyApplyResultElem 企业申请审核结果通知
type MsgUserCompanyApplyResultElem struct {
	Model          int
	Title          string
	ResultTitle    string // 审核结果标题
	Result         int    // 审核结果（0 拒绝 1 通过）
	Content        string // 拒绝理由
	DepartmentName string // 门店名称
	Date           string
}

// MsgLearningMapsElem 学习地图发布通知
type MsgLearningMapsElem struct {
	Model       int
	ID          int
	Finished    int
	Title       string
	Description string
	Date        string
}

// MsgFlowingElem 进店客流通知
type MsgFlowingElem struct {
	Model           int
	FlowingID       int
	CustomerCrowdID int
	CustomerID      int
	Name            string
	Sex             string
	Age             int
	Avatar          string
	Tags            string
	OrderMedicines  string
}

// Tag Tag
type Tag struct {
	ID   int
	Name string
}

// OrderMedicine OrderMedicine
type OrderMedicine struct {
	ID              int
	MedicineOrderID int
	MedicineID      int
	MedicineName    string
	MedicineFactory string
	MedicineDosage  string
	MedicineUnit    string
	Count           float32
}

// CouponUserElem 优惠券通知
type CouponUserElem struct {
	Model int
	Title string
}

// MsgCompanyArticleElem 企业资讯消息
type MsgCompanyArticleElem struct {
	Model       int
	Title       string
	Description string
	Date        string
	ArticleID   int
}

// MsgTeachingOperationElem 带教操作通知
type MsgTeachingOperationElem struct {
	Model       int
	Title       string
	Description string
	StartDate   string
	EndDate     string
	Cycle       string // 周期
	TeachingID  int
	RoleType    string
}

// MsgTeachingProgressElem 带教进程通知，未完成/结束/评价
type MsgTeachingProgressElem struct {
	Model       int
	Description string
	TeachingID  int
	RoleType    string
}

// MsgMedicineTaskStartElem 拿药练习开始通知
type MsgMedicineTaskStartElem struct {
	Model          int
	Title          string
	Description    string
	EndTime        string
	MedicineTaskID int
}

// MsgMedicineTaskProgressElem 拿药练习进行中通知
type MsgMedicineTaskProgressElem struct {
	Model          int
	Description    string
	MedicineTaskID int
	RoleType       string
}
