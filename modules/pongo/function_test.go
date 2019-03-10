package pongo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPages(t *testing.T) {
	Convey("测试分页页码计算：数据少于 1 页", t, func() {
		ret := pages(2, 2, 10, 1)
		So(ret.Page, ShouldEqual, 1)
		So(ret.Total, ShouldEqual, 1)
		So(len(ret.Pages), ShouldEqual, 0)
	})

	Convey("测试分页页码计算：数据等于 2 页", t, func() {
		ret := pages(12, 2, 10, 1)
		So(ret.Page, ShouldEqual, 2)
		So(ret.Total, ShouldEqual, 2)
		So(len(ret.Pages), ShouldEqual, 0)
	})

	Convey("测试分页页码计算：数据等于 3 页", t, func() {
		ret := pages(21, 2, 10, 3)
		So(ret.Page, ShouldEqual, 2)
		So(ret.Total, ShouldEqual, 3)
		So(ret.Prev, ShouldEqual, 1)
		So(ret.Next, ShouldEqual, 3)
		So(ret.Pages[0], ShouldEqual, 2)
	})

	Convey("测试分页页码计算：数据大于 3 页，可视数为 3", t, func() {
		ret := pages(31, 2, 10, 3)
		So(ret.Page, ShouldEqual, 2)
		So(ret.Total, ShouldEqual, 4)
		So(ret.Prev, ShouldEqual, 1)
		So(ret.Next, ShouldEqual, 3)
		So(ret.Pages[0], ShouldEqual, 2)

		ret = pages(31, 3, 10, 3)
		So(ret.Page, ShouldEqual, 3)
		So(ret.Total, ShouldEqual, 4)
		So(ret.Prev, ShouldEqual, 2)
		So(ret.Next, ShouldEqual, 4)
		So(ret.Pages[0], ShouldEqual, 3)

		ret = pages(31, 4, 10, 3)
		So(ret.Page, ShouldEqual, 4)
		So(ret.Total, ShouldEqual, 4)
		So(ret.Prev, ShouldEqual, 3)
		So(ret.Next, ShouldEqual, 4)
		So(ret.Pages[0], ShouldEqual, 3)
	})

	Convey("测试分页页码计算：数据等于 4 页，可视数为 4", t, func() {
		ret := pages(33, 2, 10, 4)
		So(ret.Page, ShouldEqual, 2)
		So(ret.Total, ShouldEqual, 4)
		So(ret.Prev, ShouldEqual, 1)
		So(ret.Next, ShouldEqual, 3)
		So(ret.Pages[0], ShouldEqual, 2)
		So(ret.Pages[1], ShouldEqual, 3)
	})

	Convey("测试分页页码计算：数据等于 30 页，可视数为 3", t, func() {
		ret := pages(291, 1, 10, 3)
		So(ret.Page, ShouldEqual, 1)
		So(ret.Total, ShouldEqual, 30)
		So(ret.Prev, ShouldEqual, 1)
		So(ret.Next, ShouldEqual, 2)
		So(ret.Pages[0], ShouldEqual, 2)

		ret = pages(291, 20, 10, 3)
		So(ret.Page, ShouldEqual, 20)
		So(ret.Total, ShouldEqual, 30)
		So(ret.Prev, ShouldEqual, 19)
		So(ret.Next, ShouldEqual, 21)
		So(ret.Pages[0], ShouldEqual, 20)

		ret = pages(291, 30, 10, 3)
		So(ret.Page, ShouldEqual, 30)
		So(ret.Total, ShouldEqual, 30)
		So(ret.Prev, ShouldEqual, 29)
		So(ret.Next, ShouldEqual, 30)
		So(ret.Pages[0], ShouldEqual, 29)

		ret = pages(291, 29, 10, 3)
		So(ret.Page, ShouldEqual, 29)
		So(ret.Total, ShouldEqual, 30)
		So(ret.Prev, ShouldEqual, 28)
		So(ret.Next, ShouldEqual, 30)
		So(ret.Pages[0], ShouldEqual, 29)
	})

	Convey("测试分页页码计算：数据等于 30 页，可视数为 6", t, func() {
		ret := pages(291, 1, 10, 6)
		So(ret.Page, ShouldEqual, 1)
		So(ret.Total, ShouldEqual, 30)
		So(ret.Prev, ShouldEqual, 1)
		So(ret.Next, ShouldEqual, 2)
		So(ret.Pages[0], ShouldEqual, 2)
		So(ret.Pages[1], ShouldEqual, 3)
		So(ret.Pages[2], ShouldEqual, 4)
		So(ret.Pages[3], ShouldEqual, 5)

		ret = pages(291, 30, 10, 6)
		So(ret.Page, ShouldEqual, 30)
		So(ret.Total, ShouldEqual, 30)
		So(ret.Prev, ShouldEqual, 29)
		So(ret.Next, ShouldEqual, 30)
		So(ret.Pages[0], ShouldEqual, 26)
		So(ret.Pages[1], ShouldEqual, 27)
		So(ret.Pages[2], ShouldEqual, 28)
		So(ret.Pages[3], ShouldEqual, 29)
	})
}
