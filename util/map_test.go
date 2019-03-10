package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMapMerge1(t *testing.T) {
	Convey("测试Map合并", t, func() {
		m1 := map[string]interface{}{
			"a": "a1",
			"b": "b1",
		}
		m2 := map[string]interface{}{
			"a": "a2",
			"c": "c1",
		}
		m3 := map[string]interface{}{
			"c": "c2",
			"a": "a3",
			"d": "d1",
		}
		mm := MapMerge(m1, nil)
		So(mm, ShouldEqual, m1)

		mm = MapMerge(m1, m2)
		So(mm["a"], ShouldEqual, "a2")

		mm = MapMerge(m1, m2, m3)
		So(mm["a"], ShouldEqual, "a3")
		So(mm["d"], ShouldEqual, "d1")
	})
}
