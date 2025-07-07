package collect

import (
	"bytes"
	"indicator/internal/common"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

var MEMInfoCMD = "/proc/meminfo"

var GaugeMEMData = map[string]*atomic.Pointer[string]{
	"memtotal":         new(atomic.Pointer[string]),
	"memfree":          new(atomic.Pointer[string]),
	"memavailable":     new(atomic.Pointer[string]),
	"mem_buffers":      new(atomic.Pointer[string]),
	"mem_cached":       new(atomic.Pointer[string]),
	"mem_active":       new(atomic.Pointer[string]),
	"mem_inactive":     new(atomic.Pointer[string]),
	"mem_swaptotal":    new(atomic.Pointer[string]),
	"mem_swapcached":   new(atomic.Pointer[string]),
	"mem_swapfree":     new(atomic.Pointer[string]),
	"mem_shmem":        new(atomic.Pointer[string]),
	"memused":          new(atomic.Pointer[string]),
	"memused(%)":       new(atomic.Pointer[string]),
	"mem_sreclaimable": new(atomic.Pointer[string]),
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

	data, err := os.ReadFile(MEMInfoCMD)
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
		if value, exists := GaugeMEMData[fields[0]]; exists {
			mb, err := strconv.Atoi(fields[1])
			if err != nil {
				continue
			}
			value.Store(common.NewStringPtr(strconv.Itoa(mb / 1024)))
			continue
		}
		_tit := "mem_" + fields[0]
		if value, exists := GaugeMEMData[_tit]; exists {
			mb, err := strconv.Atoi(fields[1])
			if err != nil {
				continue
			}
			value.Store(common.NewStringPtr(strconv.Itoa(mb / 1024)))
		}
	}

	used := common.PstringToInt(GaugeMEMData["memtotal"])
	used -= common.PstringToInt(GaugeMEMData["memfree"]) + common.PstringToInt(GaugeMEMData["mem_buffers"]) + common.PstringToInt(GaugeMEMData["mem_cached"]) + common.PstringToInt(GaugeMEMData["mem_sreclaimable"])
	GaugeMEMData["memused"].Store(common.NewStringPtr(strconv.FormatInt(used, 10)))
	GaugeMEMData["memused(%)"].Store(common.NewStringPtr(strconv.FormatInt(used*100/common.PstringToInt(GaugeMEMData["memtotal"]), 10)))
}
