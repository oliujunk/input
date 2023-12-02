package main

import (
	"log"
	"oliujunk/input/rn_2184568_history"
)

func init() {
	// 日志信息添加文件名行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {

	//go soilmoisture.Start()
	//
	//go tulian.Start()
	//
	//go tiannankeji.Start()
	//
	//go rn_2184568.Start()
	//
	//go xph_527046329.Start()
	//
	//go xph_527122894.Start()

	go rn_2184568_history.Start()

	select {}
}
