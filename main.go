package main

import (
	"fmt"
	"jinliconfig/Template"
	"jinliconfig/class"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// https网站创建注意事项
// 1. 必须输入邮箱
// 2. 创建完成启动后，必须在基础目录执行: docker-compose exec nginx certbot -n --nginx --agree-tos -m example@example.com --domains www.example.com

func main() {

	//配置常量
	const (
		//基础目录配置
		BASEPATH = "/var/jinli/"
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
		DockerReInstall := class.ConsoleUserConfirm("由于检测到您的系统不符合使用版本，需要对您的系统重新安装基础运行环境\n1. 会卸载docker\n2.重新安装符合版本的docker\n3.dockers原始文件目录/etc/docker将被删除\n4.docker原始目录/var/lib/docker将被删除\n 请问您是否同意重新安装docker？")
		if DockerReInstall {
			//删除老的docker 重新安装特定的docker
			class.ChkDokcerRemove()
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

		//获取最大内网数
		SiteNetMax := 2
		for k, v := range DockerComposeYamlMap["services"].(map[string]interface{}) {
			if k != "nginx" && k != "memcached" && k != "mysql" && k != "php" {
				ExistSiteSlice = append(ExistSiteSlice, strings.Replace(k, "_", ".", -1))

				//获取内网最大数字
				max := strings.Split(v.(map[string]interface{})["networks"].(map[string]interface{})["jinli_net"].(map[string]interface{})["ipv4_address"].(string), ".")
				maxNumString := max[3]

				maxNum, err := strconv.Atoi(maxNumString)

				if err == nil && SiteNetMax < maxNum {
					SiteNetMax = maxNum
				}
			}
		}
		fmt.Println(ExistSiteSlice)
		// fmt.Printf("%v\n", DockerComposeYamlMap["networks"].(map[string]interface{})["jinli_net"])

		// fmt.Println(DockerComposeYamlRead)

		//读取Caddyfile文件内容
		// DockerComposeCaddyFile := class.ReadFile(BASEPATH + "config/caddy/Caddyfile")

		//获取mysql配置文件
		// menumysql := class.ConsoleOptionsSelect("请选择您需要管理的数据库", class.MysqlInfo(DockerComposeYamlRead), "请输入选项")
		// fmt.Println(menumysql)
		// println(DockerComposeCaddyFile)

		//服务选择主菜单
	ServiceSelectFlag:
		ServiceSelect := class.ConsoleOptionsSelect("请选择您需要的服务", []string{"网站服务", "备份管理", "权限修复", "升级系统镜像", "退出"}, "请输入选项")
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

				//域名转下划线
				newDomain := NewSiteDomain
				newDomain = strings.Replace(newDomain, ".", "_", -1)

				//检测nginx是否已经存在配置文件
				if class.CheckFileExist(BASEPATH + "config/nginx/" + newDomain + ".conf") {
					if class.ConsoleOptionsSelect("您已经存在当前配置文件，如您选择继续，将覆盖配置文件", []string{"是", "否"}, "请输入选项") == "否" {
						goto ReInputSiteDomainFlag
					}
				}

				NewSiteHTTPS := class.ConsoleOptionsSelect("是否使用HTTPS", []string{"是", "否"}, "请输入选项")
				NewSiteSSLEmail := ""
				if NewSiteHTTPS == "否" {
					fmt.Println("您选择了没有https证书，如果选择错误请按Ctrl+C结束当前进程")
				} else {
					fmt.Println("您选择了没有https证书，我们将会自动为您创建HTTPS证书，请您先一步解析域名到您的服务器上，如果使用CDN请参考官方帮助文档：https://xxxxxxxxxxxxxxx")

				ReInputSiteEmailFlag:
					//开始输入邮箱
					NewSiteSSLEmail = class.ConsoleUserInput("请输入您的邮箱地址，此地址为了自动申请证书所用：")
					NewSiteSSLEmail = strings.TrimSpace(NewSiteSSLEmail)

					//检测邮箱是否输入正确
					if !class.CheckEmail(NewSiteSSLEmail) {
						fmt.Println("您输入的邮箱不正确，请重新输入！")
						goto ReInputSiteEmailFlag
					}

					// 	certPEMBlock, _ := ioutil.ReadFile("/var/discuz_deploy/config/cert/live/test1.jinli.plus/cert.pem")
					//     certDERBlockde, _ := pem.Decode(certPEMBlock)
					//     x509Cert, _ := x509.ParseCertificate(certDERBlockde.Bytes)
					// 	println(x509Cert.NotAfter.Format("2006-01-02 15:04:05"))

				}
				NewSitePhpVersion := class.ConsoleOptionsSelect("请选择您需要的php版本, sec版本为安全版本", []string{
					"5.6",
					"5.6-sec",
					"7.0",
					"7.0-sec",
					"7.1",
					"7.1-sec",
					"7.2",
					"7.2-sec",
					"7.3",
					"7.3-sec",
				}, "请输入选项")

				//再回显一次输入的内容判断是否真的要开始安装
				LastReConfirm := class.ConsoleUserConfirm("\n域名：[" + NewSiteDomain + "]\n是否启用https：[" + NewSiteHTTPS + "]\nphp版本：[" + NewSitePhpVersion + "]\n确定是否立即安装")
				if LastReConfirm != true {
					fmt.Println("已取消操作")
					os.Exit(3)
				}

				fmt.Println("新网站正在建设中，请您稍等......")

				//获取php 镜像模版
				SitePhpVersionCompose := Template.DockerComposePhp()

				//替换域名
				SitePhpVersionCompose = strings.Replace(SitePhpVersionCompose, "www_example_com", newDomain, -1)

				//替换内网ip
				SitePhpVersionCompose = strings.Replace(SitePhpVersionCompose, "ipv4_address: 10.99.2.2", "ipv4_address: 10.99.2."+strconv.Itoa(SiteNetMax+1), -1)

				//替换php版本
				SitePhpVersionCompose = strings.Replace(SitePhpVersionCompose, "jinlicode/php:latest", "jinlicode/php:v"+NewSitePhpVersion, 1)

				//生成子map
				NewSitePhpVersionComposeMap := class.YamlFileToMap(SitePhpVersionCompose)

				//插入到总的yaml文件中
				DockerComposeYamlMap["services"].(map[string]interface{})[newDomain] = NewSitePhpVersionComposeMap

				//自动创建以网站名字命名的程序目录
				class.ExecLinuxCommand("mkdir " + BASEPATH + "code/" + newDomain)
				class.ExecLinuxCommand("mkdir " + BASEPATH + "config/php/" + newDomain)

				//写入404以及index文件到置顶目录
				class.WriteFile(BASEPATH+"code/"+newDomain+"/index.html", Template.HTMLIndex())
				class.WriteFile(BASEPATH+"code/"+newDomain+"/404.html", Template.HTML404())

				//创建网站的配置文件到对应的config配置文件中
				class.ExecLinuxCommand("mkdir " + BASEPATH + "config/php/" + newDomain)
				class.WriteFile(BASEPATH+"config/php/"+newDomain+"/www.conf", Template.PhpWww())
				class.WriteFile(BASEPATH+"config/php/"+newDomain+"/php.ini", Template.PhpIni())
				class.WriteFile(BASEPATH+"config/php/"+newDomain+"/php-fpm.conf", Template.PhpFpm())

				//创建对应nginx.conf到对应目录
				if NewSiteHTTPS == "否" {
					TemplateNginxHTTPString := Template.TemplateNginxHttp()
					TemplateNginxHTTPString = strings.Replace(TemplateNginxHTTPString, "www_example_com", newDomain, -1)
					TemplateNginxHTTPString = strings.Replace(TemplateNginxHTTPString, "www.example.com", NewSiteDomain, -1)
					TemplateNginxHTTPString = strings.Replace(TemplateNginxHTTPString, "php:9000", newDomain+":9000", -1)
					class.WriteFile(BASEPATH+"config/nginx/"+newDomain+".conf", TemplateNginxHTTPString)
				} else {

					TemplateNginxHTTPSString := Template.TemplateNginxHttps()
					TemplateNginxHTTPSString = strings.Replace(TemplateNginxHTTPSString, "www_example_com", newDomain, -1)
					TemplateNginxHTTPSString = strings.Replace(TemplateNginxHTTPSString, "www.example.com", NewSiteDomain, -1)
					TemplateNginxHTTPSString = strings.Replace(TemplateNginxHTTPSString, "php:9000", newDomain+":9000", -1)

					//如果是手动输入 保存cert.key
					class.WriteFile(BASEPATH+"config/nginx/"+newDomain+".conf", TemplateNginxHTTPSString)

				}

				//写入docker-compose.yaml 文件
				NewDockerComposeYamlString, _ := class.MapToYaml(DockerComposeYamlMap)
				class.WriteFile(BASEPATH+"docker-compose.yaml", NewDockerComposeYamlString)

				//启动新网站服务
				class.ExecLinuxCommand("cd " + BASEPATH + " && docker-compose up -d " + newDomain)
				//重启nginx 配置
				class.ExecLinuxCommand("cd " + BASEPATH + " && docker-compose exec nginx nginx -s reload")

				//重启命令
				if NewSiteHTTPS == "是" {
					class.ExecLinuxCommand("cd " + BASEPATH + " && docker-compose exec nginx certbot -n --nginx --agree-tos -m " + NewSiteSSLEmail + " --domains " + NewSiteDomain)
				}

				//自动创建网站对应mysql数据
				MysqlRootPassword := DockerComposeYamlMap["services"].(map[string]interface{})["mysql"].(map[string]interface{})["environment"].(map[string]interface{})["MYSQL_ROOT_PASSWORD"]
				MysqlRootPasswordString := MysqlRootPassword.(string)

				//获取随机密码
				mysqlSiteRandPassword := class.RandomString(16)

				//自动创建数据库 用户名 密码
				class.CreateDatabase(BASEPATH, MysqlRootPasswordString, newDomain, newDomain, mysqlSiteRandPassword)

				//显示新网站内容
				fmt.Println("请将您的网站代码上传至 " + BASEPATH + "code/" + newDomain)
				fmt.Println("数据库名称：" + newDomain)
				fmt.Println("mysql用户名：" + newDomain)
				fmt.Println("mysql密码：" + mysqlSiteRandPassword)

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

					//删除对应的nginx配置
					class.ExecLinuxCommand("rm " + BASEPATH + "config/nginx/" + MapKey + ".conf")

					//重启nginx配置
					class.ExecLinuxCommand("cd " + BASEPATH + " && docker-compose exec nginx nginx -s reload")

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
					class.ExecLinuxCommand("rm -rf " + BASEPATH + MapKey + " && docker-compose exec mysql bash -c \"mysql -uroot -p" + MysqlPassword.(string) + " -e 'drop database " + MapKey + "'\"")

					//删除对应的nginx配置
					class.ExecLinuxCommand("rm " + BASEPATH + "config/nginx/" + MapKey + ".conf")

					//重启nginx配置
					class.ExecLinuxCommand("cd " + BASEPATH + " && docker-compose exec nginx nginx -s reload")

					fmt.Println("删除成功")

				}

			}

		case "备份管理":
			fmt.Println("备份管理")
			goto ServiceSelectFlag
		case "权限修复":
			fmt.Println("权限修复")
			goto ServiceSelectFlag
		case "升级系统镜像":
			fmt.Println("正在升级系统环境，预计需要5-15分钟.....")
			class.ExecLinuxCommand("cd " + BASEPATH + " && " + "docker-compose pull" + " && " + "docker-compose restart")
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
			class.ExecLinuxCommand("mkdir " + BASEPATH + "config/ && mkdir " + BASEPATH + "config/cert/ && mkdir " + BASEPATH + "config/mysql/ && mkdir " + BASEPATH + "config/nginx/ && mkdir " + BASEPATH + "config/php/")

			//设置代码目录为 10000,10000
			class.ExecLinuxCommand("chown -R 10000:10000 " + BASEPATH + "code/")

			//自动创建yaml标准模版

			// fmt.Println("安装过程")

			DockerComposeVersion := Template.DockerComposeVersion()
			DockerComposeNginxMap := class.YamlFileToMap(Template.DockerComposeNginx())
			DockerComposeMysqlString := Template.DockerComposeMysql()

			//自动生成mysql密码
			mysqlRandPassword := class.RandomString(16)
			//替换compose中的密码为随机密码
			DockerComposeMysqlString = strings.Replace(DockerComposeMysqlString, "MYSQL_ROOT_PASSWORD: root", "MYSQL_ROOT_PASSWORD: "+mysqlRandPassword, 1)
			DockerComposeMysqlMap := class.YamlFileToMap(DockerComposeMysqlString)

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

			//创建mysql my.cnf
			class.WriteFile(BASEPATH+"config/mysql/my.cnf", Template.MysqlCnf())

			//启动docker-compose
			fmt.Println("服务正在启动中，预计需要10分钟，请您耐心等待......")
			class.ExecLinuxCommand("cd " + BASEPATH + " && docker-compose up -d")

			//回显数据库密码
			fmt.Println("\n====================")
			fmt.Println("mysql 数据库信息")
			fmt.Println("用户名：root")
			fmt.Println("密码  " + mysqlRandPassword)
			fmt.Println("====================")

			goto CreateNewSiteFlag
		}
	}
}
