package Class

import "os"

//检查文件是否存在
func CheckFileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
