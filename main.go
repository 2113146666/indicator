package main

import (
	"flag"
	"fmt"
	"indicator/cmd/collect"
	"indicator/cmd/localclient"
	"indicator/cmd/logger"
	"io"
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
	GLOBAL_VAR MONITOR
	lock       sync.Mutex
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
	http.HandleFunc("/metrices", reply_data)
	fmt.Printf("start listen on %v\n", port)
	address := ":" + strconv.Itoa(port)
	http.ListenAndServe(address, nil)
}

func reply_data(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()
	var data = "111"
	io.WriteString(w, data)
}

func run_client_mode() {
	logger.LogConsole("runmode - client")
}

func run_controller_mode() {
	logger.LogConsole("runmode - controller")
}

func run_prometheus_client(interval int) {
	//1、初始化prometheus对象
	//2、起一个服务监听对应端口
	totalTime := time.Duration(GLOBAL_VAR.runtotaltime) * time.Second
	intervalTime := time.Duration(interval) * time.Second

	go start_http_server(9101)
	logger.LogConsole("runmode - prometheus client")
	ticker := time.NewTicker(intervalTime)
	defer ticker.Stop()
	timeout := time.After(totalTime)
loop:
	for {
		select {
		case <-ticker.C:
			collect.GetAllDatas()
		case <-timeout:
			break loop
		}
	}
}

func run_update_mode() {
	logger.LogConsole("runmode - update")
}

func run_test_mode(remote_host string, remote_port int) {
	logger.LogConsole("runmode - test")
	localclient.RunCMD("d: & cd d:/Git_Code/indicator/indicator/ & dir")
	localclient.RunCMD("set GOOS=linux& go build -o indicator ./main.go & set GOOS=windows")
	// localclient.RunCMD(string(fmt.Sprintf(`scp -P %v "D:\Git_Code\indicator\indicator\indicator" root@%v:/root/`, remote_port, remote_host)))
	upload_cmd := fmt.Sprintf("scp -P %v D:/Git_Code/indicator/indicator/indicator root@%v:/root/", remote_port, remote_host)
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
	case "agent":
		run_prometheus_client(GLOBAL_VAR.monitorinterval)
	case "update":
		run_update_mode()
	case "test":
		// 目前作为upload使用
		run_test_mode("106.15.6.164", 22)
	}
	// logger.LogConsole(GLOBAL_VAR)
}
