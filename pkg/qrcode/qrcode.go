package qrcode

import (
	"fmt"
	"go-gin-blog-api/pkg/file"
	"go-gin-blog-api/pkg/logging"
	"go-gin-blog-api/pkg/setting"
	"go-gin-blog-api/pkg/util"
	"image/jpeg"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type QrCode struct {
	URL    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const (
	EXT_JPG = ".jpg"
)

// 实例化
func NewQrCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:    url,
		Width:  width,
		Height: height,
		Level:  level,
		Mode:   mode,
		Ext:    EXT_JPG,
	}
}

// 保存相对路径
func GetQrCodePath() string {
	return setting.AppSetting.QrCodeSavePath
}

// 保存
func GetQrCodeFullPath() string {
	return setting.AppSetting.PublicRootPath + GetQrCodePath()
}

// 二维码 URL 地址
func GetQrCodeFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetQrCodePath() + name
}

// 二维码文件名
func GetQrCodeFileName(value string) string {
	return util.EncodeMd5(value)
}

// 二维码后缀
func (q *QrCode) GetQrCodeExt() string {
	return q.Ext
}

// 判定是否存在
func (q *QrCode) CheckEncode(path string) bool {
	src := path + GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	if file.CheckNotExist(src) == true {
		return false
	}

	return true
}

// 编码处理
func (q *QrCode) Encode(path string) (string, string, error) {
	name := GetQrCodeFileName(q.URL) + q.GetQrCodeExt()

	// 获取二维码生成路径
	src := path + name
	fmt.Println(src)
	if file.CheckNotExist(src) == true {
		fmt.Println(11)
		// 创建二维码
		code, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			logging.Info(err)
			return "", "", err
		}

		// 缩放二维码到指定大小
		code, err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			logging.Info(err)
			return "", "", err
		}

		// 新建存放二维码图片的文件
		f, err := file.MustOpen(name, path)
		if err != nil {
			logging.Info(err)
			return "", "", err
		}
		defer f.Close()

		// 图片和二维码合体，设置图像质量，默认值为 75
		err = jpeg.Encode(f, code, nil)
		if err != nil {
			logging.Info(err)
			return "", "", err
		}
	}

	return name, path, nil
}
