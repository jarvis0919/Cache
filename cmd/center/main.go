package main

import (
	"cache/config"
	"cache/internal/data"

	// "cache/internal/service/registry"
	"cache/internal/service/registry"
)

func main() {
	//0 生产环境 1 发布环境
	config.Configinit(0)
	//数据库初始化
	data.Sqlinit()
	registry := registry.Registry{}
	registry.Run()
	// time.Sleep(2 * time.Second)
	// node := node.Node{}
	// fmt.Println("<==>")
	// node.Run()
	// select {}
	//启动节点

}
