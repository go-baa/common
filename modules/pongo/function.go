package pongo

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"strconv"
	"strings"

	"git.code.tencent.com/xinhuameiyu/common/util"
	"github.com/go-baa/setting"
	baa "gopkg.in/baa.v1"
)

// 解析VUE构建文件
var manifest map[string]string

// Functions 模板函数库
func Functions(b *baa.Baa) map[string]interface{} {
	return map[string]interface{}{
		"config":        config,
		"assets":        assets,
		"assetsUrl":     assetsURL,
		"stylesheetTag": stylesheetTag,
		"scriptTag":     scriptTag,
		"url": func(name string, params ...interface{}) string {
			n := make([]interface{}, len(params))
			for k, v := range params {
				n[k] = fmt.Sprint(v)
			}
			return b.URLFor(name, n...)
		},
		"strcut": func(str string, length int, dot string) string {
			return util.StrNatCut(str, length, dot)
		},
		"pages": pages,
		"range": func(start, end int) []int {
			ret := []int{}
			for start <= end {
				ret = append(ret, start)
				start++
			}
			return ret
		},
		"replace":    replace,
		"replaceInt": replaceInt,
		"stripTags":  stripTags,
	}
}

func config(args ...string) string {
	key := args[0]
	dft := ""
	if len(args) == 2 {
		dft = args[1]
	}
	if len(key) == 0 {
		return dft
	}
	return setting.Config.MustString(key, dft)
}

func assets(path string) string {
	if val, ok := manifest[path]; ok {
		return setting.Config.MustString("assets.baseUri", "") + val
	}
	return ""
}

func assetsURL(str ...string) string {
	path := strings.Join(str, "")
	if val, ok := manifest["hash"]; ok {
		if strings.Contains(path, "?") {
			return path + "&" + val
		}
		return path + "?" + val
	}
	return path
}

func stylesheetTag(name string) template.HTML {
	uri := assets(name)
	tag := ""
	if len(uri) > 0 {
		tag = "<link rel=\"stylesheet\" href=\"" + uri + "\">"
	}
	return template.HTML(tag)
}

func scriptTag(name string) template.HTML {
	uri := assets(name)
	tag := ""
	if len(uri) > 0 {
		tag = "<script src=\"" + uri + "\"></script>"
	}
	return template.HTML(tag)
}

type pagesResult struct {
	Page  int
	Total int
	Prev  int
	Next  int
	Pages []int
	Start int
	Stop  int
}

func pages(items, page, pagesize, visible int) pagesResult {
	ret := pagesResult{}
	ret.Pages = []int{}

	if page < 1 {
		page = 1
	}
	ret.Page = page

	// 分页数不合法
	if pagesize < 1 {
		return ret
	}

	// 去除首页和末页，最少 1 个页码可见
	if visible < 3 {
		visible = 3
	}

	// 只有一页，无需处理
	ret.Total = int(math.Ceil(float64(items) / float64(pagesize)))
	if ret.Page > ret.Total {
		ret.Page = ret.Total
	}
	if ret.Total < 2 {
		return ret
	}

	// 上一页
	ret.Prev = ret.Page - 1
	if ret.Prev < 1 {
		ret.Prev = 1
	}

	// 下一页
	ret.Next = ret.Page + 1
	if ret.Next > ret.Total {
		ret.Next = ret.Total
	}

	// 只有两页的时候特殊处理
	if ret.Total == 2 {
		return ret
	}

	// 页码总数刚好等于可见数
	if visible > ret.Total {
		visible = ret.Total
	}
	if ret.Total == visible {
		start := 2
		for start < ret.Total {
			ret.Pages = append(ret.Pages, start)
			start++
		}
	} else {
		visible -= 2
		offset := (visible - 1) / 2

		// 如果当前页码靠近总页码
		var start, stop int
		if ret.Page+visible >= ret.Total {
			stop = ret.Total - 1
			start = stop - visible + 1
		} else {
			start = ret.Page - offset
			if start < 2 {
				start = 2
			}
			stop = start + visible - 1
		}
		if start < 2 {
			start = 2
		}
		if stop >= ret.Total {
			stop = ret.Total - 1
		}
		for i := start; i <= stop; i++ {
			ret.Pages = append(ret.Pages, i)
		}
	}

	ret.Start = ret.Pages[0]
	ret.Stop = ret.Pages[len(ret.Pages)-1]

	return ret
}

func replace(subject, search, replace string) string {
	return strings.Replace(subject, search, replace, -1)
}

func replaceInt(subject, search string, replace int) string {
	return strings.Replace(subject, search, strconv.Itoa(replace), -1)
}

func stripTags(str string, tags ...string) string {
	return string(util.StripTags([]byte(str), tags...))
}

func loadManifest() {
	data, err := util.ReadFile(setting.Config.MustString("assets.buildPath", "") + "rev-manifest.json")
	if err == nil {
		json.Unmarshal(data, &manifest)
	}
}

func init() {
	loadManifest()
}
