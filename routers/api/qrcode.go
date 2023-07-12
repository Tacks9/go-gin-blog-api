package api

import (
	"go-gin-blog-api/pkg/app"
	"go-gin-blog-api/pkg/e"
	"go-gin-blog-api/pkg/qrcode"
	"net/http"

	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
)

const (
	QRCODE_URL = "https://github.com/Tacks9/go-gin-blog-api"
)

// 二维码生成示例
func Generate(c *gin.Context) {
	appG := app.Gin{C: c}
	qrc := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto)
	path := qrcode.GetQrCodeFullPath()
	codeName, _, err := qrc.Encode(path)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	// 返回数据
	data := make(map[string]string)
	data["filename"] = codeName
	data["url"] = qrcode.GetQrCodeFullUrl(codeName)

	appG.Response(http.StatusOK, e.SUCCESS, data)
}
