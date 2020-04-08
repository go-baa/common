package pinyin

import pinyin "github.com/mozillazg/go-pinyin"

// StringToPinyin 获取一个字符串的拼音，关闭多音字模式
func StringToPinyin(s string) [][]string {
	engine := pinyin.NewArgs()
	return pinyin.Pinyin(s, engine)
}

// GetPinyinByString 获取一个字符串的拼音，开启多音字模式，每个汉字返回一个拼音切片
func GetPinyinByString(s string) [][]string {
	engine := pinyin.NewArgs()
	engine.Heteronym = true
	return pinyin.Pinyin(s, engine)
}

// GetPinyinByRune 获取一个汉字的拼音，开启多音字模式，返回可能的多个拼音
func GetPinyinByRune(c rune) []string {
	engine := pinyin.NewArgs()
	engine.Heteronym = true
	return pinyin.SinglePinyin(c, engine)
}

// ComparePinyinByRune 计算两个汉字的读音是否相同
func ComparePinyinByRune(c1 rune, c2 rune) bool {
	engine := pinyin.NewArgs()
	engine.Heteronym = true
	c1py := pinyin.SinglePinyin(c1, engine)
	c2py := pinyin.SinglePinyin(c2, engine)
	hasSamePinyin := false
	for _, p1 := range c1py {
		for _, p2 := range c2py {
			if p1 == p2 {
				hasSamePinyin = true
				break
			}
		}
	}
	return hasSamePinyin
}
