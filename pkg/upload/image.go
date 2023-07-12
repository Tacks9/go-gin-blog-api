package upload

import (
	"fmt"
	"go-gin-blog-api/pkg/file"
	"go-gin-blog-api/pkg/logging"
	"go-gin-blog-api/pkg/setting"
	"go-gin-blog-api/pkg/util"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// 获取图片保存相对路径
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

// 获取图片保存绝对路径
func GetImageFullPath() string {
	return setting.AppSetting.PublicRootPath + GetImagePath()
}

// 获取图片链接
func GetImageFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetImagePath() + name
}

// 获取MD5图片名
func GetImageMd5Name(name string) string {
	ext := path.Ext(name)
	// 获取文件名
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMd5(fileName)
	return fileName + ext
}

// 验证图片后缀(是否是与配置允许的一致)
func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)

	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

// 验证图片上传大小(是否与配置的范围一致)
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Panicln(err)
		logging.Warn(err)
		return false
	}

	return size <= setting.AppSetting.ImageMaxSize
}

// 检查照片
func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err:%v", err)
	}
	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err:%v", err)
	}

	perm := file.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}
	return nil
}
