package main

import (
	"fmt"
	"jinliconfig/Class"
	"os"
)

func main() {

	//配置常量
	const (
		//基础目录配置
		BASEPATH = "/var/discuz_deploy/"
	)

	if Class.CheckFileExist(BASEPATH + "docker-compose.yaml") {

		//读取docker-compose配置文件
		DockerComposeYamlRead := Class.ReadFile(BASEPATH + "docker-compose.yaml")

		//读取Caddyfile文件内容
		// DockerComposeCaddyFile := Class.ReadFile(BASEPATH + "config/caddy/Caddyfile")

		//获取mysql配置文件
		DockerComposeMysqlConfig := Class.MysqlInfo(DockerComposeYamlRead)

		fmt.Println(DockerComposeMysqlConfig)
		// println(DockerComposeCaddyFile)
		menu := Class.ConsoleOptionsSelect("请选择您需要的服务", []string{"网站服务", "备份管理", "退出"}, "请输入选项")
		fmt.Println(menu)
	} else {
		fmt.Println("您未安装锦鲤部署，是否要安装？")
		NewInstall := Class.ConsoleOptionsSelect("请输入您的选项", []string{"否", "是"}, "请选择是否重新安装")
		if NewInstall == "否" {
			os.Exit(3)
		} else {
			fmt.Println("安装过程")
		}
	}
}
