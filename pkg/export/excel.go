package export

import "go-gin-blog-api/pkg/setting"

// 获取URL
func GetExcelFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}

// 获取相对路径
func GetExcelPath() string {
	return setting.AppSetting.ExportSavePath
}

// 获取绝对路径
func GetExcelFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExcelPath()
}
