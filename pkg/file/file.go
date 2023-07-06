package file

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

// 获取文件大小
// mime/multipart 包， 实现了 MIME 的 multipart 解析
func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)

	return len(content), err
}

// 获取文件后缀
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// 判断文件是否存在
func CheckNotExist(src string) bool {
	_, err := os.Create(src)

	return os.IsNotExist(err)
}

// 判断文件权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

// 判断目录不存在就创建
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// 创建目录
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// 打开指定文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}
