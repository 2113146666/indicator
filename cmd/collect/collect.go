package collect

import (
	"indicator/cmd/localclient"
	"indicator/cmd/logger"
)

var Gauge_data = make(map[string]string)

func getCPUInfo() string {
	res := localclient.RunCMD("cat /proc/cpuinfo")
	logger.LogConsole(res)
	return res
}

func GetAllDatas() {
	Gauge_data["cpu"] = getCPUInfo()
}
