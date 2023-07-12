package article_service

import (
	"go-gin-blog-api/pkg/file"
	"go-gin-blog-api/pkg/qrcode"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

// 文章海报
type ArticlePoster struct {
	PosterName string
	Article    *Article
	Qr         *qrcode.QrCode
}

// 海报背景
type ArticlePosterBg struct {
	Name string
	*ArticlePoster
	*Rect
	*Pt
}

type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

type Pt struct {
	X int
	Y int
}

// 实例化海报
func NewArticlePoster(posterName string, article *Article, qr *qrcode.QrCode) *ArticlePoster {
	return &ArticlePoster{
		PosterName: posterName,
		Article:    article,
		Qr:         qr,
	}
}

// 创建一个海报背景
func NewArticlePosterBg(name string, ap *ArticlePoster, rect *Rect, pt *Pt) *ArticlePosterBg {
	return &ArticlePosterBg{
		Name:          name,
		ArticlePoster: ap,
		Rect:          rect,
		Pt:            pt,
	}
}

func GetPosterFlag() string {
	return "poster"
}

// 判断海报是否生成
func (a *ArticlePoster) CheckMergedImage(path string) bool {
	if file.CheckNotExist(path+a.PosterName) == true {
		return false
	}

	return true
}

// 打开海报文件
func (a *ArticlePoster) OpenMergedImage(path string) (*os.File, error) {
	f, err := file.MustOpen(a.PosterName, path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// 生成文章海报
func (a *ArticlePosterBg) Generate() (string, string, error) {
	// 获取路径
	fullPath := qrcode.GetQrCodeFullPath()
	fileName, path, err := a.Qr.Encode(fullPath)
	if err != nil {
		return "", "", err
	}

	// 判断有无图片
	if !a.CheckMergedImage(path) {
		// 打开
		mergedF, err := a.OpenMergedImage(path)
		if err != nil {
			return "", "", err
		}
		defer mergedF.Close()

		// 强制打开
		bgF, err := file.Open(a.Name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return "", "", err
		}
		defer bgF.Close()

		// 二维码位置
		qrF, err := file.MustOpen(fileName, path)
		if err != nil {
			return "", "", err
		}
		defer qrF.Close()

		// 图片处理
		bgImage, err := jpeg.Decode(bgF)
		if err != nil {
			return "", "", err
		}

		// 图片处理
		qrImage, err := jpeg.Decode(qrF)
		if err != nil {
			return "", "", err
		}

		// 创建一个新的 RGBA 画板
		jpg := image.NewRGBA(image.Rect(a.Rect.X0, a.Rect.Y0, a.Rect.X1, a.Rect.Y1))

		// 绘制：在 RGBA 图像上绘制 背景图（bgF）
		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)

		// 绘制：在已绘制背景图的 RGBA 图像上，在指定 Point 上绘制二维码图像（qrF）
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)

		// 将绘制好的 RGBA 图像以 JPEG 4：2：0 基线格式写入合并后的图像文件（mergedF）
		jpeg.Encode(mergedF, jpg, nil)
	}

	return fileName, path, nil
}
