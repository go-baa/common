// Package excel 提供了一个excel读写的封装
// 读取时需要注意，循环Exce.Fetch()直到返回nil就没有数据了，但是如果表格中有合并等情况，返回的row中并不能反应出来，需要自己check row中的数据
package excel

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/go-baa/log"
	"github.com/tealeg/xlsx"
)

const (
	modeRead uint = iota
	modeRwrite
)

// MapValue Excel读取出来的每一行的值
type MapValue map[string]interface{}

// HeaderColumn 表头字段定义
type HeaderColumn struct {
	Field   string // 字段，数据映射到的数据字段名
	Title   string // 标题，表格中的列名称
	Require bool   // 是否必选
}

// Get 获取一个interface{}的值
func (t MapValue) Get(key string) interface{} {
	return t[key]
}

// GetString 获取一个string类型的值
func (t MapValue) GetString(key string) string {
	var s string
	if v, ok := t[key]; ok {
		switch v.(type) {
		case string:
			s = v.(string)
		case fmt.Stringer:
			s = v.(fmt.Stringer).String()
		default:
			s = fmt.Sprintf("%v", v)
		}
	}
	return strings.TrimSpace(s)
}

// GetInt 获取一个int类型的值
func (t MapValue) GetInt(key string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(t.GetString(key)))
	return i
}

// Excel 操作Excel的结构
type Excel struct {
	header         []*HeaderColumn
	fields         []string
	file           *xlsx.File
	sheet          *xlsx.Sheet
	ignoreOverCell bool // 是否忽略超出的列
	cellNum        int  // 表格的列数
	rowi           int  // 当前读到了第几行
	mode           uint // 是否是写模式
}

// NewWriter 创建一个Excel写操作实例
func NewWriter() (*Excel, error) {
	t := new(Excel)
	t.mode = modeRwrite
	t.file = xlsx.NewFile()
	var err error
	t.sheet, err = t.file.AddSheet("Sheet1")
	if err != nil {
		return nil, err
	}
	return t, nil
}

// NewReader 创建一个Excel读操作实例
// 接受的 header 为表头和字段映射关系，如：{Field: "title", Title: "标题", Requre: true}
// headerlineNum 表头行号是第几行，默认是第1行
// igoreNewLine 如果表头有换行，是否删除掉换行以后的字符
func NewReader(file string, header []*HeaderColumn, headerlineNum int, igoreNewLine bool) (*Excel, error) {
	if len(header) == 0 {
		return nil, errors.New("Excel.NewReader 错误: 表头不能为空")
	}

	xlFile, err := xlsx.OpenFile(file)
	if err != nil {
		return nil, errors.New("Excel.NewReader 错误: " + err.Error())
	}

	t := new(Excel)
	t.header = header
	t.cellNum = len(header)
	t.rowi = headerlineNum // 从第几行开始读取内容
	t.mode = modeRead
	t.file = xlFile

	// 检测是否为空文件
	if len(t.file.Sheets) == 0 {
		return nil, errors.New("Excel.NewReader 错误: 文件中没有表格")
	}
	// 检测表格中是否有数据
	if len(t.file.Sheets[0].Rows) == 0 {
		return nil, errors.New("Excel.NewReader 错误: 表格中没有数据")
	}
	// 检测表头
	err = t.checkHeader(igoreNewLine)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// checkHeader 读模式下，检测表头是否和要求的字段一致
// fieldsMap 初始化所有的字段，默认为false,后面发现一个即设为true,到最后如果还有false,那么字段不匹配
// igoreNewLine 是否将表头中的换行下面的内容忽略
func (t *Excel) checkHeader(igoreNewLine bool) error {
	header := make(map[string]string)
	fieldsMap := make(map[string]*HeaderColumn)
	for _, column := range t.header {
		fieldsMap[column.Field] = column
		header[column.Title] = column.Field
	}

	// 检查表头列是否和指定的表头对应
	if t.rowi == 0 {
		t.rowi = 1
	}
	headerRow := t.file.Sheets[0].Rows[t.rowi-1]

	// 获取所有表头对应的字段，并加入字段列表，同时验证字段是否存在
	var v string
	var ok bool
	headerText := make([]string, 0)
	for i, cell := range headerRow.Cells {
		cellText := cell.String()
		cellText = strings.TrimSpace(cellText)
		if igoreNewLine {
			if pos := strings.Index(cellText, "\n"); pos > 0 {
				cellText = cellText[:pos]
			}
		}
		cellText = strings.Replace(cellText, "\n", "", -1)
		cellText = strings.TrimSpace(cellText)
		headerText = append(headerText, cellText)
		v, ok = header[cellText]
		if cellText == "" || ok == false {
			v = strconv.Itoa(i)
			t.fields = append(t.fields, v)
			continue
		}
		t.fields = append(t.fields, v)
		fieldsMap[v] = nil
	}

	// 验证是否所有的字段，都读到了
	for _, column := range fieldsMap {
		if column != nil && column.Require {
			fmt.Printf("Excel.CheckHeader 错误：所有的表头如下：\n%v\n", headerText)
			return errors.New("Excel.CheckHeader 错误: 表头中的列 " + column.Title + " 在文件中不存在")
		}
	}

	header = nil
	fieldsMap = nil

	return nil
}

// Fetch 逐行读取内容，请循环调用该值，直到返回 nil
// 该方法会直接跳过空行
func (t *Excel) Fetch() (MapValue, error) {
	var err error
	if t.mode != modeRead {
		return nil, fmt.Errorf("Excel.Fetch 读取错误: 当前表格不允许读取")
	}

	if t.rowi == 0 {
		t.rowi = 1
	}

	// 没有更多的数据了
	if t.rowi >= len(t.file.Sheets[0].Rows) {
		return nil, nil
	}

	blank := true
	row := make(MapValue)
	for i, cell := range t.file.Sheets[0].Rows[t.rowi].Cells {
		v := cell.String()
		v = strings.TrimSpace(v)
		if blank && len(v) > 0 {
			blank = false
		}
		if i < len(t.fields) {
			row[t.fields[i]] = v
		} else {
			// 检测到行的列超出了字段设置，但是读取者并没有强制限制，报错，但忽略错误
			log.Errorf("Excel.Fetch 读取错误: 行 %d 列 %d 未在表头中指定", t.rowi, i)
		}
	}

	t.rowi++

	// 如果有空行直接跳过
	if blank {
		row, err = t.Fetch()
	}

	return row, err
}

// Reset 重置游标
func (t *Excel) Reset() {
	if t.mode != modeRead {
		return
	}

	t.rowi = 0
}

// SetHeader 写模式下，设置字段表头和字段顺序
// 参数 header 为表头和字段映射关系，如：HeaderColumn{Field:"title", Title: "标题", Requre: true}
// 参数 width  为表头每列的宽度，单位 CM：map[string]float64{"title": 0.8}
func (t *Excel) SetHeader(header []*HeaderColumn, width map[string]float64) error {
	if t.mode != modeRwrite {
		return errors.New("Excel.SetHeader 错误: 当前为读模式，不支持写入操作")
	}

	if len(header) == 0 {
		return errors.New("Excel.SetHeader 错误: 表头不能为空")
	}

	// 表头样式
	font := xlsx.DefaultFont()
	font.Bold = true

	alignment := xlsx.DefaultAlignment()
	alignment.Vertical = "center"

	style := xlsx.NewStyle()
	style.Font = *font
	style.Alignment = *alignment

	style.ApplyFont = true
	style.ApplyAlignment = true

	// 设置表头字段
	t.header = header
	row := t.sheet.AddRow()
	row.SetHeightCM(1.0)
	for _, column := range header {
		t.fields = append(t.fields, column.Field)
		cell := row.AddCell()
		cell.Value = column.Title
		cell.SetStyle(style)
	}

	// 表格列，宽度
	if len(t.fields) > 0 {
		for k, v := range t.fields {
			if width[v] > 0.0 {
				t.sheet.SetColWidth(k, k, width[v]*10)
			}
		}
	}

	return nil
}

// Write 写入行
func (t *Excel) Write(data MapValue) error {
	if t.mode != modeRwrite {
		return errors.New("Excel.Write 错误: 当前为读模式，不支持写入操作")
	}

	if len(data) == 0 {
		return nil
	}

	row := t.sheet.AddRow()
	row.SetHeightCM(0.8)
	for _, field := range t.fields {
		row.AddCell().Value = data.GetString(field)
	}

	return nil
}

// Save 输出到io.Writer
func (t *Excel) Save(w io.Writer) error {
	if t.mode != modeRwrite {
		return errors.New("Excel.Save 错误: 当前为读模式，不支持写入操作")
	}
	return t.file.Write(w)
}

// SaveToFile 直接输出到文件
func (t *Excel) SaveToFile(path string) error {
	if t.mode != modeRwrite {
		return errors.New("Excel.SaveToFile 错误: 当前为读模式，不支持写入操作")
	}
	return t.file.Save(path)
}
