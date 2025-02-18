package localclient

import (
	"context"
	"os/exec"
	"runtime"
)

// 每次下发重新打开需启动shell, 性能差
func RunCMD(ctx context.Context, command string) (string, error) {
	var cmd *exec.Cmd

	// 兼容windows
	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(ctx, "cmd.exe", "/C", command)
	default:
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), err
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
