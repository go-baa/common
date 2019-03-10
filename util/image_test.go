package util

import (
	"fmt"
	"image/color"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/image/colornames"
)

func TestGenerateTextAvatar(t *testing.T) {
	Convey("测试生成文字头像", t, func() {
		file := os.TempDir() + "gen_avatar.png"
		err := GenerateTextAvatar(file, "杨恒飞", nil, "", 300)
		So(err, ShouldBeNil)
		os.Remove(file)
	})
}

func TestColorHex(t *testing.T) {
	Convey("测试16进制颜色转换RBGA", t, func() {
		Convey("绿色", func() {
			rgbColor := colornames.Green
			hexColor, _ := HexColor2RGBA(ColorGreen)
			So(color2Text(hexColor), ShouldEqual, color2Text(rgbColor))
		})

		Convey("橙色", func() {
			rgbColor := colornames.Orange
			hexColor, _ := HexColor2RGBA(ColorOrange)
			So(color2Text(hexColor), ShouldEqual, color2Text(rgbColor))
		})
	})

	Convey("测试RGBA转换16进制颜色", t, func() {
		Convey("绿色", func() {
			rgbColor := colornames.Green
			hexColor := RGBA2HexColor(rgbColor)
			So(hexColor, ShouldEqual, ColorGreen)
		})

		Convey("橙色", func() {
			rgbColor := colornames.Orange
			hexColor := RGBA2HexColor(rgbColor)
			So(hexColor, ShouldEqual, ColorOrange)
		})
	})

}
func color2Text(c color.RGBA) string {
	return fmt.Sprintf("%v", c)
}
