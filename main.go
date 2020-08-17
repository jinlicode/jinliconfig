package main

import (
	"fmt"
	"jinliconfig/Template"
	"jinliconfig/class"
	"os"
	"os/exec"
	"strings"
)

// https网站创建注意事项
// 1. 必须输入邮箱
// 2. 创建完成启动后，必须在基础目录执行: docker-compose exec nginx certbot -n --nginx --agree-tos -m example@example.com --domains www.example.com

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

	//检测安装的docker是否符合运行环境要求
	DockerClaim := class.ChkDokcerInstall()
	//如果不符合提示重装
	if DockerClaim != true {
		DockerReInstall := class.ConsoleUserConfirm("由于检测到您的系统不符合使用版本，需要对您的系统重新安装基础运行环境\n1. 会卸载docker\n2.重新安装符合版本的docker\n3.dockers原始文件目录/etc/docker改为/etc/docker_bak\n4.docker原始目录/var/lib/docker改为/var/lib/docker_bak\n 请问您是否同意重新安装docker？")
		if DockerReInstall {
			class.ExecDockerInstall()
		} else {
			os.Exit(3)
		}
	}

CreateNewSiteFlag:

	if class.CheckFileExist(BASEPATH + "docker-compose.yaml") {

		//读取docker-compose配置文件
		DockerComposeYamlRead := class.ReadFile(BASEPATH + "docker-compose.yaml")
		DockerComposeYamlMap := class.YamlFileToMap(DockerComposeYamlRead)

		//获取已经存在的网站
		ExistSiteSlice := []string{}
		for k := range DockerComposeYamlMap["services"].(map[string]interface{}) {
			if k != "nginx" && k != "memcached" && k != "mysql" && k != "php" {
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
	ServiceSelectFlag:
		ServiceSelect := class.ConsoleOptionsSelect("请选择您需要的服务", []string{"网站服务", "备份管理", "权限修复", "木马查杀", "退出"}, "请输入选项")
		switch ServiceSelect {
		case "网站服务":
			//网站服务选择主菜单
		WebServiceSelectFlag:
			WebServiceSelectOption := []string{}
			WebServiceSelectOption = append(ExistSiteSlice, "新增网站", "返回上层")
			WebServiceSelect := class.ConsoleOptionsSelect("请选择您需要管理的网站", WebServiceSelectOption, "请输入选项")
			switch WebServiceSelect {
			case "返回上层":
				fmt.Println("返回上层")
				goto ServiceSelectFlag
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
							goto WebServiceSelectFlag
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

				//自动创建以网站名字命名的程序目录
				class.ExecLinuxCommand("mkdir " + BASEPATH + "/code/" + newDomain)

				//创建网站的配置文件到对应的config配置文件中
				class.ExecLinuxCommand("mkdir " + BASEPATH + "/config/php/" + newDomain)
				class.WriteFile(BASEPATH+"config/php"+newDomain+"/www.conf", Template.PhpWww())
				class.WriteFile(BASEPATH+"config/php"+newDomain+"/php.ini", Template.PhpIni())
				class.WriteFile(BASEPATH+"config/php"+newDomain+"/php-fpm.conf", Template.PhpFpm())

				//写入docker-compose.yaml 文件
				NewDockerComposeYamlString, _ := class.MapToYaml(DockerComposeYamlMap)
				class.WriteFile(BASEPATH+"docker-compose.yaml", NewDockerComposeYamlString)
				// fmt.Println(NewDockerComposeYamlString)

			case WebServiceSelect:
				WebConfigSelect := class.ConsoleOptionsSelect("请选择您需要管理的网站服务", []string{
					WebServiceSelect + "的" + "nginx配置",
					WebServiceSelect + "的" + "php配置",
					WebServiceSelect + "的" + "数据库配置",
					"暂停" + WebServiceSelect + "网站服务",
					"删除" + WebServiceSelect + "的网站(不删除数据)",
					"删除" + WebServiceSelect + "的网站(删除数据，包含数据库和程序)",
					"返回上层"}, "请输入选项")
				//网站内服务修改主菜单
			WebConfigSelectFlag:
				switch WebConfigSelect {
				case "返回上层":
					fmt.Println("返回上层")
					goto WebServiceSelectFlag
				case WebServiceSelect + "的" + "nginx配置":
					fmt.Println(WebServiceSelect + "的" + "nginx配置")
					switch WebConfigSelect {
					case "返回上层":
						fmt.Println("返回上层")
						goto WebConfigSelectFlag
					}
				case WebServiceSelect + "的" + "php配置":
					fmt.Println(WebServiceSelect + "的" + "php配置")
				case WebServiceSelect + "的" + "数据库配置":
					fmt.Println(WebServiceSelect + "的" + "数据库配置")
					fmt.Println(WebServiceSelect + "的" + "php配置")
				case "暂停" + WebServiceSelect + "网站服务":
					//确定是否需要暂停
					ReStopSiteConfirm := class.ConsoleUserConfirm("确定暂停" + WebServiceSelect + "服务吗？")
					if ReStopSiteConfirm != true {
						fmt.Println("已取消操作")
						goto WebServiceSelectFlag
					}

					//输入命令 暂停容器
					// fmt.Println("cd " + BASEPATH + " && docker-compose stop " + strings.Replace(WebServiceSelect, ".", "_", -1))
					cmd := exec.Command("cd " + BASEPATH + " && docker-compose stop " + strings.Replace(WebServiceSelect, ".", "_", -1))
					cmd.Stdout = os.Stdout
					cmd.Run()
					fmt.Println("暂停成功")

				case "删除" + WebServiceSelect + "的网站(不删除数据)":
					//确定是否需要删除
					ReDelSiteConfirm := class.ConsoleUserConfirm("确定删除" + WebServiceSelect + "网站吗？(不删除数据)")
					if ReDelSiteConfirm != true {
						fmt.Println("已取消操作")
						goto WebServiceSelectFlag
					}
					//输入命令 删除yaml中服务
					MapKey := strings.Replace(WebServiceSelect, ".", "_", -1)
					class.ExecLinuxCommand("cd " + BASEPATH + " && docker-compose stop " + MapKey + " && docker-compose rm " + MapKey)

					//执行完之后删除yaml中对应的map
					delete(DockerComposeYamlMap["services"].(map[string]interface{}), MapKey)

					//重新写入到yaml
					NewDockerComposeYamlString, _ := class.MapToYaml(DockerComposeYamlMap)
					class.WriteFile(BASEPATH+"docker-compose.yaml", NewDockerComposeYamlString)
					fmt.Println("删除成功")

				case "删除" + WebServiceSelect + "的网站(删除数据，包含数据库和程序)":
					//确定是否需要删除
					ReDelSiteConfirm := class.ConsoleUserConfirm("确定删除" + WebServiceSelect + "网站吗？(删除数据，包含数据库和程序)")
					if ReDelSiteConfirm != true {
						fmt.Println("已取消操作")
						goto WebServiceSelectFlag
					}
					//输入命令 删除yaml中服务
					MapKey := strings.Replace(WebServiceSelect, ".", "_", -1)
					// fmt.Println("cd " + BASEPATH + " && docker-compose stop " + MapKey + " && docker-compose rm " + MapKey)
					class.ExecLinuxCommand("cd " + BASEPATH + " && docker-compose stop " + MapKey + " && docker-compose rm " + MapKey)

					//执行完之后删除yaml中对应的map
					delete(DockerComposeYamlMap["services"].(map[string]interface{}), MapKey)

					//重新写入到yaml
					NewDockerComposeYamlString, _ := class.MapToYaml(DockerComposeYamlMap)
					class.WriteFile(BASEPATH+"docker-compose.yaml", NewDockerComposeYamlString)

					//操作删除工作 删除代码目录  删除  数据库 drop database bbbbbbbb
					MysqlPassword := DockerComposeYamlMap["services"].(map[string]interface{})["mysql"].(map[string]interface{})["environment"].(map[string]interface{})["MYSQL_ROOT_PASSWORD"]
					class.ExecLinuxCommand("rm -rf " + BASEPATH + "/" + MapKey + " && docker-compose exec mysql bash -c \"mysql -uroot -p" + MysqlPassword.(string) + " -e 'drop database " + MapKey + "'\"")

					fmt.Println("删除成功")

				}

			}

		case "备份管理":
			fmt.Println("备份管理")
			goto ServiceSelectFlag
		case "权限修复":
			fmt.Println("权限修复")
			goto ServiceSelectFlag
		case "木马查杀":
			fmt.Println("木马查杀")
			goto ServiceSelectFlag
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
			//协议内容
			AgreementText := Template.Agreement()
			fmt.Println(AgreementText)
			Agree := class.ConsoleOptionsSelect("您是否同意此协议", []string{"否", "是"}, "请选择是否同意协议")
			if Agree == "否" {
				fmt.Println("您未同意安装协议，程序已退出")
				os.Exit(3)
			}
			fmt.Println("开始安装程序，请稍等...")

			//执行安装docker
			class.ExecDockerInstall()

			//创建项目目录
			class.ExecLinuxCommand("mkdir " + BASEPATH)

			//创建代码目录
			class.ExecLinuxCommand("mkdir " + BASEPATH + "code/")

			//创建各配置项目录
			class.ExecLinuxCommand("mkdir " + BASEPATH + "/config/ && mkdir " + BASEPATH + "/config/cert/ && mkdir " + BASEPATH + "/config/mysql/ && mkdir " + BASEPATH + "/config/nginx/ && mkdir " + BASEPATH + "/config/php/")

			//设置代码目录为 10000,10000
			class.ExecLinuxCommand("chown -R 10000:10000 " + BASEPATH + "code/")

			//自动创建yaml标准模版

			// fmt.Println("安装过程")

			DockerComposeVersion := Template.DockerComposeVersion()
			DockerComposeNginxMap := class.YamlFileToMap(Template.DockerComposeNginx())
			DockerComposeMysqlMap := class.YamlFileToMap(Template.DockerComposeMysql())
			DockerComposeNetWorksMap := class.YamlFileToMap(Template.DockerComposeNetWorks())

			DockerComposeMap := make(map[string]interface{})
			DockerComposeMap["version"] = DockerComposeVersion
			DockerComposeMap["services"] = make(map[string]interface{})
			DockerComposeMap["services"].(map[string]interface{})["nginx"] = DockerComposeNginxMap
			DockerComposeMap["services"].(map[string]interface{})["mysql"] = DockerComposeMysqlMap
			DockerComposeMap["networks"] = DockerComposeNetWorksMap

			//写入yaml文件 跳转到新建网站
			DockerComposeYamlString, _ := class.MapToYaml(DockerComposeMap)
			class.WriteFile(BASEPATH+"docker-compose.yaml", DockerComposeYamlString)
			goto CreateNewSiteFlag
		}
	}
}
