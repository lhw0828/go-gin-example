package upload

import (
	"fmt"
	"github.com/lhw0828/go-gin-example/pkg/file"
	"github.com/lhw0828/go-gin-example/pkg/logging"
	"github.com/lhw0828/go-gin-example/pkg/setting"
	"github.com/lhw0828/go-gin-example/pkg/util"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// GetImageFullUrl 获取图片完整的url
func GetImageFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetImagePath() + name
}

// GetImageName 获取图片名称
func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

// GetImagePath 获取图片保存路径
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

// GetImageFullPath 获取图片完整的保存路径
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

// CheckImageExt 检查图片后缀是否合法
func CheckImageExt(fileName string) bool {
	ext := path.Ext(fileName)
	ext = strings.ToUpper(ext)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == ext {
			return true
		}
	}
	return false
}

// CheckImageSize 检查图片大小是否合法
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}
	return size <= setting.AppSetting.ImageMaxSize
}

// CheckImage 检查图片是否合法
func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}
	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}
	perm := file.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}
	return nil
}
