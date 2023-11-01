package main

import (
	"cache/config"
	"cache/internal/data"
	"cache/internal/service/node"
)

func main() {
	//0 生产环境 1 发布环境
	config.Configinit(0)
	//数据库初始化
	data.Sqlinit()
	node := node.Node{}
	node.Run()
	select {}
	//启动节点

}
