package pinyin

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPinyinCompare(t *testing.T) {
	Convey("测试两个汉字的拼音是否相同", t, func() {
		Convey("药 = 要", func() {
			So(ComparePinyinByRune('药', '要'), ShouldBeTrue)
		})
		Convey("换 = 还", func() {
			So(ComparePinyinByRune('换', '还'), ShouldBeTrue)
		})
		Convey("海 = 还", func() {
			So(ComparePinyinByRune('海', '还'), ShouldBeTrue)
		})
		Convey("飞 =/= 黑", func() {
			So(ComparePinyinByRune('飞', '黑'), ShouldBeFalse)
		})
	})
}
