package main

import (
	"fmt"
	"jinliconfig/Template"
	"jinliconfig/class"
	"os"
	"os/exec"
	"strings"
)

func main() {

	//配置常量
	const (
		//基础目录配置
		BASEPATH = "/var/discuz_deploy/"
	)
	//检查是否为root启动
	if os.Getuid() != 0 {
		fmt.Println("您的用户权限太低，请使用root用户执行，命令为：sudo su")
		os.Exit(3)
	}

	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()

	Welcome := `

		     ██╗██╗███╗   ██╗██╗     ██╗ ██████╗ ██████╗ ███╗   ██╗███████╗██╗ ██████╗ 
		     ██║██║████╗  ██║██║     ██║██╔════╝██╔═══██╗████╗  ██║██╔════╝██║██╔════╝ 
		     ██║██║██╔██╗ ██║██║     ██║██║     ██║   ██║██╔██╗ ██║█████╗  ██║██║  ███╗
		██   ██║██║██║╚██╗██║██║     ██║██║     ██║   ██║██║╚██╗██║██╔══╝  ██║██║   ██║
		╚█████╔╝██║██║ ╚████║███████╗██║╚██████╗╚██████╔╝██║ ╚████║██║     ██║╚██████╔╝
		╚════╝ ╚═╝╚═╝  ╚═══╝╚══════╝╚═╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝     ╚═╝ ╚═════╝ 
																			  
						欢迎使用锦鲤网站管理系统 v1.1
		`
	fmt.Println(Welcome)

	if class.CheckFileExist(BASEPATH + "docker-compose.yaml") {

		//读取docker-compose配置文件
		DockerComposeYamlRead := class.ReadFile(BASEPATH + "docker-compose.yaml")
		DockerComposeYamlMap := class.YamlFileToMap(DockerComposeYamlRead)

		//获取已经存在的网站
		ExistSiteSlice := []string{}
		for k := range DockerComposeYamlMap["services"].(map[string]interface{}) {
			if k != "caddy" && k != "memcached" && k != "mysql" && k != "php" {
				ExistSiteSlice = append(ExistSiteSlice, strings.Replace(k, "_", ".", -1))
			}
		}
		fmt.Println(ExistSiteSlice)
		// fmt.Printf("%v\n", DockerComposeYamlMap["networks"].(map[string]interface{})["discuz"])

		// fmt.Println(DockerComposeYamlRead)

		//读取Caddyfile文件内容
		// DockerComposeCaddyFile := class.ReadFile(BASEPATH + "config/caddy/Caddyfile")

		//获取mysql配置文件
		// menumysql := class.ConsoleOptionsSelect("请选择您需要管理的数据库", class.MysqlInfo(DockerComposeYamlRead), "请输入选项")
		// fmt.Println(menumysql)
		// println(DockerComposeCaddyFile)

		//服务选择主菜单
	ServiceSelect:
		ServiceSelect := class.ConsoleOptionsSelect("请选择您需要的服务", []string{"网站服务", "备份管理", "退出"}, "请输入选项")
		switch ServiceSelect {
		case "网站服务":
			//网站服务选择主菜单
		WebServiceSelect:
			WebServiceSelectOption := []string{}
			WebServiceSelectOption = append(ExistSiteSlice, "新增网站", "返回上层")
			WebServiceSelect := class.ConsoleOptionsSelect("请选择您需要管理的网站", WebServiceSelectOption, "请输入选项")
			switch WebServiceSelect {
			case "返回上层":
				fmt.Println("返回上层")
				goto ServiceSelect
			case "新增网站":
			ReInputSiteDomainFlag:
				NewSiteDomain := class.ConsoleUserInput("请输入您需要添加的域名：")
				NewSiteDomain = strings.TrimSpace(NewSiteDomain)

				//检测网站域名是否输入正确
				if !class.CheckDomain(NewSiteDomain) {
					fmt.Println("您输入的域名不正确，已退出操作请重新输入！")
					goto ReInputSiteDomainFlag
				}

				NewSiteHTTPS := class.ConsoleOptionsSelect("是否使用HTTPS", []string{"是", "否"}, "请输入选项")
				if NewSiteHTTPS == "否" {
					fmt.Println("您选择了没有https证书，如果选择错误请按Ctrl+C结束当前进程")
				} else {
					NewSiteSSLHave := class.ConsoleOptionsSelect("您是否有自己的证书", []string{"是", "否"}, "请输入选项")
					if NewSiteSSLHave == "否" {
						fmt.Println("您选择了没有https证书，我们将会自动为您创建HTTPS证书，请您先一步解析域名到您的服务器上，如果使用CDN请参考官方帮助文档：https://xxxxxxxxxxxxxxx")
					} else {
						fmt.Println("请您准备好证书需要用到的两个文件，如果有选择请选择下载nginx使用版本，马上您会被要求粘贴两个文件内容")
						NewSiteSSLHaveConfirm := class.ConsoleUserConfirm("您是否已经准备好证书，如果准备好请选择")
						if NewSiteSSLHaveConfirm == true {
							CertCERInput := class.ConsoleUserText("输入证书CER文件内容，写入完成请按两次回车即可")
							CertKEYInput := class.ConsoleUserText("输入证书KEY文件内容，写入完成请按两次回车即可")
							fmt.Println(CertCERInput)
							fmt.Println(CertKEYInput)
						} else {
							goto WebServiceSelect
						}
					}
				}
				NewSitePhpVersion := class.ConsoleOptionsSelect("请选择您需要的php版本", []string{"5.6", "7.0", "7.1", "7.2", "7.3", "7.4", "8.0"}, "请输入选项")

				//再回显一次输入的内容判断是否真的要开始安装
				LastReConfirm := class.ConsoleUserConfirm("\n域名：[" + NewSiteDomain + "]\n是否启用https：[" + NewSiteHTTPS + "]\nphp版本：[" + NewSitePhpVersion + "]\n确定是否立即安装")
				if LastReConfirm != true {
					fmt.Println("已取消操作")
					os.Exit(3)
				}

				//获取php 镜像模版
				SitePhpVersionCompose := Template.DockerComposePhp()

				newDomain := NewSiteDomain
				newDomain = strings.Replace(newDomain, ".", "_", -1)

				//替换域名
				SitePhpVersionCompose = strings.Replace(SitePhpVersionCompose, "- ./code:", "- ./code/"+newDomain+":", 1)

				//替换php版本
				SitePhpVersionCompose = strings.Replace(SitePhpVersionCompose, "image: jinlicode/discuz_docker:latest", "image: jinlicode/discuz_docker:"+NewSitePhpVersion, 1)

				//生成子map
				NewSitePhpVersionComposeMap := class.YamlFileToMap(SitePhpVersionCompose)

				//插入到总的yaml文件中
				DockerComposeYamlMap["services"].(map[string]interface{})[newDomain] = NewSitePhpVersionComposeMap

				//写入docker-compose.yaml 文件
				NewDockerComposeYamlString, _ := class.MapToYaml(DockerComposeYamlMap)
				class.WriteFile(BASEPATH+"docker-compose.yaml", NewDockerComposeYamlString)
				// fmt.Println(NewDockerComposeYamlString)

			case WebServiceSelect:
				WebConfigSelect := class.ConsoleOptionsSelect("请选择您需要管理的网站服务", []string{WebServiceSelect + "的" + "nginx配置", WebServiceSelect + "的" + "php配置", WebServiceSelect + "的" + "数据库配置", "返回上层"}, "请输入选项")
				//网站内服务修改主菜单
			WebConfigSelect:
				switch WebConfigSelect {
				case "返回上层":
					fmt.Println("返回上层")
					goto WebServiceSelect
				case WebServiceSelect + "的" + "nginx配置":
					fmt.Println(WebServiceSelect + "的" + "nginx配置")
					switch WebConfigSelect {
					case "返回上层":
						fmt.Println("返回上层")
						goto WebConfigSelect
					}
				case WebServiceSelect + "的" + "php配置":
					fmt.Println(WebServiceSelect + "的" + "php配置")
				case WebServiceSelect + "的" + "数据库配置":
					fmt.Println(WebServiceSelect + "的" + "数据库配置")
				}

			}
		case "备份管理":
			fmt.Println("备份管理")
		case "木马查杀":
			fmt.Println("木马查杀")
		case "退出":
			fmt.Println("退出")
			break //可以添加
			// goto WebServiceSelect //随便添加的一会删除
		}
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
