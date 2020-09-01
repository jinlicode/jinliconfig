package main

import (
	"fmt"
	"io/ioutil"
	"jinliconfig/class"
	"jinliconfig/gotop"
	"jinliconfig/manage"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// https网站创建注意事项
// 1. 必须输入邮箱
// 2. 创建完成启动后，必须在基础目录执行: docker-compose exec nginx certbot -n --nginx --agree-tos -m example@example.com --domains www.example.com

func main() {

	//配置常量
	const (
		//基础目录配置
		BASEPATH = "/var/jinli/"
		JINLIVER = 1.3
	)
	//检查是否为root启动
	if os.Getuid() != 0 {
		fmt.Println("您的用户权限太低，请使用root用户执行，命令为：sudo su")
		os.Exit(3)
	}

	//命令行备份
	manage.FlagBackupExec(BASEPATH)

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
																			  
						欢迎使用锦鲤网站管理系统 v1.3
		`
	fmt.Println(Welcome)

	//检测版本更新 3秒
	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	seedURL := "https://api.github.com/repos/jinlicode/jinliconfig/releases"
	resp, err := client.Get(seedURL)
	if err == nil {

		body, error := ioutil.ReadAll(resp.Body)

		if error == nil {

			resTag := gjson.Get(string(body), "#.tag_name")

			resTagSlice := resTag.Array()
			jinliVersion := resTagSlice[0].String()

			jinliVersionString := strings.Replace(jinliVersion, "v", "", -1)

			//判断版本
			if class.CompareVersion(strconv.FormatFloat(JINLIVER, 'f', -1, 64), jinliVersionString) == -1 {
				//程序升级提示需
				fmt.Println("程序正在升级中，请勿退出程序......")
				class.ExecLinuxCommand("wget -O /tmp/jinliconfig https://release.jinli.plus/linux/x86_64/jinliconfig && mv /usr/sbin/jinliconfig /usr/sbin/jinliconfig_" + strconv.FormatFloat(JINLIVER, 'f', -1, 64) + " && mv /tmp/jinliconfig /usr/sbin/jinliconfig && chmod +x /usr/sbin/jinliconfig")
				fmt.Println("程序升级成功，请重新运行")
				os.Exit(1)

			}
		}

	}

	//检测安装的docker是否符合运行环境要求
	DockerClaim := class.ChkDokcerInstall()
	//如果不符合提示重装
	if DockerClaim != true {
		DockerReInstall := class.ConsoleUserConfirm("由于检测到您的系统不符合使用版本，需要对您的系统重新安装基础运行环境\n1. 会卸载docker\n2.重新安装符合版本的docker\n3.dockers原始文件目录/etc/docker将被删除\n4.docker原始目录/var/lib/docker将被删除\n 请问您是否同意重新安装docker？")
		if DockerReInstall {
			//删除老的docker 重新安装特定的docker
			class.ChkDokcerRemove()
			class.ExecDockerInstall()
			fmt.Println("安装基础环境完成。")
		} else {
			os.Exit(3)
		}
	}

CreateNewSiteFlag:

	if class.CheckFileExist(BASEPATH + "docker-compose.yaml") {

		//读取docker-compose配置文件
		DockerComposeYamlRead := class.ReadFile(BASEPATH + "docker-compose.yaml")
		DockerComposeYamlMap := class.YamlFileToMap(DockerComposeYamlRead)

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
		ServiceSelect := class.ConsoleOptionsSelect("请选择您需要的服务", []string{"网站服务", "备份管理", "数据库管理", "权限修复", "升级系统镜像", "服务器监控", "Mysql服务", "退出"}, "请输入选项")
		switch ServiceSelect {
		case "网站服务":
			//网站服务选择主菜单
		WebServiceSelectFlag:

			//获取已经存在的网站
			ExistSiteSlice := []string{}
			ExistStopSiteSlice := []string{}

			ExistSiteSlice = class.GetPathFiles(BASEPATH+"config/nginx/", false)
			if class.CheckFileExist(BASEPATH+"config/nginx_stop/") == false {
				class.ExecLinuxCommand("mkdir " + BASEPATH + "config/nginx_stop/")
			}
			ExistStopSiteSlice = class.GetPathFiles(BASEPATH+"config/nginx_stop/", false)

			for k, v := range ExistSiteSlice {
				ExistSiteSlice[k] = strings.Replace(v, ".conf", "", -1)
				ExistSiteSlice[k] = strings.Replace(ExistSiteSlice[k], "_", ".", -1)
			}

			for k, v := range ExistStopSiteSlice {
				ExistStopSiteSlice[k] = strings.Replace(v, ".conf", "", -1)
				ExistStopSiteSlice[k] = strings.Replace(ExistStopSiteSlice[k], "_", ".", -1) + "（已暂停）"
			}

			WebServiceSelectOption := []string{}
			WebServiceSelectOption = append(ExistSiteSlice, ExistStopSiteSlice...)
			WebServiceSelectOption = append(WebServiceSelectOption, "新增网站", "返回上层")

			WebServiceSelect := class.ConsoleOptionsSelect("请选择您需要管理的网站", WebServiceSelectOption, "请输入选项")
			switch WebServiceSelect {
			case "返回上层":
				fmt.Println("返回上层")
				goto ServiceSelectFlag
			case "新增网站":
				manage.CreateSite(BASEPATH, DockerComposeYamlMap)

			case WebServiceSelect:
				if WebServiceSelect == "interrupt" {
					fmt.Println("您已强制退出")
					os.Exit(1)
				}
				if manage.SiteManage(BASEPATH, WebServiceSelect, DockerComposeYamlMap) == false {
					goto WebServiceSelectFlag
				}
			}

		case "备份管理":

			//获取已经存在的网站
			ExistSiteSlice := []string{}
			ExistSiteSlice = class.GetPathFiles(BASEPATH+"config/nginx/", false)
			for k, v := range ExistSiteSlice {
				ExistSiteSlice[k] = strings.Replace(v, ".conf", "", -1)
				ExistSiteSlice[k] = strings.Replace(ExistSiteSlice[k], "_", ".", -1)
			}

			if manage.BackupSiteManage(BASEPATH, ExistSiteSlice) == false {
				goto ServiceSelectFlag
			}

		case "数据库管理":
			if manage.PhpMyAdminManage(BASEPATH) == false {
				goto ServiceSelectFlag
			}

		case "权限修复":
			class.ExecLinuxCommand("cd " + BASEPATH + "code && chown -R 10000:10000 *")
			fmt.Println("权限修复成功")
			goto ServiceSelectFlag
		case "升级系统镜像":
			fmt.Println("正在升级系统环境，预计需要5-15分钟.....")
			class.ExecLinuxCommand("cd " + BASEPATH + " && " + "docker-compose pull" + " && " + "docker-compose restart")
			// goto ServiceSelectFlag
		case "服务器监控":
			WebServiceTopSelect := class.ConsoleOptionsSelect("请选择您查看的监控服务", []string{"系统整体监控", "服务监控", "返回上层"}, "请选择选项")

			switch WebServiceTopSelect {
			case "系统整体监控":
				gotop.GetGoTop()
			case "服务监控":
				class.ExecLinuxCommand("docker stats --format \"table {{.Name}}\t{{.CPUPerc}}  {{.MemUsage}}\t{{.NetIO}}\"")
			case "返回上层":
				goto ServiceSelectFlag

			}
		case "Mysql服务":
			MysqlServerSelect := class.ConsoleOptionsSelect("请选择mysql服务", []string{"重启", "暂停", "返回上层"}, "请选择选项")

			switch MysqlServerSelect {
			case "重启":
				class.ExecLinuxCommand("cd " + BASEPATH + " && " + "docker-compose up -d mysql")
				fmt.Println("重启成功")
				goto ServiceSelectFlag

			case "暂停":
				class.ExecLinuxCommand("cd " + BASEPATH + " && " + "docker-compose stop mysql")
				fmt.Println("暂停成功")
				goto ServiceSelectFlag

			case "返回上层":
				goto ServiceSelectFlag

			}
		case "退出":
			fmt.Println("退出")
			break //可以添加
			// goto WebServiceSelect //随便添加的一会删除
		}
	} else {
		if manage.InstallJinliCode(BASEPATH, JINLIVER) == false {
			goto CreateNewSiteFlag
		}
	}
}
