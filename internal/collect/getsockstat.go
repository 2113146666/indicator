package collect

import (
	"indicator/internal/common"
	"indicator/internal/logger"
	"os"
	"strings"
	"sync/atomic"
)

var GaugeSockStatData = make(map[string]*atomic.Pointer[string])
var SockStat = "/proc/net/sockstat"
var Sock6Stat = "/proc/net/sockstat6"

func getSockStatInfo() {
	data, err := os.ReadFile(SockStat)
	if err != nil {
		return
	}

	data6, err := os.ReadFile(Sock6Stat)
	if err != nil {
		return
	}

	allDataSlice := strings.Split(string(data)+"\n"+string(data6), "\n")

	logger.LogConsole(allDataSlice)
	for _, data := range allDataSlice {
		dataslice := strings.Fields(data)
		if len(dataslice) < 3 {
			continue
		}
		logger.LogConsole(dataslice)
		head := strings.Replace(dataslice[0], ":", "", -1)
		for i := 1; i < len(dataslice[1:]); i += 2 {
			title := head + dataslice[i]
			value, exists := GaugeSockStatData[title]
			if !exists || value == nil {
				GaugeSockStatData[title] = new(atomic.Pointer[string])
			}
			GaugeSockStatData[title].Store(common.NewStringPtr(dataslice[i+1]))
		}
	}
}
