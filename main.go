package main

import (
	"context"
	"flag"
	"indicator/cmd/localclient"
	"indicator/cmd/logger"
	"strings"
	"time"
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
}

func run_update_mode() {
	logger.LogConsole("runmode - update")
}

func run_test_mode() {
	logger.LogConsole("runmode - test")
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	std, err := localclient.RunCMD(ctx, "dir")
	logger.LogConsole(std, err)

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
		run_test_mode()
	}
	logger.LogConsole(GLOBAL_VAR)
}
