package logging

import (
	"fmt"
	"go-gin-blog-api/pkg/file"
	"go-gin-blog-api/pkg/setting"
	"os"
	"time"
)

// 获取日志目录
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

// 获取日志文件名
func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}

// 打开日志文件
func openLogFile(fileName, filePath string) (*os.File, error) {
	// 获取根目录
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	// 获取日志目录路径
	logFilePath := dir + "/" + filePath
	// 获取权限
	perm := file.CheckPermission(logFilePath)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", logFilePath)
	}

	// 获取目录是否存在
	err = file.IsNotExistMkDir(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", logFilePath, err)
	}

	// 打开文件
	handle, err := file.Open(logFilePath+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return handle, err
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
