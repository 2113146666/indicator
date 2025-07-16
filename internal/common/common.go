package common

import (
	"strconv"
	"sync/atomic"
)

// 辅助函数：将字符串转为 *string
func NewStringPtr(s string) *string {
	return &s
}

func PtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func PstringToInt(s *atomic.Pointer[string]) int64 {
	var res int64

	if value := s.Load(); value != nil {
		str := *value
		num, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			res = num
		}
	}
	return res
}

func ParseUint(b []byte) uint64 {
	v, _ := strconv.ParseUint(string(b), 10, 64)
	return v
}

func Int(st string) int {
	value, err := strconv.Atoi(st)
	if err != nil {
		return 0
	}
	return value
}
