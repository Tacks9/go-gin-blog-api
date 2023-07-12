package article_service

import (
	"go-gin-blog-api/models"
	"go-gin-blog-api/pkg/file"
	"go-gin-blog-api/pkg/qrcode"
	"go-gin-blog-api/pkg/setting"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
)

// 文章海报
type ArticlePoster struct {
	PosterName string
	Article    *models.Article
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

// 文字
type DrawText struct {
	JPG    draw.Image
	Merged *os.File

	Title string
	X0    int
	Y0    int
	Size0 float64

	SubTitle string
	X1       int
	Y1       int
	Size1    float64
}

// 实例化海报
func NewArticlePoster(posterName string, article *models.Article, qr *qrcode.QrCode) *ArticlePoster {
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

		// 创建文案
		err = a.DrawPoster(&DrawText{
			JPG:    jpg,
			Merged: mergedF,

			Title: a.Article.Title,
			X0:    100,
			Y0:    180,
			Size0: 42,

			SubTitle: "---" + a.Article.CreatedBy,
			X1:       400,
			Y1:       240,
			Size1:    36,
		}, "msyhbd.ttc")

		if err != nil {
			return "", "", err
		}

		// 将绘制好的 RGBA 图像以 JPEG 4：2：0 基线格式写入合并后的图像文件（mergedF）
		jpeg.Encode(mergedF, jpg, nil)
	}

	return fileName, path, nil
}

// 海报写文字
func (a *ArticlePosterBg) DrawPoster(d *DrawText, fontName string) error {
	// 字体文件的完整路径
	fontSource := setting.AppSetting.PublicRootPath + setting.AppSetting.FontSavePath + fontName

	// 读取内容为字节切片
	fontSourceBytes, err := ioutil.ReadFile(fontSource)
	if err != nil {
		return err
	}

	// 解析字体文件的字节切片为 TrueType 字体
	trueTypeFont, err := freetype.ParseFont(fontSourceBytes)
	if err != nil {
		return err
	}

	// freetype.NewContext：创建一个新的 Context
	fc := freetype.NewContext()
	// 屏幕每英寸的分辨率
	fc.SetDPI(72)
	// 设置用于绘制文本的字体
	fc.SetFont(trueTypeFont)
	// 以磅为单位设置字体大小
	fc.SetFontSize(d.Size0)
	// 设置剪裁矩形以进行绘制
	fc.SetClip(d.JPG.Bounds())
	// 目标图像
	fc.SetDst(d.JPG)
	// 设置绘制操作的源图颜色
	fc.SetSrc(image.White)
	// fc.SetSrc(image.Black)
	// red := color.RGBA{255, 0, 0, 255} // 红色 (R: 255, G: 0, B: 0, A: 255)
	// fc.SetSrc(image.NewUniform(red))

	// 设置坐标
	pt := freetype.Pt(d.X0, d.Y0)
	// 根据 Pt 的坐标值绘制给定的文本内容
	_, err = fc.DrawString(d.Title, pt)
	if err != nil {
		return err
	}

	fc.SetFontSize(d.Size1)
	_, err = fc.DrawString(d.SubTitle, freetype.Pt(d.X1, d.Y1))
	if err != nil {
		return err
	}

	// 使用 jpeg.Encode 将 d.JPG 保存到 d.Merged 中
	err = jpeg.Encode(d.Merged, d.JPG, nil)
	if err != nil {
		return err
	}

	return nil
}
