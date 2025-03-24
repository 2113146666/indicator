package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var logFile *os.File

func init() {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	dir := "./output"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return
		}
	}

	filePath := filepath.Join(currentPath, "./output/output.log")
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)
	log.SetPrefix("[")
	log.SetFlags(log.LstdFlags | log.LUTC)
}

// 自定义式日志打印
func LogConsole(strings ...interface{}) {
	fileName, codeline := getLineNumber(2)
	console := ""
	for _, inter := range strings {
		console += fmt.Sprintf("%v ", inter)
	}
	log.Printf("%v][%d]: %s\n", fileName, codeline, console)
}

// 返回调用该函数对应栈帧的行号
func getLineNumber(skip int) (fileName string, num int) {
	if skip == 0 {
		skip = 1
	}
	_, file, line, _ := runtime.Caller(skip)
	fileName = filepath.Base(file)
	return fileName, line
}

func CloseLogFile() {
	logFile.Close()
}
