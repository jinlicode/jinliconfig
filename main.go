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

	//命令行备份
	class.FlagBackupExec(BASEPATH)

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
		// fmt.Println(ExistSiteSlice)
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
				class.CreateSite(BASEPATH, DockerComposeYamlMap, SiteNetMax)

			case WebServiceSelect:
				if class.SiteManage(BASEPATH, WebServiceSelect, DockerComposeYamlMap) == false {
					goto WebServiceSelectFlag
				}
			}

		case "备份管理":

			if class.BackupSiteManage(BASEPATH, ExistSiteSlice) == false {
				goto ServiceSelectFlag
			}

		case "权限修复":
			class.ExecLinuxCommand("cd " + BASEPATH + "code && chown -R 10000:10000 *")
			fmt.Println("权限修复成功")
		case "升级系统镜像":
			fmt.Println("正在升级系统环境，预计需要5-15分钟.....")
			class.ExecLinuxCommand("cd " + BASEPATH + " && " + "docker-compose pull" + " && " + "docker-compose restart")
			// goto ServiceSelectFlag
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

			//创建备份目录
			class.ExecLinuxCommand("mkdir " + BASEPATH + "backup/")
			class.ExecLinuxCommand("mkdir " + BASEPATH + "backup/database")
			class.ExecLinuxCommand("mkdir " + BASEPATH + "backup/site")

			//创建自动备份目录
			class.ExecLinuxCommand("mkdir " + BASEPATH + "autobackup/")
			class.ExecLinuxCommand("mkdir " + BASEPATH + "autobackup/database")
			class.ExecLinuxCommand("mkdir " + BASEPATH + "autobackup/site")

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
			fmt.Println("\n=======================数据库ROOT信息===========================")
			fmt.Println("。数据库服务器地址：mysql")
			fmt.Println("。数据库用户名：root")
			fmt.Println("。数据库密码：" + mysqlRandPassword)
			fmt.Println("================================================================")

			goto CreateNewSiteFlag
		}
	}
}
