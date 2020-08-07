package class

import (
	"fmt"
	"io/ioutil"
	"os"
)

//检查文件是否存在，输入路径返回布尔值
func CheckFileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

//读取文件，传入文件路径，返回文件内容
func ReadFile(path string) string {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%s\n", err)
		panic(err)
	}
	return string(f)
}

//写入文件，传入文件路径
func WriteFile(path string, content string) bool {
	err := ioutil.WriteFile(path, []byte(content), 0666)

	if err != nil {
		fmt.Printf("%s\n", err)
		panic(err)
	}

	return true
}
