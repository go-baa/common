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

func TestImageCutWithFormat(t *testing.T) {
	te := []string{
		"5c8b0a832f6501601.jpg_e680.webp",
		"5cbb2969ab7757923.jpg_e680.png",
		"5c8b0a832f6501601.jpg_e680.jpg",
		"5c6bfd26ab24a3846.jpg_e680 (3).tiff",
		"timg.gif",
		"unk.jpg",
	}
	Convey("测试识别图片格式剪裁", t, func() {
		for i, f := range te {
			in := "/tmp/" + f
			out := "/tmp/" + fmt.Sprintf("%d.jpg", i)
			err := ImageCutWithFormat(in, out, 0, 0, 300, 300)
			So(err, ShouldBeNil)
		}
	})
}

func TestImageFormatDiscernment(t *testing.T) {
	te := map[string]string{
		"PNG": "5cbb2969ab7757923.jpg_e680.png",
		// "JPEG": "5c8b0a832f6501601.jpg_e680.jpg",
		"WEBP": "5c8b0a832f6501601.jpg_e680.webp",
		"GIF":  "timg.gif",
		"TIFF": "5c6bfd26ab24a3846.jpg_e680 (3).tiff",
		"JPEG": "unk.jpg",
	}
	Convey("测试识别图片格式识别", t, func() {
		for ft, src := range te {
			f, _ := os.Open("/tmp/" + src)
			format, err := ImageFormatDiscernment(f)
			So(err, ShouldBeNil)
			So(format, ShouldEqual, ft)
		}
	})
}

func TestImageSizeWithFormat(t *testing.T) {
	te := []string{
		"5cbb2969ab7757923.jpg_e680.png",
		"5c8b0a832f6501601.jpg_e680.jpg",
		"5c8b0a832f6501601.jpg_e680.webp",
		"timg.gif",
		"5c6bfd26ab24a3846.jpg_e680 (3).tiff",
	}
	Convey("测试识别图片格式识别", t, func() {
		for _, src := range te {
			h, w, err := ImageSizeWithFormat("/tmp/" + src)
			fmt.Println(h)
			fmt.Println(w)
			So(err, ShouldBeNil)
		}
	})
}
