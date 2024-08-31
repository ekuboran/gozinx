package utils

import (
	"log"
	"os"
)

// Debugging
const Debug = true

func init() {
	// 打开日志文件
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}

	// 将日志输出重定向到文件
	log.SetOutput(file)
}

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

func DPrintln(a ...interface{}) (n int, err error) {
	if Debug {
		log.Println(a...)
	}
	return
}
