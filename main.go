package main

import (
	"fmt"
	"jinliconfig/class"
	"os"
)

func main() {

	//配置常量
	const (
		//基础目录配置
		BASEPATH = "/var/discuz_deploy/"
	)

	if class.CheckFileExist(BASEPATH + "docker-compose.yaml") {

		//读取docker-compose配置文件
		DockerComposeYamlRead := class.ReadFile(BASEPATH + "docker-compose.yaml")

		//读取Caddyfile文件内容
		// DockerComposeCaddyFile := class.ReadFile(BASEPATH + "config/caddy/Caddyfile")

		//获取mysql配置文件
		menumysql := class.ConsoleOptionsSelect("请选择您需要管理的数据库", class.MysqlInfo(DockerComposeYamlRead), "请输入选项")
		fmt.Println(menumysql)
		// println(DockerComposeCaddyFile)
		menu := class.ConsoleOptionsSelect("请选择您需要的服务", []string{"网站服务", "备份管理", "退出"}, "请输入选项")
		fmt.Println(menu)
	} else {
		fmt.Println("您未安装锦鲤部署，是否要安装？")
		NewInstall := class.ConsoleOptionsSelect("请输入您的选项", []string{"否", "是"}, "请选择是否重新安装")
		if NewInstall == "否" {
			os.Exit(3)
		} else {
			fmt.Println("安装过程")
		}
	}
}
