package main

import (
	"log"
	"oliujunk/input/soilmoisture"
	"oliujunk/input/tiannankeji"
	"oliujunk/input/tulian"
	"oliujunk/input/xph_527122894"
)

func init() {
	// 日志信息添加文件名行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {

	go soilmoisture.Start()

	go tiannankeji.Start()

	go tulian.Start()

	go xph_527122894.Start()

	select {}
}
