package collect

import (
	"fmt"
	"indicator/internal/common"
	"sync/atomic"
)

var GaugeAllData atomic.Pointer[string]

func gaugeData() {

	output := "# HELP cpu_percent use CPU%\n# TYPE cpu_percent gauge\n"
	for mode, value := range GaugeCPUData {
		output += fmt.Sprintf("cpu_percent{mode=\"%s\"} %v\n", mode, common.PtrToString(value.Load()))
	}

	output += "# HELP Vmstat times\n# TYPE Vmstat gauge\n"
	for mode, value := range GaugeVmstatData {
		output += fmt.Sprintf("vmstat{mode=\"%s\"} %v\n", mode, common.PtrToString(value.Load()))
	}

	output += "# HELP MEM avail MB\n# TYPE mem_info gauge\n"
	for mode, value := range GaugeMEMData {
		output += fmt.Sprintf("mem_mb{mode=\"%s\"} %v\n", mode, common.PtrToString(value.Load()))
	}

	output += "# HELP IO use %\n# TYPE io_percent gauge\n"
	for mode, value := range GaugeIOData {
		output += fmt.Sprintf("io_percent{mode=\"%s\"} %v\n", mode, common.PtrToString(value.Load()))
	}

	output += "# HELP Disk Storage use%\n# TYPE disk_percent gauge\n"
	for mode, value := range GaugeDiskData {
		output += fmt.Sprintf("disk_percent{mode=\"%s\"} %v\n", mode, common.PtrToString(value.Load()))
	}

	GaugeAllData.Store(common.NewStringPtr(output))
}

// 数据采集入口
func GetAllDatas() {
	getCPUInfo()
	getVMStatData()
	getMEMInfo()
	getIOInfo()
	getDiskInfo()
	// 数据拼接
	gaugeData()
}
