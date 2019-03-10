package util

import "github.com/safeie/pdf"

// PdfPageCount 读取PDF文件的页码数
func PdfPageCount(file string) int {
	r, err := pdf.Open(file)
	if err != nil {
		return 0
	}
	return r.NumPage()
}
