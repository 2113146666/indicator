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

var GaugeDiskData = make(map[string]*atomic.Pointer[string])
var getDiskCommand = "df -Thm"

type DiskInfo struct {
	filesystem string
	diskType   string
	blocks     string
	used       string
	available  string
	usePercent string
	mounted    string
}

// 解析df -Thm数据
func getDiskInfo() {
	data := localclient.RunCMD(getDiskCommand)
	dataSlice := strings.Split(data, "\n")
	for _, info := range dataSlice {
		if strings.Contains(strings.ToLower(info), "used") {
			continue
		}

		re := regexp.MustCompile(`\s{2,}`)
		info := re.ReplaceAllString(strings.TrimSpace(info), " ")
		diskInfo := strings.Split(info, " ")
		if len(diskInfo) < 7 {
			continue
		}

		diskname := diskInfo[0] + "#storage"
		_, exists := GaugeDiskData[diskname]
		if !exists {
			GaugeDiskData[diskname] = new(atomic.Pointer[string])
		}

		diskstorage, err := strconv.Atoi(strings.Trim(diskInfo[5], "%"))
		if err != nil {
			logger.LogConsole("failed, reason is use% is not int: %v", diskInfo[5])
		}

		_storageuse := float64(diskstorage) / 100

		GaugeDiskData[diskname].Store(common.NewStringPtr(fmt.Sprintf("%.2f", _storageuse)))
	}
}
