package collect

import (
	"indicator/cmd/localclient"
	"indicator/cmd/logger"
)

func GetCPUInfo() {
	res := localclient.RunCMD("cat /proc/cpuinfo")
	logger.LogConsole(res)
}
