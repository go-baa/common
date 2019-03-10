package image

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	"github.com/sillydong/fastimage"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

// Info 图片信息
type Info struct {
	Type   string // bmp,jpg,png,gif
	Width  int
	Height int
}

// Detect 返回一个图片的类型，宽和高
func Detect(file string) (*Info, error) {
	// 判断远程图片
	if strings.Contains(file, "://") {
		d := fastimage.NewFastImage(3, nil)
		it, is, err := d.Detect(file)
		if err != nil {
			return nil, err
		}
		return &Info{Type: imageTypes[it], Width: int(is.Width), Height: int(is.Height)}, nil
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	img, format, err := image.DecodeConfig(f)
	if err != nil {
		return nil, err
	}

	return &Info{Type: strings.ToUpper(format), Width: img.Width, Height: img.Height}, nil
}

var imageTypes map[fastimage.ImageType]string

func init() {
	imageTypes = make(map[fastimage.ImageType]string)
	imageTypes[fastimage.JPEG] = "JPEG"
	imageTypes[fastimage.PNG] = "PNG"
	imageTypes[fastimage.GIF] = "GIF"
	imageTypes[fastimage.BMP] = "BMP"
	imageTypes[fastimage.WEBP] = "WEBP"
	imageTypes[fastimage.Unknown] = "Unkown"
}
