package errors

var (
	// ErrUserNeedLogin 用户未登录
	ErrUserNeedLogin = New(1000, "用户未登录")
	// ErrUserDisable 用户已被禁用
	ErrUserDisable = New(1001, "用户已被禁用")
	// ErrUserNotExist 用户不存在
	ErrUserNotExist = New(1002, "用户不存在")
	// ErrUserInvalidPass 用户名或密码错误
	ErrUserInvalidPass = New(1003, "用户名或密码错误")
	// ErrUserInvalidMobile 手机号码格式不正确
	ErrUserInvalidMobile = New(1004, "手机号码格式不正确")
	// ErrUserCodeSendTooFrequently 您的验证码请求频率太过频繁
	ErrUserCodeSendTooFrequently = New(1005, "您的验证码请求频率太过频繁")
	// ErrUserCodeSendFailed 验证码发送失败
	ErrUserCodeSendFailed = New(1006, "验证码发送失败")
	// ErrUserNotBind 用户未绑定
	ErrUserNotBind = New(1007, "用户未绑定")
	// ErrUserInvalidEmail 邮箱格式不正确
	ErrUserInvalidEmail = New(1008, "邮箱格式不正确")
	// ErrUserUploadInvalid 上传出现错误
	ErrUserUploadInvalid = New(1009, "上传出现错误")
	// ErrUserInvalidOldPassword 原密码错误
	ErrUserInvalidOldPassword = New(1010, "原密码错误")
	// ErrUserUploadEmpty 没有文件被上传
	ErrUserUploadEmpty = New(1011, "没有文件被上传")
	// ErrUserCodeValidateFailed 验证码不正确
	ErrUserCodeValidateFailed = New(1012, "验证码不正确")
	// ErrUserCodeValidateOver 验证码已过期
	ErrUserCodeValidateOver = New(1013, "验证码已过期")
	// ErrUserEmpty 用户不能为空
	ErrUserEmpty = New(1013, "用户不能为空")
	// ErrUserBindMobileExist 该手机号已绑定其它用户
	ErrUserBindMobileExist = New(1020, "该手机号已绑定其它用户")
	// ErrUserBindEmailExist 该邮箱已绑定其它用户
	ErrUserBindEmailExist = New(1021, "该邮箱已绑定其它用户")
	// ErrUserEmailExist 该邮箱已存在
	ErrUserEmailExist = New(1022, "该邮箱已存在")
	// ErrUserMobileExist 该手机号码已存在
	ErrUserMobileExist = New(1023, "该手机号码已存在")
	// ErrUserMobileFailed 该手机号码绑定失败
	ErrUserMobileFailed = New(1023, "该手机号码绑定失败")
	// ErrUserPasswordFailed 密码修改失败
	ErrUserPasswordFailed = New(1023, "密码修改失败")
	// ErrUserInvalidActionAuth 无操作权限
	ErrUserInvalidActionAuth = New(1026, "无操作权限")
	// ErrUserPasswdCodeEmpty 密码或验证码不能为空
	ErrUserPasswdCodeEmpty = New(1027, "密码或验证码不能为空")
	// ErrUserCodeEmpty 验证码不能为空
	ErrUserCodeEmpty = New(1028, "验证码不能为空")
	// ErrUserWechatBindFailed 微信绑定失败
	ErrUserWechatBindFailed = New(1029, "微信绑定失败")
	// ErrUserWechatUNBindFailed 微信解除绑定失败
	ErrUserWechatUNBindFailed = New(1029, "微信解除绑定失败")
	// ErrUserCollecEmpty 用户收藏为空
	ErrUserCollecEmpty = New(1101, "用户收藏为空")
	// ErrUserCollectDelete 用户删除收藏失败
	ErrUserCollectDelete = New(1102, "用户删除收藏失败")
	// ErrUserCollecFailed 用户收藏失败
	ErrUserCollecFailed = New(1103, "用户收藏失败")
	// ErrUserHistoryLogEmpty 用户浏览记录为空
	ErrUserHistoryLogEmpty = New(1104, "用户浏览记录为空")
	// ErrUserHistoryLogDelete 用户删除收藏失败
	ErrUserHistoryLogDelete = New(1105, "用户删除浏览记录失败")
	// ErrUserInfoChanged 用户信息有变更，请重新登录
	ErrUserInfoChanged = New(10000, "用户信息有变更，请重新登录")
)
