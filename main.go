// 主包，负责功能模块初始化
package main

import (
	"runtime"

	"github.com/denghongcai/MessageHive/modules/command"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	command.Execute()
}
