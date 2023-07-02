package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = "runtime/logs/"
	LogSaveName = "log"
	LogFileExt  = "log"
	TimeFormat  = "20060102"
)

// 获取日志目录
func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

// 获取日志文件位置
func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func openLogFile(filepath string) *os.File {
	//获取文件描述信息
	_, err := os.Stat(filepath)
	switch {
	// 是否存在
	case os.IsNotExist(err):
		mkDir()
	// 是否有权限
	case os.IsPermission(err):
		log.Fatalf("Permission:%v", err)
	}

	// 打开文件
	handle, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile:%v", err)
	}

	return handle
}

// 创建目录
func mkDir() {
	// 获取当前目录
	dir, _ := os.Getwd()
	// 创建子目录
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
