package util

import (
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSliceInt1(t *testing.T) {
	Convey("测试相等的切片", t, func() {
		So(SliceIntEqual([]int{1, 2, 3}, []int{1, 2, 3}), ShouldBeTrue)
	})
	Convey("测试元素相同但顺序不同的切片", t, func() {
		So(SliceIntEqual([]int{2, 3, 1}, []int{1, 3, 2}), ShouldBeTrue)
	})
	Convey("测试不等长切片", t, func() {
		So(SliceIntEqual([]int{1, 2, 3}, []int{1, 3}), ShouldBeFalse)
	})
	Convey("测试等长不相等切片", t, func() {
		So(SliceIntEqual([]int{1, 2, 3}, []int{1, 4, 3}), ShouldBeFalse)
	})
}

func TestSliceIntToString1(t *testing.T) {
	Convey("测试数字切片转换为字符切片", t, func() {
		So(strings.Join(SliceIntToString([]int{1, 2, 3}), ","), ShouldEqual, "1,2,3")
		So(strings.Join(SliceIntToString([]int{12, 23, 34}), ","), ShouldEqual, "12,23,34")
	})
}

func TestSliceIntRand1(t *testing.T) {
	Convey("打乱顺序返回切片元素测试", t, func() {
		a := []int{1, 2, 3, 4, 5}
		b := SliceIntRand(a)
		fmt.Println(b)
		So(b, ShouldNotBeNil)
	})
}

func TestInSlice1(t *testing.T) {
	Convey("判断一个元素是否在切片中出现", t, func() {
		a := []int{1, 2, 3, 4, 5}
		b := []string{"aa", "bb", "cc"}
		t1 := time.Now()
		c := []time.Time{time.Now(), t1}
		So(InSlice(a, 1, "int"), ShouldBeTrue)
		So(InSlice(a, 2, "int"), ShouldBeTrue)
		So(InSlice(a, 6, "int"), ShouldBeFalse)
		So(InSlice(b, "", "string"), ShouldBeFalse)
		So(InSlice(b, "a", "string"), ShouldBeFalse)
		So(InSlice(b, "aa", "string"), ShouldBeTrue)
		So(InSlice(b, "bb", "string"), ShouldBeTrue)
		So(InSlice(c, t1, "time"), ShouldBeTrue)
	})
}

func TestSliceIntDiff1(t *testing.T) {
	Convey("获取数字切片差集", t, func() {
		a := []int{1, 2, 3}
		b := []int{1, 3, 4}
		c := fmt.Sprintf("%v", SliceIntDiff(a, b))
		So(c, ShouldEqual, "[2]")
	})
}

func TestSliceStringDiff1(t *testing.T) {
	Convey("获取字符串切片差集", t, func() {
		a := []string{"1", "2", "3"}
		b := []string{"1", "3", "4"}
		c := fmt.Sprintf("%v", SliceStringDiff(a, b))
		So(c, ShouldEqual, "[2]")
	})
}
