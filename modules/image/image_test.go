package image

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestImageType(t *testing.T) {
	Convey("测试图片类型", t, func() {
		info, err := Detect("http://upload.vodjk.com/2013/1202/1385947610696.jpg")
		So(err, ShouldBeNil)
		So(info.Type, ShouldEqual, "JPEG")
	})
}
func TestImageSize(t *testing.T) {
	Convey("测试相等比较", t, func() {
		info, err := Detect("http://upload.vodjk.com/2013/1202/1385947610696.jpg")
		So(err, ShouldBeNil)
		So(info.Width, ShouldEqual, 329)
	})
}
