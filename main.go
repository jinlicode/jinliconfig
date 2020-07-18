package main

import (
	"fmt"
	"jinlicode/Class"
)

func main() {

	const BasePath = "/var/jinli_deploy/"
	File := Class.CheckFileExist(BasePath + "docker-compose.yaml")
	fmt.Println(File)
	if Class.CheckFileExist(BasePath + "docker-compose.yaml") {
		fmt.Println("目录存在")
	} else {
		fmt.Println("目录不存在")
	}
	menu := Class.ConsoleOptionsSelect("run", []string{"", "重新安装（不删除数据）", "重新安装（删除数据）", "管理服务"}, "请输入选项")

	print(menu)
}
