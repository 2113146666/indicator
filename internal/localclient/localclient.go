package localclient

import (
	"context"
	"indicator/internal/logger"
	"os/exec"
	"reflect"
	"runtime"
	"time"
)

var cmdTimeout = 60 * time.Second

// 每次下发重新打开需启动shell, 性能差
func runCmd(ctx context.Context, command string) (string, error) {
	var cmd *exec.Cmd

	// 兼容windows
	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(ctx, "cmd.exe", "/C", command)
	default:
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}

	output, err := cmd.CombinedOutput()

	return string(output), err
}

func RunCMD(command string) string {
	var ctx, cancel = context.WithTimeout(context.Background(), cmdTimeout*time.Second)
	defer cancel()

	res, err := runCmd(ctx, command)
	if err != nil {
		logger.LogConsole(err)
	}
	logger.LogConsole("\nrun cmd ==> [", command, "]\nresult: ", res, reflect.TypeOf(res))
	// fmt.Printf("\n[string]%#v\n", string(res))

	return res
}

// // shell长连接方案

// func NewSessionPool() {

// }

// var defaultPool = NewSessionPool(runtime.NumCPU())

// func RunCMD(ctx context.Context, command string) (string, error) {
// 	// 拿到会话对象
// 	session := defaultPool.Get()
// 	// 保护措施
// 	defer defaultPool.Put(session)
// }
