package main

import (
	"flag"
	"fmt"
	"indicator/internal/collect"
	"indicator/internal/localclient"
	"indicator/internal/logger"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MONITOR struct {
	lowperformance   string
	logconsole       string
	logfilelevel     string
	runmode          string
	runtotaltime     int
	needmonitorpname []string
	severtype        string
	monitorinterval  int
}

var (
	GLOBAL_VAR  MONITOR
	lock        sync.Mutex
	listen_fail bool
)

func init() {
	var pnameList string = ""
	flag.StringVar(&GLOBAL_VAR.lowperformance, "lowperformance", "false", "")
	flag.StringVar(&GLOBAL_VAR.logconsole, "logconsole", "warning", "")
	flag.StringVar(&GLOBAL_VAR.logfilelevel, "logfilelevel", "info", "")
	flag.StringVar(&GLOBAL_VAR.runmode, "runmode", "agent", "")
	flag.IntVar(&GLOBAL_VAR.runtotaltime, "runtotaltime", 31536000, "")
	flag.StringVar(&pnameList, "needmonitorpname", "", "")
	flag.StringVar(&GLOBAL_VAR.severtype, "severtype", "", "")
	flag.IntVar(&GLOBAL_VAR.monitorinterval, "monitorinterval", 15, "")

	flag.Parse()

	GLOBAL_VAR.needmonitorpname = strings.Split(pnameList, ",")
}

func start_http_server(port int) {

	maxRetry := 99
	address := ":" + strconv.Itoa(port)

	http.HandleFunc("/metrics", metricsHandler)

loop:
	for i := maxRetry; i > 0; i-- {
		logger.LogConsole("start server on port %v...", port)
		err := http.ListenAndServe(address, nil)

		if err != nil {
			fmt.Printf("listen failed, reason is %s", err)
			time.Sleep(1 * time.Second)
			continue
		}

		break loop
	}

	listen_fail = true
	logger.LogConsole("The startup service listening on port 9101 has failed 99 times, and the main program will exit")

}

func run_client_mode() {
	logger.LogConsole("runmode - client")
}

func run_controller_mode() {
	logger.LogConsole("runmode - controller")
}

// 逻辑: 起线程监听9101端口并响应, 主程序循环采集数据
func run_prometheus_client(totalTime int, interval int, port int) {

	// 监听并响应采集数据
	go start_http_server(port)

	startTime := time.Now().UnixNano()
	futureEndTime := startTime + int64(totalTime)*1e9

loop:
	for {

		// 记录开始时间
		_start := time.Now().UnixNano()

		collect.GetAllDatas()

		// 记录结束时间
		_end := time.Now().UnixNano()

		// 记录时间差
		_diff := _end - _start

		_seconds := float64(_diff) / 1e9

		_need_reload_log := fmt.Sprintf("gauge data spend %v nano second", _seconds)

		logger.LogConsole(_need_reload_log)

		_spendTime := float64(interval) - _seconds
		if _spendTime > 0 {
			time.Sleep(time.Duration(_spendTime * float64(time.Second)))
		}

		if _end > futureEndTime || listen_fail {
			break loop
		}
	}

}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	output := "# HELP cpu_percent CPU 使用百分比\n# TYPE cpu_percent gauge\n"
	for mode, value := range collect.GaugeCPUData {
		output += fmt.Sprintf("cpu_percent{mode=\"%s\"} %s\n", mode, value)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(output))
}

// 调试逻辑
func run_test_mode() {
	logger.LogConsole("runmode - update")
}

// 将项目编译并上传到远端
func run_upload_mode(remote_host string, remote_port int) {
	logger.LogConsole("runmode - test")
	localclient.RunCMD("d: & cd d:/Git_Code/indicator/indicator/ & dir")
	localclient.RunCMD("set GOOS=linux& go build -o indicator ./main.go & set GOOS=windows")
	// localclient.RunCMD(string(fmt.Sprintf(`scp -P %v "D:\Git_Code\indicator\indicator\indicator" root@%v:/root/`, remote_port, remote_host)))
	upload_cmd := fmt.Sprintf("scp -P %v D:/golang/project/indicator/indicator root@%v:/root/", remote_port, remote_host)
	logger.LogConsole(upload_cmd, reflect.TypeOf(upload_cmd))
	localclient.RunCMD(upload_cmd)
	logger.LogConsole("upload end")
}

func main() {

	switch GLOBAL_VAR.runmode {

	case "client":
		run_client_mode()

	case "controller":
		run_controller_mode()

	// 作为prometheus client运行
	case "agent":
		run_prometheus_client(GLOBAL_VAR.runtotaltime, GLOBAL_VAR.monitorinterval, 9101)

	// 编译并上传代码到指定服务器
	case "upload":
		// 目前作为upload使用
		run_upload_mode("106.15.6.164", 22)

	// 测试代码
	case "test":
		run_test_mode()
	}
}
