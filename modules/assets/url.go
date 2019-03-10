package assets

import (
	"net/url"
	"strings"

	"github.com/go-baa/setting"
)

// AbsoluteUploadURL ...
func AbsoluteUploadURL(rawURL string) string {
	if len(rawURL) == 0 {
		return rawURL
	}

	uploadURI := setting.Config.MustString("upload.baseUri", "")
	if len(uploadURI) == 0 {
		return rawURL
	}

	if strings.HasPrefix(rawURL, uploadURI) {
		return rawURL
	}

	url, err := url.Parse(rawURL)
	if err != nil || len(url.Host) > 0 {
		return rawURL
	}

	return strings.TrimRight(uploadURI, "/") + "/" + strings.TrimLeft(rawURL, "/")
}

// RelativeUploadURL ...
func RelativeUploadURL(rawURL string) string {
	if len(rawURL) == 0 {
		return rawURL
	}

	uploadURI := setting.Config.MustString("upload.baseUri", "")
	if len(uploadURI) == 0 {
		return rawURL
	}

	if strings.HasPrefix(rawURL, uploadURI) {
		return strings.TrimPrefix(rawURL, uploadURI)
	}

	return rawURL
}

// AbsoluteAssetsURL ...
func AbsoluteAssetsURL(rawURL string) string {
	if len(rawURL) == 0 {
		return rawURL
	}

	assetsURI := setting.Config.MustString("assets.baseUri", "")
	if len(assetsURI) == 0 {
		return rawURL
	}

	if strings.HasPrefix(rawURL, assetsURI) {
		return rawURL
	}

	url, err := url.Parse(rawURL)
	if err != nil || len(url.Host) > 0 {
		return rawURL
	}

	return strings.TrimRight(assetsURI, "/") + "/" + strings.TrimLeft(rawURL, "/")
}

// RelativeAssetsURL ...
func RelativeAssetsURL(rawURL string) string {
	if len(rawURL) == 0 {
		return rawURL
	}

	assetsURI := setting.Config.MustString("assets.baseUri", "")
	if len(assetsURI) == 0 {
		return rawURL
	}

	if strings.HasPrefix(rawURL, assetsURI) {
		return strings.TrimPrefix(rawURL, assetsURI)
	}

	return rawURL
}
