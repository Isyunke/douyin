package filex

import (
	"os"
	"path"
	"strings"
)

// PathExists 判断是否存在路径
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileNameWithOutExt(fileName string) string {
	ext := path.Ext(fileName)
	return strings.TrimSuffix(fileName, ext)
}
