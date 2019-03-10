package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCamelCase1(t *testing.T) {
	Convey("将一个字符串转为大驼峰命名", t, func() {
		So(CamelCase("user name"), ShouldEqual, "UserName")
		So(CamelCase("user_name"), ShouldEqual, "UserName")
		So(CamelCase("userName"), ShouldEqual, "UserName")
		So(CamelCase("UserName"), ShouldEqual, "UserName")
		So(CamelCase("User name"), ShouldEqual, "UserName")
		So(CamelCase("User Name"), ShouldEqual, "UserName")
		So(CamelCase("User-Name"), ShouldNotEqual, "UserName")
	})
}

func TestCamelCaseInitialism1(t *testing.T) {
	Convey("强制首字母缩写命名规范转换", t, func() {
		So(CamelCaseInitialism("user name"), ShouldEqual, "UserName")
		So(CamelCaseInitialism("userid"), ShouldEqual, "Userid")
		So(CamelCaseInitialism("user_id"), ShouldEqual, "UserID")
		So(CamelCaseInitialism("username"), ShouldEqual, "Username")
		So(CamelCaseInitialism("jsonrpc"), ShouldEqual, "Jsonrpc")
		So(CamelCaseInitialism("json_rpc"), ShouldEqual, "JSONRPC")
		So(CamelCaseInitialism("UserIP"), ShouldEqual, "UserIP")
		So(CamelCaseInitialism("User-ID"), ShouldNotEqual, "UserID")
	})
}

func TestStrPad1(t *testing.T) {
	Convey("使用另一个字符串填充字符串为指定长度", t, func() {
		So(StrPad(1, 9, "0", STR_PAD_LEFT), ShouldEqual, "000000001")
		So(StrPad("1", 9, "0", STR_PAD_LEFT), ShouldEqual, "000000001")
		So(StrPad("1", 9, "0", STR_PAD_RIGHT), ShouldEqual, "100000000")
		So(StrPad("1", 9, "0", STR_PAD_BOTH), ShouldEqual, "000010000")
	})
}

func TestStrNatCut(t *testing.T) {
	Convey("使用自然的方式截取字符串", t, func() {
		s := "我是123的abc的456吗？会"
		So(StrNatCut(s, 3), ShouldEqual, "我是12...")
		So(StrNatCut(s, 5, ""), ShouldEqual, "我是123的a")
		So(StrNatCut("测试长度相等", 6), ShouldEqual, "测试长度相等")
		So(StrNatCut(s, 100), ShouldEqual, s)
	})
}

func TestIsNumberic(t *testing.T) {
	Convey("测试给定的字符串是否是纯数字", t, func() {
		So(IsNumeric(""), ShouldBeFalse)
		So(IsNumeric("123a"), ShouldBeFalse)
		So(IsNumeric("123"), ShouldBeTrue)
	})
}

func TestConcat(t *testing.T) {
	Convey("连接字符串", t, func() {
		So(Concat(), ShouldEqual, "")
		So(Concat(""), ShouldEqual, "")
		So(Concat("a"), ShouldEqual, "a")
		So(Concat("a", "b"), ShouldEqual, "ab")
		So(Concat("a", "", "b"), ShouldEqual, "ab")
	})
}
