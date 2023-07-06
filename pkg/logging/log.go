package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

var (
	FHandle *os.File

	// 默认日志前缀
	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	// 日志记录器
	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

// 创建 枚举值
const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// 日志启动  Setup initialize the log instance
func Setup() {
	var err error
	filePath := getLogFilePath()
	fileName := getLogFileName()
	FHandle, err = openLogFile(fileName, filePath)
	if err != nil {
		log.Fatalf("logging.Setup err: %v", err)
	}
	// LstdFlags 日志记录的属性，日期-时间

	// 实例化日志记录器
	logger = log.New(FHandle, DefaultPrefix, log.LstdFlags)
}

// 设置日志行前缀
func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		// 日志前缀
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v)
}
