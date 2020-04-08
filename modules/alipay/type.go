package alipay

const (
	alipayGateway = "https://openapi.alipay.com/gateway.do"
	oauthTokenAPI = "alipay.system.oauth.token"
	userInfoAPI   = "alipay.user.info.share"
	scopeAuthUser = "auth_user"
	authBase      = "auth_base"
)

// APIRequestTimeout 请求超时时间
var APIRequestTimeout = 10

type errorResponse struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

type oauthTokenResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ReExpiresIn  int    `json:"re_expires_in"`
}

type userInfoResponse struct {
	UserID   string `json:"user_id"`
	Avatar   string `json:"avatar"`
	Province int    `json:"province"`
	City     string `json:"city"`
	Nickname int    `json:"nick_name"`
	Gender   int    `json:"gender"`

	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

type bizOauthTokenContent struct {
	GrantType string `json:"grant_type"`
	Code      string `json:"code"`
}
