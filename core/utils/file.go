package utils

import (
	"errors"
	"os"
)

//判断目录是否存在
func DirExists(path string) (bool, error) {
	file, err := os.Stat(path)
	if err == nil {
		//存在
		if file.IsDir() {
			return true, nil
		}
		return false, errors.New("文件已存在")
	}

	if os.IsNotExist(err) {
		//不存在
		return false, err
	}

	//不确定是否存在
	return false, err
}

func Mkdir() {

}
