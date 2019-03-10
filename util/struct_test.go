package util

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type User struct {
	Name string
	ID   int `map:"userid"`
	Age  int8
	Date time.Time
}

type Student struct {
	User
	Class string
}

func TestMapFillStruct1(t *testing.T) {
	Convey("测试map填充struct", t, func() {
		Convey("标准Map和结构体字段对应", func(){
			data := make(map[string]interface{})
			data["Name"] = "张三"
			data["Age"] = 26
			data["ID"] = 1
			data["Date"] = "2015-09-29 00:00:00"

			result := &Student{}
			err := MapFillStruct(data, result)
			So(err, ShouldBeNil)
			So(result.ID, ShouldEqual, 1)
		})
		Convey("驼峰命名转换规则", func(){
			data := make(map[string]interface{})
			data["name"] = "张三"
			data["age"] = 26
			data["id"] = 1
			data["Date"] = "2015-09-29 00:00:00"

			result := &Student{}
			err := MapFillStruct(data, result)
			So(err, ShouldBeNil)
			So(result.ID, ShouldEqual, 1)
		})
		Convey("通过mapTag命名规则转换", func(){
			data := make(map[string]interface{})
			data["name"] = "张三"
			data["age"] = 26
			data["userid"] = 2
			data["Date"] = "2015-09-29 00:00:00"

			result := &Student{}
			err := MapFillStruct(data, result)
			So(err, ShouldBeNil)
			So(result.ID, ShouldEqual, 2)
		})
	})
}

func TestStructToMap1(t *testing.T) {
	Convey("测试struct转换Map", t, func() {
		data := &User{}
		data.Name = "张三"
		data.Age = 26
		data.Date, _ = time.Parse("2006-01-02", "2015-09-29")

		result := make(map[string]interface{})
		StructToMap(data, result)
		So(result, ShouldNotBeNil)
		fmt.Println(result)

		data = &User{}
		data.Name = "张三"
		data.Age = 26
		result = make(map[string]interface{})
		StructToMap(data, result)
		So(result, ShouldNotBeNil)
		fmt.Println(result)
	})
}

func TestStructToStruct1(t *testing.T) {
	Convey("测试结构体复制", t, func() {
		d1 := &User{
			Name: "张三",
			Age:  26,
			Date: time.Now(),
		}
		d2 := new(User)
		StructToStruct(d1, d2)
		So(d2.Age, ShouldEqual, 26)
	})
}

func TestSortByID(t *testing.T) {
	Convey("测试结构体slice排序", t, func() {
		d1 := &User{
			Name: "张大",
			Age:  26,
			Date: time.Now(),
		}
		d2 := &User{
			Name: "张二",
			Age:  25,
			Date: time.Now(),
		}
		d3 := &User{
			Name: "张三",
			Age:  24,
			Date: time.Now(),
		}
		var dSlice = []*User{d1, d2, d3}
		sortIDs := []int{24, 26, 25}
		sorted := SortByID(sortIDs, dSlice, "Age").([]*User)
		So(sorted[0].Age, ShouldEqual, 24)
		So(sorted[1].Age, ShouldEqual, 26)
		So(sorted[2].Age, ShouldEqual, 25)
	})
}
