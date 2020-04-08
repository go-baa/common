package excel

import (
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExcelRead1(t *testing.T) {
	Convey("测试Excel读取", t, func() {
		e, err := NewReader("example/店面导入模板.xlsx", []*HeaderColumn{
			&HeaderColumn{Field: "Name", Title: "店面名称"},
			&HeaderColumn{Field: "Concact", Title: "店长"},
			&HeaderColumn{Field: "Mobile", Title: "手机"},
			&HeaderColumn{Field: "Tel", Title: "电话"},
			&HeaderColumn{Field: "Fax", Title: "传真"},
			&HeaderColumn{Field: "ProvinceName", Title: "地区（省）"},
			&HeaderColumn{Field: "CityName", Title: "地区（市）"},
			&HeaderColumn{Field: "AreaName", Title: "地区（区/县）"},
			&HeaderColumn{Field: "Address", Title: "地址"},
			&HeaderColumn{Field: "Zip", Title: "邮编"},
			&HeaderColumn{Field: "Description", Title: "描述"},
		}, 0, false)

		So(err, ShouldBeNil)

		row, err := e.Fetch()
		So(row, ShouldNotBeNil)
	})
}

func TestExcelRead2(t *testing.T) {
	Convey("测试Excel读取错误的内容", t, func() {
		e, err := NewReader("example/店面导入模板.xlsx", []*HeaderColumn{
			&HeaderColumn{Field: "Name", Title: "店面名称"},
			&HeaderColumn{Field: "Concact", Title: "店长2", Require: true},
			&HeaderColumn{Field: "Mobile", Title: "手机"},
		}, 0, false)

		So(err, ShouldNotBeNil)
		So(e, ShouldBeNil)
	})
}

func TestExcelRead3(t *testing.T) {
	Convey("测试Excel读取带注释的行", t, func() {
		e, err := NewReader("example/店员导入模板.xlsx", []*HeaderColumn{
			&HeaderColumn{Field: "Department", Title: "店面名称"},
			&HeaderColumn{Field: "Name", Title: "店员名称"},
			&HeaderColumn{Field: "Mobile", Title: "手机号"},
			&HeaderColumn{Field: "Sex", Title: "性别"},
			&HeaderColumn{Field: "Address", Title: "地址"},
		}, 0, false)

		So(err, ShouldBeNil)

		var i int
		for {
			row, _ := e.Fetch()
			if row == nil {
				break
			}
			fmt.Println(i, row)
			i++
		}
	})
}
func TestExcelRead4(t *testing.T) {
	Convey("测试Excel读取带说明头部的行", t, func() {
		e, err := NewReader("example/带说明的模板.xlsx", []*HeaderColumn{
			&HeaderColumn{Field: "Department", Title: "店面名称"},
			&HeaderColumn{Field: "Name", Title: "店员名称"},
			&HeaderColumn{Field: "Mobile", Title: "手机号"},
			&HeaderColumn{Field: "Sex", Title: "性别"},
			&HeaderColumn{Field: "Address", Title: "地址"},
		}, 3, false)

		So(err, ShouldBeNil)

		var i int
		for {
			row, _ := e.Fetch()
			if row == nil {
				break
			}
			fmt.Println(i, row)
			i++
		}
	})
}

func TestExcelRead5(t *testing.T) {
	Convey("测试Excel读取带大段说明头部的行", t, func() {
		e, err := NewReader("example/试题导入模板.xlsx", []*HeaderColumn{
			&HeaderColumn{Field: "type", Title: "题型（必填）"},
			&HeaderColumn{Field: "title", Title: "题目（必填）"},
			&HeaderColumn{Field: "score", Title: "分数（必填）"},
			&HeaderColumn{Field: "correct", Title: "正确答案（必填）"},
			&HeaderColumn{Field: "explain", Title: "试题解析"},
			&HeaderColumn{Field: "optiona", Title: "选项A（必填）"},
			&HeaderColumn{Field: "optionb", Title: "选项B（必填）"},
			&HeaderColumn{Field: "optionc", Title: "选项C"},
			&HeaderColumn{Field: "optiond", Title: "选项D"},
			&HeaderColumn{Field: "optione", Title: "选项E"},
			&HeaderColumn{Field: "optionf", Title: "选项F"},
			&HeaderColumn{Field: "optiong", Title: "选项G"},
		}, 10, false)

		So(err, ShouldBeNil)

		var i int
		for {
			row, _ := e.Fetch()
			if row == nil {
				break
			}
			fmt.Println(i, row)
			i++
		}
	})
}

func TestExcelWrite1(t *testing.T) {
	Convey("测试Excel写入", t, func() {
		e, err := NewWriter()

		So(err, ShouldBeNil)

		e.SetHeader([]*HeaderColumn{
			&HeaderColumn{Field: "Name", Title: "店员名称"},
			&HeaderColumn{Field: "Concact", Title: "店长"},
			&HeaderColumn{Field: "Mobile", Title: "手机"},
			&HeaderColumn{Field: "Tel", Title: "电话"},
			&HeaderColumn{Field: "Fax", Title: "传真"},
			&HeaderColumn{Field: "ProvinceName", Title: "地区（省）"},
			&HeaderColumn{Field: "CityName", Title: "地区（市）"},
			&HeaderColumn{Field: "AreaName", Title: "地区（区/县）"},
			&HeaderColumn{Field: "Address", Title: "地址"},
			&HeaderColumn{Field: "Zip", Title: "邮编"},
			&HeaderColumn{Field: "Description", Title: "描述"},
		}, nil)

		err = e.Write(map[string]interface{}{
			"Name":         "测试店1",
			"Concact":      "店长",
			"Mobile":       "18611111111",
			"Tel":          "010-11111111",
			"Fax":          "010-111",
			"ProvinceName": "北京",
			"CityName":     "市辖区",
			"AreaName":     "海淀区",
			"Address":      "北京 海淀 西山",
			"Zip":          "100000",
			"Description":  "描述描述描述",
		})
		So(err, ShouldBeNil)

		err = e.SaveToFile("example/output.xlsx")
		So(err, ShouldBeNil)

		os.Remove("example/output.xlsx")
	})
}
