package util

// MapMerge 合并字典，如果m1和m2中有同名KEY，使用m2中的值覆盖m1中的值
func MapMerge(m1 map[string]interface{}, mm ...map[string]interface{}) map[string]interface{} {
	if len(mm) == 0 {
		return m1
	}
	for _, mv := range mm {
		for k, v := range mv {
			m1[k] = v
		}
	}
	return m1
}
