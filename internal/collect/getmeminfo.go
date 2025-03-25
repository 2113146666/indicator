package collect

import (
	"bytes"
	"indicator/internal/common"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

var GaugeMEMData = map[string]*atomic.Pointer[string]{
	"memtotal":     new(atomic.Pointer[string]),
	"memfree":      new(atomic.Pointer[string]),
	"buffers":      new(atomic.Pointer[string]),
	"cached":       new(atomic.Pointer[string]),
	"used":         new(atomic.Pointer[string]),
	"sreclaimable": new(atomic.Pointer[string]),
}

// 初始化默认值
func init() {
	empty := "0"
	for _, ptr := range GaugeMEMData {
		ptr.Store(&empty)
	}
}

// 获取数据
func getMEMInfo() {

	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}

	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		line := strings.Replace(string(line), " ", "", -1)
		line = strings.ToLower(line)
		line = strings.Replace(line, "kb", "", -1)
		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}
		value, exists := GaugeMEMData[fields[0]]
		if exists {
			value.Store(common.NewStringPtr(fields[1]))
		}
	}

	used := common.PstringToInt(GaugeMEMData["memtotal"])
	used -= common.PstringToInt(GaugeMEMData["memfree"]) + common.PstringToInt(GaugeMEMData["buffers"]) + common.PstringToInt(GaugeMEMData["cached"]) + common.PstringToInt(GaugeMEMData["sreclaimable"])
	GaugeMEMData["used"].Store(common.NewStringPtr(strconv.FormatInt(used, 10)))
}
