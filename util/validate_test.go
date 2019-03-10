package util

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMobile(t *testing.T) {
	Convey("测试手机号码", t, func() {
		So(Validate("13683216194", "mobile"), ShouldBeTrue)
	})
}
