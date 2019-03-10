package util

import "strings"

// UserAgentInfo 根据UA获取系统和浏览器
func UserAgentInfo(ua string) (os string, browser string) {
	if strings.Contains(ua, "Android") {
		os = "Android"
	} else if strings.Contains(ua, "iPhone") {
		os = "iOS"
	} else {
		os = "other"
	}

	if strings.Contains(ua, "baiduboxapp") {
		browser = "手机百度"
	} else if strings.Contains(ua, "Baiduspider-render") {
		browser = "百度Render"
	} else if strings.Contains(ua, "MiuiBrowser") {
		browser = "miui浏览器"
	} else if strings.Contains(ua, "MQQBrowser") {
		browser = "QQ浏览器"
	} else if strings.Contains(ua, "UCBrowser") {
		browser = "UC浏览器"
	} else if strings.Contains(ua, "MicroMessenger") {
		browser = "微信"
	} else if strings.Contains(ua, "OppoBrowser") {
		browser = "oppo浏览器"
	} else if strings.Contains(ua, "VivoBrowser") {
		browser = "vivo浏览器"
	} else if strings.Contains(ua, "SogouMobileBrowser") {
		browser = "搜狗浏览器"
	} else if strings.Contains(ua, "Safari") && strings.Contains(ua, "iPhone") {
		browser = "Safari"
	} else if strings.Contains(ua, "MZBrowser") {
		browser = "魅族浏览器"
	} else if strings.Contains(ua, "HUAWEI") {
		browser = "华为浏览器"
	} else if strings.Contains(ua, "SamsungBrowser") {
		browser = "三星浏览器"
	}

	return
}
