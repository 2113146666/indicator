package collect

import (
	"indicator/internal/common"
	"indicator/internal/localclient"
	"strings"
	"sync/atomic"
)

var vmstatCMD = "vmstat -a -w -S K 1 2"
var GaugeVmstatData = make(map[string]*atomic.Pointer[string])

func getVMStatInfo() {
	result := localclient.RunCMD(vmstatCMD)
	resSlice := strings.Split(result, "\n")
	if len(resSlice) < 4 {
		return
	}
	title := strings.Fields(resSlice[1])
	data := strings.Fields(resSlice[3])
	for index, ti := range title {
		ti = "vm_" + ti
		value, exists := GaugeVmstatData[ti]
		if !exists || value == nil {
			GaugeVmstatData[ti] = new(atomic.Pointer[string])
		}
		GaugeVmstatData[ti].Store(common.NewStringPtr(data[index]))
	}
}
