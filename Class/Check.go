package Class

import (
	"fmt"
	"io/ioutil"
	"os"
)

//检查文件是否存在
func CheckFileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

//读取文件，传入文件路径
func ReadFile(path string) string {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%s\n", err)
		panic(err)
	}
	return string(f)
}
