package base

// SMSProvider 短信提供商
type SMSProvider interface {
	// SendSMSCode 发送短信验证码
	SendSMSCode(mobile string, code string) (string, error)
}

// VoiceProvider 语音短信提供商
type VoiceProvider interface {
	// SendVoiceCode 发送语音验证码
	SendVoiceCode(mobile string, code string) (string, error)
}

// Provider 提供商存储器
var _store map[string]interface{}

func init() {
	_store = make(map[string]interface{})
}

// Register 注册一个提供商
func Register(name string, p interface{}) {
	_store[name] = p
}

// GetSMSProvider 获取一个短信提供商
func GetSMSProvider(name string) SMSProvider {
	p := _store[name]
	if p == nil {
		return nil
	}
	if v, ok := p.(SMSProvider); ok {
		return v
	}
	return nil
}

// GetVoiceProvider 获取一个语音短信提供商
func GetVoiceProvider(name string) VoiceProvider {
	p := _store[name]
	if p == nil {
		return nil
	}
	if v, ok := p.(VoiceProvider); ok {
		return v
	}
	return nil
}
