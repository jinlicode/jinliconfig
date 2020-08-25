package manage

import (
	"fmt"
	"jinliconfig/Template"
	"jinliconfig/class"
	"os"
	"strconv"
	"strings"
)

//CreateSite 创建新网站
func CreateSite(basepath string, DockerComposeYamlMap map[string]interface{}, SiteNetMax int) {
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

	//检测是否还有项目目录存在
	if class.CheckFileExist(basepath + "code/" + newDomain) {
		if class.ConsoleOptionsSelect("已存在"+newDomain+"目录，是否继续？继续将删除此目录和以此目录为名的数据库", []string{"是", "否"}, "请输入选项") == "否" {
			goto ReInputSiteDomainFlag
		} else {
			//删除对应的项目目录和数据库
			RootPassword := class.ReadMysqlRootPassword(basepath)
			class.ExecLinuxCommand("cd " + basepath + " && rm -rf " + basepath + "code/" + newDomain + " && docker-compose exec mysql bash -c \"mysql -uroot -p" + RootPassword + " -e 'drop database " + newDomain + "'\"")

		}
	}

	//检测nginx是否已经存在配置文件
	if class.CheckFileExist(basepath + "config/nginx/" + newDomain + ".conf") {
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

	//自动创建网站对应mysql数据
	MysqlRootPassword := class.ReadMysqlRootPassword(basepath)
	MysqlRootPasswordString := MysqlRootPassword
	//获取随机密码
	mysqlSiteRandPassword := class.RandomString(16)

	//替换mysql信息到环境变量
	SitePhpVersionCompose = strings.Replace(SitePhpVersionCompose, "MYSQL_HOST=MYSQL_HOST", "MYSQL_HOST="+class.ReadMysqlHost(basepath), 1)
	SitePhpVersionCompose = strings.Replace(SitePhpVersionCompose, "MYSQL_USER=MYSQL_USER", "MYSQL_USER="+newDomain, 1)
	SitePhpVersionCompose = strings.Replace(SitePhpVersionCompose, "MYSQL_PASS=MYSQL_PASS", "MYSQL_PASS="+mysqlSiteRandPassword, 1)

	//生成子map
	NewSitePhpVersionComposeMap := class.YamlFileToMap(SitePhpVersionCompose)

	//插入到总的yaml文件中
	DockerComposeYamlMap["services"].(map[string]interface{})[newDomain] = NewSitePhpVersionComposeMap

	//自动创建以网站名字命名的程序目录
	class.ExecLinuxCommand("mkdir " + basepath + "code/" + newDomain)
	class.ExecLinuxCommand("mkdir " + basepath + "config/php/" + newDomain)

	//写入404以及index文件到置顶目录
	class.WriteFile(basepath+"code/"+newDomain+"/index.html", Template.HTMLIndex())
	class.WriteFile(basepath+"code/"+newDomain+"/404.html", Template.HTML404())

	//创建网站的配置文件到对应的config配置文件中
	class.ExecLinuxCommand("mkdir " + basepath + "config/php/" + newDomain)
	class.WriteFile(basepath+"config/php/"+newDomain+"/www.conf", Template.PhpWww())
	class.WriteFile(basepath+"config/php/"+newDomain+"/php.ini", Template.PhpIni())
	class.WriteFile(basepath+"config/php/"+newDomain+"/php-fpm.conf", Template.PhpFpm())

	//创建对应nginx.conf到对应目录
	if NewSiteHTTPS == "否" {
		TemplateNginxHTTPString := Template.TemplateNginxHttp()
		TemplateNginxHTTPString = strings.Replace(TemplateNginxHTTPString, "www_example_com", newDomain, -1)
		TemplateNginxHTTPString = strings.Replace(TemplateNginxHTTPString, "www.example.com", NewSiteDomain, -1)
		TemplateNginxHTTPString = strings.Replace(TemplateNginxHTTPString, "php:9000", newDomain+":9000", -1)
		class.WriteFile(basepath+"config/nginx/"+newDomain+".conf", TemplateNginxHTTPString)

	} else {

		TemplateNginxHTTPSString := Template.TemplateNginxHttps()
		TemplateNginxHTTPSString = strings.Replace(TemplateNginxHTTPSString, "www_example_com", newDomain, -1)
		TemplateNginxHTTPSString = strings.Replace(TemplateNginxHTTPSString, "www.example.com", NewSiteDomain, -1)
		TemplateNginxHTTPSString = strings.Replace(TemplateNginxHTTPSString, "php:9000", newDomain+":9000", -1)
		class.WriteFile(basepath+"config/nginx/"+newDomain+".conf", TemplateNginxHTTPSString)

	}

	//写入docker-compose.yaml 文件
	NewDockerComposeYamlString, _ := class.MapToYaml(DockerComposeYamlMap)
	class.WriteFile(basepath+"docker-compose.yaml", NewDockerComposeYamlString)

	//启动新网站服务
	class.ExecLinuxCommand("cd " + basepath + " && docker-compose up -d " + newDomain)
	//重启nginx 配置
	class.ExecLinuxCommand("cd " + basepath + " && docker-compose exec nginx nginx -s reload")

	//重启命令
	if NewSiteHTTPS == "是" {
		class.ExecLinuxCommand("cd " + basepath + " && docker-compose exec nginx certbot -n --nginx --agree-tos -m " + NewSiteSSLEmail + " --domains " + NewSiteDomain)
	}

	//自动创建数据库 用户名 密码
	class.CreateDatabase(basepath, MysqlRootPasswordString, newDomain, newDomain, mysqlSiteRandPassword)

	//显示新网站内容

	fmt.Println("\n=======================您的网站对应信息==========================")
	fmt.Println("。请将您的网站代码上传至 【" + basepath + "code/" + newDomain + "】 目录")
	fmt.Println("。数据库服务器地址：" + class.ReadMysqlHost(basepath))
	fmt.Println("。数据库库名：" + newDomain)
	fmt.Println("。数据库用户名：" + newDomain)
	fmt.Println("。数据库密码：" + mysqlSiteRandPassword)
	fmt.Println("================================================================")

	// fmt.Println(NewDockerComposeYamlString)
}

// SiteManage 网站管理
func SiteManage(basepath string, WebServiceSelect string, DockerComposeYamlMap map[string]interface{}) bool {
	WebConfigSelect := class.ConsoleOptionsSelect("请选择您需要管理的网站服务", []string{
		// WebServiceSelect + "的" + "nginx配置",
		// WebServiceSelect + "的" + "php配置",
		// WebServiceSelect + "的" + "数据库配置",
		"查看" + WebServiceSelect + "数据库信息",
		"重置" + WebServiceSelect + "数据库密码",
		"重启" + WebServiceSelect + "网站服务",
		"暂停" + WebServiceSelect + "网站服务",
		"删除" + WebServiceSelect + "的网站(不删除数据)",
		"删除" + WebServiceSelect + "的网站(删除数据，包含数据库和程序)",
		"返回上层"}, "请输入选项")

	MapKey := strings.Replace(WebServiceSelect, ".", "_", -1)

	//网站内服务修改主菜单
WebConfigSelectFlag:
	switch WebConfigSelect {
	case "返回上层":
		fmt.Println("返回上层")
		return false
	case "查看" + WebServiceSelect + "数据库信息":
		class.PrintHr()
		fmt.Println("数据库服务器地址：" + class.ReadMysqlHost(basepath))
		fmt.Println(WebServiceSelect + "的数据库用户名：" + class.ReadSiteMysqlInfo(basepath, MapKey, "user"))
		fmt.Println(WebServiceSelect + "的数据库密码：" + class.ReadSiteMysqlInfo(basepath, MapKey, "pass"))
		class.PrintHr()
		return false

	case "重置" + WebServiceSelect + "数据库密码":

		ReConfirm := class.ConsoleUserConfirm("确定重置" + WebServiceSelect + "数据库密码吗？")

		if ReConfirm != true {
			fmt.Println("已取消操作")
			return false
		}

		newPass := MysqlSiteEditPass(basepath, MapKey)
		if newPass != "" {
			class.PrintHr()
			fmt.Println(WebServiceSelect + "的新数据库密码为：" + newPass)
			class.PrintHr()

		}
		return false

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
	case "重启" + WebServiceSelect + "网站服务":
		//确定是否需要暂停
		ReStartSiteConfirm := class.ConsoleUserConfirm("确定重启" + WebServiceSelect + "服务吗？")
		if ReStartSiteConfirm != true {
			fmt.Println("已取消操作")
			return false
		}

		class.ExecLinuxCommand("cd " + basepath + " && docker-compose restart " + strings.Replace(WebServiceSelect, ".", "_", -1))
		fmt.Println("重启成功")

	case "暂停" + WebServiceSelect + "网站服务":
		//确定是否需要暂停
		ReStopSiteConfirm := class.ConsoleUserConfirm("确定暂停" + WebServiceSelect + "服务吗？")
		if ReStopSiteConfirm != true {
			fmt.Println("已取消操作")
			return false
		}

		//输入命令 暂停容器
		// fmt.Println("cd " + basepath + " && docker-compose stop " + strings.Replace(WebServiceSelect, ".", "_", -1))
		class.ExecLinuxCommand("cd " + basepath + " && docker-compose stop " + strings.Replace(WebServiceSelect, ".", "_", -1))
		fmt.Println("暂停成功")
		goto WebConfigSelectFlag

	case "删除" + WebServiceSelect + "的网站(不删除数据)":
		//确定是否需要删除
		ReDelSiteConfirm := class.ConsoleUserConfirm("确定删除" + WebServiceSelect + "网站吗？(不删除数据)")
		if ReDelSiteConfirm != true {
			fmt.Println("已取消操作")
			return false
		}
		//输入命令 删除yaml中服务
		class.ExecLinuxCommand("cd " + basepath + " && docker-compose stop " + MapKey + " && docker-compose rm " + MapKey)

		//执行完之后删除yaml中对应的map
		delete(DockerComposeYamlMap["services"].(map[string]interface{}), MapKey)

		//重新写入到yaml
		NewDockerComposeYamlString, _ := class.MapToYaml(DockerComposeYamlMap)
		class.WriteFile(basepath+"docker-compose.yaml", NewDockerComposeYamlString)

		//删除对应的nginx配置
		class.ExecLinuxCommand("rm " + basepath + "config/nginx/" + MapKey + ".conf")

		//重启nginx配置
		class.ExecLinuxCommand("cd " + basepath + " && docker-compose exec nginx nginx -s reload")

		fmt.Println("删除成功")
		return false

	case "删除" + WebServiceSelect + "的网站(删除数据，包含数据库和程序)":
		//确定是否需要删除
		ReDelSiteConfirm := class.ConsoleUserConfirm("确定删除" + WebServiceSelect + "网站吗？(删除数据，包含数据库和程序)")
		if ReDelSiteConfirm != true {
			fmt.Println("已取消操作")
			return false
		}
		//输入命令 删除yaml中服务
		MapKey := strings.Replace(WebServiceSelect, ".", "_", -1)
		// fmt.Println("cd " + basepath + " && docker-compose stop " + MapKey + " && docker-compose rm " + MapKey)
		class.ExecLinuxCommand("cd " + basepath + " && docker-compose stop " + MapKey + " && docker-compose rm " + MapKey)

		//执行完之后删除yaml中对应的map
		delete(DockerComposeYamlMap["services"].(map[string]interface{}), MapKey)

		//重新写入到yaml
		NewDockerComposeYamlString, _ := class.MapToYaml(DockerComposeYamlMap)
		class.WriteFile(basepath+"docker-compose.yaml", NewDockerComposeYamlString)

		//操作删除工作 删除代码目录  删除  数据库 drop database bbbbbbbb
		MysqlPassword := class.ReadMysqlRootPassword(basepath)
		class.ExecLinuxCommand("cd " + basepath + " && rm -rf " + basepath + "code/" + MapKey + " && docker-compose exec mysql bash -c \"mysql -uroot -p" + MysqlPassword + " -e 'drop database " + MapKey + "'\"")

		//删除对应的nginx配置
		class.ExecLinuxCommand("rm " + basepath + "config/nginx/" + MapKey + ".conf")

		//重启nginx配置
		class.ExecLinuxCommand("cd " + basepath + " && docker-compose exec nginx nginx -s reload")

		fmt.Println("删除成功")
		return false

	}

	return true
}
