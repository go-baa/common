package util

import (
	"testing"

	"os"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMD5(t *testing.T) {
	Convey("测试MD5", t, func() {
		So(MD5("123456"), ShouldEqual, "e10adc3949ba59abbe56e057f20f883e")
	})
}

func TestMD5File(t *testing.T) {
	Convey("", t, func() {
		f := os.TempDir() + "/testmd5file.tmp"
		_, err := WriteFile(f, []byte("123456"))
		So(err, ShouldBeNil)
		So(MD5("123456"), ShouldEqual, "e10adc3949ba59abbe56e057f20f883e")
	})
}
