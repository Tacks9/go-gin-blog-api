package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

// 全局变量声明
var (
	// 配置
	Cfg *ini.File

	// 运行模式
	RunMode string

	// API服务端口
	HTTPPort int

	// 超时时间
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// 分页配置
	PageSize int

	// JWT 密钥
	JwtSecret string
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	// APP 环境配置
	LoadBase()

	// HTTP Server
	LoadServer()

	// APP 业务配置
	LoadApp()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	JwtSecret = sec.Key("JWT_SECRET").MustString("@A#B$C%D^E&F*G!1234567")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}
