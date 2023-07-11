package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

// APP 配置相关
type App struct {
	JwtSecret       string
	PageSize        int
	RuntimeRootPath string

	ImagePrefixUrl string
	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string

	ExportSavePath string
	PrefixUrl      string
}

// 服务器配置相关
type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// database 数据库配置相关
type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

// Redis 相关配置
type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

// 全局变量
var (
	AppSetting      = &App{}
	ServerSetting   = &Server{}
	DatabaseSetting = &Database{}
	RedisSetting    = &Redis{}
)

// 相同包下的 init 函数：按照源文件编译顺序决定执行顺序
// 不同包下的 init 函数：按照包导入的依赖关系决定先后顺序
// 如果希望执行顺序按照自己的要求来，就自己封装对应的方法进行统一调用，映射全局变量

// 映射全局配置整体加载
func Setup() {
	var err error
	Cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	// [APP] 相关配置
	err = Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo AppSetting err: %v", err)
	}
	// 上传图片大小
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024

	// [Server] 相关配置
	err = Cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo ServerSetting err: %v", err)
	}

	// 服务器超时时间
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	// [Database] 相关配置
	err = Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo DatabaseSetting err: %v", err)
	}

	// [redis] 相关配置
	err = Cfg.Section("redis").MapTo(RedisSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}

}
