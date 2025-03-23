package collect

import (
	"bytes"
	"fmt"
	"indicator/cmd/logger"
	"os"
	"strconv"
	"strings"
	"time"
)

var GaugeCPUData = make(map[string]string)
var procStatSlice = []string{"user", "nice", "system", "idle", "iowait", "irq", "softirq", "steal", "guest"}

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
				res[key] = parseUint(fields[index+1])
			}
		}
	}
	return res
}

func parseUint(b []byte) uint64 {
	v, _ := strconv.ParseUint(string(b), 10, 64)
	return v
}

func calculateUsage(first, second map[string]uint64) {
	for key, v1 := range first {
		v2, exists := second[key]
		if !exists {
			_sync := fmt.Sprintf("%s is not exists", key)
			logger.LogConsole(_sync)
		}

		GaugeCPUData[key] = fmt.Sprintf("%.2f", float64(v2-v1))

	}
}

func getCPUInfo() {
	firstData := getCPUData()
	time.Sleep(1 * time.Second)
	secondData := getCPUData()
	calculateUsage(firstData, secondData)
}
