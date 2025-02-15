package main

import (
	"fmt"
	"indicator/cmd/logger"
)

func main() {
	fmt.Printf("%v", "hello, world!\n")
	logger.LogConsole("logger print")
}
