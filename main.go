package main

import (
	"flag"
	"fmt"
	"indicator/cmd/collect"
	"indicator/cmd/localclient"
	"indicator/cmd/logger"
	"reflect"
	"strings"
)

type MONITOR struct {
	lowperformance   string
	logconsole       string
	logfilelevel     string
	runmode          string
	runtotaltime     string
	needmonitorpname []string
	severtype        string
	monitorinterval  string
}

var GLOBAL_VAR MONITOR

func init() {
	var pnameList string = ""
	flag.StringVar(&GLOBAL_VAR.lowperformance, "lowperformance", "false", "")
	flag.StringVar(&GLOBAL_VAR.logconsole, "logconsole", "warning", "")
	flag.StringVar(&GLOBAL_VAR.logfilelevel, "logfilelevel", "info", "")
	flag.StringVar(&GLOBAL_VAR.runmode, "runmode", "agent", "")
	flag.StringVar(&GLOBAL_VAR.runtotaltime, "runtotaltime", "31536000", "")
	flag.StringVar(&pnameList, "needmonitorpname", "", "")
	flag.StringVar(&GLOBAL_VAR.severtype, "severtype", "", "")
	flag.StringVar(&GLOBAL_VAR.monitorinterval, "monitorinterval", "", "")

	flag.Parse()

	GLOBAL_VAR.needmonitorpname = strings.Split(pnameList, ",")
}

func run_client_mode() {
	logger.LogConsole("runmode - client")

}

func run_controller_mode() {
	logger.LogConsole("runmode - controller")
}

func run_agent_mode() {
	logger.LogConsole("runmode - agent")
	collect.GetCPUInfo()
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
		run_agent_mode()
	case "update":
		run_update_mode()
	case "test":
		// 目前作为upload使用
		run_test_mode("106.15.6.164", 22)
	}
	// logger.LogConsole(GLOBAL_VAR)
}
