package api

import (
	"go-gin-blog-api/pkg/e"
	"go-gin-blog-api/pkg/logging"
	"go-gin-blog-api/pkg/upload"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 上传图片
func UploadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]string)

	// 表单上传的第一个文件
	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
		return
	}

	// 如果为空
	if image == nil {
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageMd5Name(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()

		// 要保存的路径
		imageSrc := fullPath + imageName

		// 检查后缀和大小
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				// 检查图片所需目录权限等
				logging.Warn(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, imageSrc); err != nil {
				// 上传文件
				logging.Warn(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				// 上传成功
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
