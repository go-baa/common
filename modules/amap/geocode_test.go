package amap

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegeo(t *testing.T) {

	Convey("逆地理编码测试", t, func() {
		point := Point{Latitude: "34.228352", Longitude: "108.8825"}
		address := "陕西省西安市雁塔区丈八沟街道科技一路"
		Convey("字符相等测试", func() {
			res, err := Regeo(RegeoParams{Location: point})
			So(err, ShouldBeNil)
			So(res.FormattedAddress, ShouldEqual, address)
		})
	})
}
