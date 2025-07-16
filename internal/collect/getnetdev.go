package collect

import (
	"fmt"
	"indicator/internal/common"
	"indicator/internal/localclient"
	"indicator/internal/logger"
	"strings"
	"sync/atomic"
)

var GaugeNetDevData = make(map[string]*atomic.Pointer[string])
var getNetDevCMD = "/proc/pid/net/dev"
var DataSlice = []string{"bytes", "packets", "errs", "drop", "fifo", "frame", "compressed", "multicast"}

// 对外开放, 进程监控可以使用
func CalNetDev(pid string) {
	getNetDevCMD = strings.Replace(getNetDevCMD, "pid", pid, -1)
	cmd := fmt.Sprintf("cat %v && echo \"*****\" && sleep 1 && cat %v", getNetDevCMD, getNetDevCMD)
	res := localclient.RunCMD(cmd)

	timesData := strings.Split(res, "*****")
	if len(timesData) < 2 {
		logger.LogConsole(fmt.Sprintf("get /proc/%v/net/dev error", pid))
		return
	}
	logger.LogConsole(timesData[0])
	logger.LogConsole(timesData[1])
	firstData := analyseData(timesData[0])
	secondData := analyseData(timesData[1])
	for head, value := range secondData {
		value1, exists := firstData[head]
		if !exists || value1 == "" {
			continue
		}
		v, exists := GaugeNetDevData[head]
		if v == nil || !exists {
			GaugeNetDevData[head] = new(atomic.Pointer[string])
		}
		GaugeNetDevData[head].Store(common.NewStringPtr(fmt.Sprintf("%v", (common.Int(value)-common.Int(value1))/1000)))
	}
}

func analyseData(data string) map[string]string {
	dataMap := map[string]string{}

	dataSlice := strings.Split(data, "\n")
	for _, info := range dataSlice {
		info = strings.TrimSpace(info)
		if strings.HasPrefix(info, "Inter") || strings.HasPrefix(info, "face") {
			continue
		}
		logger.LogConsole(info)
		infoSlice := strings.Fields(info)
		if len(infoSlice) < 16 {
			continue
		}
		devName := strings.ReplaceAll(infoSlice[0], ":", "")
		for index, num := range infoSlice[1:] {
			prefix := "rx"
			if index > 8 {
				prefix = "tx"
			}
			head := DataSlice[index%8]
			dataMap[devName+"_"+prefix+head] = num
		}
	}
	return dataMap
}

func getNetDev() {
	CalNetDev("")
}
