package collect

import (
	"fmt"
	"indicator/internal/common"
	"indicator/internal/localclient"
	"indicator/internal/logger"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
)

/*
字段	描述	单位
1	主设备号	无
2	次设备号	无
3	设备名称	无
4	成功完成的读请求次数	次数
5	合并的读请求次数	次数
6	读取的扇区数	扇区数（每扇区通常512字节）
7	读操作耗时	毫秒（ms）
8	成功完成的写请求次数	次数
9	合并的写请求次数	次数
10	写入的扇区数	扇区数
11	写操作耗时	毫秒（ms）
12	当前正在进行的I/O操作数	次数
13	I/O操作总耗时（包括读写等）	毫秒（ms）
14	加权I/O操作总耗时（考虑队列深度）	毫秒（ms）
15	成功完成的Discard操作次数	次数
16	合并的Discard请求次数	次数
17	被Discard的扇区数	扇区数
18	Discard操作耗时	毫秒（ms）
19	成功完成的Flush操作次数	次数
20	Flush操作耗时	毫秒（ms）
*/

var GaugeIOData = make(map[string]*atomic.Pointer[string])
var getIOCommand = "cat /proc/diskstats; sleep 1; echo '###split###'; cat /proc/diskstats"

type IOInfo struct {
	majorNum          string
	minorNum          string
	devName           string
	readCount         string
	readMergeCount    string
	readSection       string
	readSpendTime     string
	writeCount        string
	writeMergeCount   string
	writeSection      string
	writeSpendTime    string
	ioRequests        string
	totalTimeSpendIO  string
	weightTimeSpendIO int
}

// 解析/proc/diskstats数据返回一个IOInfo类型的对象
func parseIOData(data string) map[string]IOInfo {
	res := make(map[string]IOInfo)

	dataSlice := strings.Split(data, "\n")
	for _, disk := range dataSlice {
		re := regexp.MustCompile(`\s{2,}`)
		disk = re.ReplaceAllString(strings.TrimSpace(disk), " ")
		diskSlice := strings.Split(disk, " ")
		if len(diskSlice) >= 13 {
			dataObj := IOInfo{}
			weightTimeNum, err := strconv.Atoi(diskSlice[13])
			if err == nil {
				dataObj.weightTimeSpendIO = weightTimeNum
				res[diskSlice[2]] = dataObj
			}
		}
	}
	return res
}

func calculateUtill(first, second map[string]IOInfo) {
	// get io_utill
	for diskname := range second {
		diskUtil := diskname + "_ioutil"
		first_value, exists := first[diskname]
		if exists {
			second_value := second[diskname]

			value, exists := GaugeIOData[diskUtil]
			if !exists || value == nil {
				GaugeIOData[diskUtil] = new(atomic.Pointer[string])
			}
			logger.LogConsole(first_value.weightTimeSpendIO, second_value.weightTimeSpendIO)
			_disk_ioutil := float64(second_value.weightTimeSpendIO-first_value.weightTimeSpendIO) / 1000
			GaugeIOData[diskUtil].Store(common.NewStringPtr(fmt.Sprintf("%.3f", _disk_ioutil)))
		}
	}
}

// 获取IO数据
func getIOInfo() {
	res := localclient.RunCMD(getIOCommand)
	ioSlice := strings.Split(res, "###split###")
	if len(ioSlice) == 2 {
		first_data := parseIOData(ioSlice[0])
		second_data := parseIOData(ioSlice[1])
		calculateUtill(first_data, second_data)
	}
}
