package main

import (
	"myapp/internal/cmd"
	_ "myapp/packed"

	"github.com/gogf/gf/v2/os/gctx"
)

// 入口调用gf cmd
func main() {
	cmd.Main.Run(gctx.New())
}
