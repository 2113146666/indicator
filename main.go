package main

import (
	"flag"
	"indicator/cmd/logger"
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

func main() {
	logger.LogConsole(GLOBAL_VAR)
}
