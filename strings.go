package pkg

import (
	"sort"
	"unicode"
)

// Contains check if element in list, used to test
func Contains(list []string, elem string) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}
	return false
}

// Index returns the index of the first instance of elem in list, or -1 if elem is not present in list.
func Index(list []string, elem string) int {
	for i, v := range list {
		if v == elem {
			return i
		}
	}
	return -1
}

// StringSliceEquals 比较两个字符串数组是否相等
func StringSliceEquals(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// StringSliceSortEquals 比较两个字符串数组排序后是否相等
func StringSliceSortEquals(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}

	sort.Strings(a)
	sort.Strings(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// Map 模拟python的map函数
func Map(fn func(string) string, ss []string) []string {
	var result []string
	for _, field := range ss {
		result = append(result, fn(field))
	}
	return result
}

// RemoveDuplicatesUnordered 字符串列表不保序去重
func RemoveDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

// IsSpaceOr 用作strings.FieldsFunc的分割函数生成器, 效果上相当于用sep分割，同时去掉空格
// strings.FieldsFunc("parking|  def", isSpaceOr('|')) -> [parking def], 而不是[parking   def]
func IsSpaceOr(sep rune) func(c rune) bool {
	return func(c rune) bool {
		if c == sep {
			return true
		}
		return unicode.IsSpace(c)
	}
}

// ToLowerFirst 将第一个字符转为小写
func ToLowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// SplitByLength 按长度分割字符串
func SplitByLength(s string, length int) []string {
	var ss []string
	for i := 0; i < len(s); i = i + length {
		last := i + length
		if i+length >= len(s) {
			last = len(s)
		}
		ss = append(ss, s[i:last])
	}
	return ss
}

// Compare 比较新旧两个列表，返回新增和删除列表
func Compare(olds, news []string) (adds, deletes []string) {
	for _, new := range news {
		if !Contains(olds, new) {
			adds = append(adds, new)
		}
	}
	for _, old := range olds {
		if !Contains(news, old) {
			deletes = append(deletes, old)
		}
	}
	return
}
