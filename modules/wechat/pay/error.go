package pay

const (
	// ErrNoAuth 无权限
	ErrNoAuth = "NOAUTH"
	// ErrNotEnough 余额不足
	ErrNotEnough = "NOTENOUGH"
	// ErrOrderPaid 商户订单已支付
	ErrOrderPaid = "ORDERPAID"
	// ErrOrderClosed 订单已关闭
	ErrOrderClosed = "ORDERCLOSED"
	// ErrSystemError 系统错误
	ErrSystemError = "SYSTEMERROR"
	// ErrAppidNotExist APPID不存在系统错误
	ErrAppidNotExist = "APPID_NOT_EXIST"
	// ErrMchidNotExist MCHID不存在
	ErrMchidNotExist = "MCHID_NOT_EXIST"
	// ErrAppidMchidNotMatch appid和mch_id不匹配
	ErrAppidMchidNotMatch = "APPID_MCHID_NOT_MATCH"
	// ErrLackParams 缺少参数
	ErrLackParams = "LACK_PARAMS"
	// ErrOutTradeNoUsed 商户订单号重复
	ErrOutTradeNoUsed = "OUT_TRADE_NO_USED"
	// ErrSignError 签名错误
	ErrSignError = "SIGNERROR"
	// ErrXMLFormatError XML格式错误
	ErrXMLFormatError = "XML_FORMAT_ERROR"
	// ErrRequirePostMethod 请使用post方法
	ErrRequirePostMethod = "REQUIRE_POST_METHOD"
	// ErrPostDataEmpty post数据为空
	ErrPostDataEmpty = "POST_DATA_EMPTY"
	// ErrNotUTF8 	编码格式错误
	ErrNotUTF8 = "NOT_UTF8"
	// ErrOrderNotExist 此交易订单号不存在
	ErrOrderNotExist = "ORDERNOTEXIST"
	// ErrBizerrNeedRetry 退款业务流程错误
	ErrBizerrNeedRetry = "BIZERR_NEED_RETRY"
	// ErrTradeOverdue 订单已经超过退款期限
	ErrTradeOverdue = "TRADE_OVERDUE"
	// ErrError 业务错误
	ErrError = "ERROR"
	// ErrUserAccountAbnormal 退款请求失败
	ErrUserAccountAbnormal = "USER_ACCOUNT_ABNORMAL"
	// ErrInvalidReqTooMuch 无效请求过多
	ErrInvalidReqTooMuch = "INVALID_REQ_TOO_MUCH"
	// ErrInvalidTransactionID 无效transaction_id
	ErrInvalidTransactionID = "INVALID_TRANSACTIONID"
	// ErrFrequencyLimited 频率限制
	ErrFrequencyLimited = "FREQUENCY_LIMITED"
	// ErrRefundNotExist 退款订单查询失败
	ErrRefundNotExist = "REFUNDNOTEXIST"
	// ErrNoComment 对应的时间段没有用户的评论数据
	ErrNoComment = "NO_COMMENT"
	// ErrTimeExpire 拉取的时间超过3个月
	ErrTimeExpire = "TIME_EXPIRE"
)

// ErrDesc 错误描述
var ErrDesc = map[string]string{
	ErrNoAuth:               "无接口权限",
	ErrNotEnough:            "余额不足",
	ErrOrderPaid:            "订单已支付",
	ErrOrderClosed:          "订单已关闭",
	ErrSystemError:          "系统错误",
	ErrAppidNotExist:        "APPID不存在",
	ErrMchidNotExist:        "MCHID不存在",
	ErrAppidMchidNotMatch:   "appid和mch_id不匹配",
	ErrLackParams:           "缺少参数",
	ErrOutTradeNoUsed:       "订单号重复",
	ErrSignError:            "签名错误",
	ErrXMLFormatError:       "XML格式错误",
	ErrRequirePostMethod:    "要求post方法",
	ErrPostDataEmpty:        "post数据为空",
	ErrNotUTF8:              "编码格式错误",
	ErrOrderNotExist:        "订单号不存在",
	ErrBizerrNeedRetry:      "退款业务流程错误，需要商户触发重试来解决",
	ErrTradeOverdue:         "订单已经超过退款期限",
	ErrError:                "业务错误",
	ErrUserAccountAbnormal:  "退款请求失败",
	ErrInvalidReqTooMuch:    "无效请求过多",
	ErrInvalidTransactionID: "无效transaction_id",
	ErrFrequencyLimited:     "频率限制",
	ErrRefundNotExist:       "退款订单查询失败",
	ErrNoComment:            "对应的时间段没有用户的评论数据",
	ErrTimeExpire:           "拉取的时间超过3个月",
}
