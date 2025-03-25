package collect

import (
	"bytes"
	"fmt"
	"indicator/internal/common"
	"indicator/internal/logger"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var GaugeCPUData = make(map[string]*atomic.Pointer[string])
var procStatSlice = []string{"user", "nice", "system", "idle", "iowait", "irq", "softirq", "steal", "guest"}

// 初始化默认值
func init() {
	for _, tit := range procStatSlice {
		GaugeCPUData[tit] = new(atomic.Pointer[string])
		GaugeCPUData[tit].Store(common.NewStringPtr("0.00"))
	}
}

// 获取数据
func getCPUData() map[string]uint64 {
	res := make(map[string]uint64)

	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil
	}

	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte("cpu")) {
			fields := bytes.Fields(line)
			if len(fields) < 11 { // 确保至少有11个字段
				continue
			}

			for index, t := range procStatSlice {
				key := string(fields[0]) + "_" + t
				if strings.HasPrefix(key, "cpu_") {
					key = strings.Replace(key, "cpu_", "cpuall_", 1)
				}
				res[key] = common.ParseUint(fields[index+1])
			}
		}
	}
	return res
}

// 作差
func calculateUsage(first, second map[string]uint64) {
	for key, v1 := range first {
		v2, exists := second[key]
		if !exists {
			_sync := fmt.Sprintf("%s is not exists", key)
			logger.LogConsole(_sync)
		}
		value, exists := GaugeCPUData[key]
		if !exists || value == nil {
			GaugeCPUData[key] = new(atomic.Pointer[string])
		}

		GaugeCPUData[key].Store(common.NewStringPtr(fmt.Sprintf("%.2f", float64(v2-v1))))
	}
}

func getCPUInfo() {
	firstData := getCPUData()
	time.Sleep(1 * time.Second)
	secondData := getCPUData()
	calculateUsage(firstData, secondData)
}
