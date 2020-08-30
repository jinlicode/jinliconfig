package manage

import (
	"fmt"
	"jinliconfig/Template"
	"jinliconfig/class"
	"os"
	"strconv"
	"strings"
)

// InstallJinliCode 首次安装jinlicode
func InstallJinliCode(basepath string, jinliVersion float32) bool {
	fmt.Println("您初始化锦鲤部署，是否要初始化？")
	NewInstall := class.ConsoleOptionsSelect("请输入您的选项", []string{"否", "是"}, "请选择是否重新初始化")
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
		class.ExecLinuxCommand("mkdir " + basepath)

		//创建代码目录
		class.ExecLinuxCommand("mkdir " + basepath + "code/")

		//创建各配置项目录
		class.ExecLinuxCommand("mkdir " + basepath + "config/")
		class.ExecLinuxCommand("mkdir " + basepath + "config/cert/")
		class.ExecLinuxCommand("mkdir " + basepath + "config/mysql/")
		class.ExecLinuxCommand("mkdir " + basepath + "config/nginx/")
		class.ExecLinuxCommand("mkdir " + basepath + "config/php/")
		class.ExecLinuxCommand("mkdir " + basepath + "config/rewrite/")
		// class.ExecLinuxCommand("mkdir " + basepath + "config/ && mkdir " + basepath + "config/cert/ && mkdir " + basepath + "config/mysql/ && mkdir " + basepath + "config/nginx/ && mkdir " + basepath + "config/php/")

		//创建nginx网址停止目录
		class.ExecLinuxCommand("mkdir " + basepath + "config/nginx_stop/")

		//创建备份目录
		class.ExecLinuxCommand("mkdir " + basepath + "backup/")
		class.ExecLinuxCommand("mkdir " + basepath + "backup/database")
		class.ExecLinuxCommand("mkdir " + basepath + "backup/site")

		//创建自动备份目录
		class.ExecLinuxCommand("mkdir " + basepath + "autobackup/")
		class.ExecLinuxCommand("mkdir " + basepath + "autobackup/database")
		class.ExecLinuxCommand("mkdir " + basepath + "autobackup/site")

		//设置代码目录为 10000,10000
		class.ExecLinuxCommand("chown -R 10000:10000 " + basepath + "code/")

		//自动创建yaml标准模版

		// fmt.Println("安装过程")

		DockerComposeVersion := Template.DockerComposeVersion()
		DockerComposeNginxString := Template.DockerComposeNginx()
		DockerComposeMysqlString := Template.DockerComposeMysql()

		//替换code版本号
		jinliVersionString := strconv.FormatFloat(float64(jinliVersion), 'f', 2, 64)
		DockerComposeNginxString = strings.Replace(DockerComposeNginxString, "JINLIVER=VERSION", "JINLIVER="+jinliVersionString, 1)

		DockerComposeNginxMap := class.YamlFileToMap(DockerComposeNginxString)

		//自动生成mysql密码
		mysqlRandPassword := class.RandomString(16)
		//替换compose中的密码为随机密码
		DockerComposeMysqlString = strings.Replace(DockerComposeMysqlString, "MYSQL_ROOT_PASSWORD: root", "MYSQL_ROOT_PASSWORD: "+mysqlRandPassword, 1)
		DockerComposeMysqlMap := class.YamlFileToMap(DockerComposeMysqlString)

		DockerComposeNetWorksMap := class.YamlFileToMap(Template.DockerComposeNetWorks())

		//获取phpmyadmin
		DockerComposePhpmyadminMap := class.YamlFileToMap(Template.DockerComposePhpmyadmin())

		DockerComposeMap := make(map[string]interface{})
		DockerComposeMap["version"] = DockerComposeVersion
		DockerComposeMap["services"] = make(map[string]interface{})
		DockerComposeMap["services"].(map[string]interface{})["nginx"] = DockerComposeNginxMap
		DockerComposeMap["services"].(map[string]interface{})["mysql"] = DockerComposeMysqlMap
		DockerComposeMap["services"].(map[string]interface{})["phpmyadmin"] = DockerComposePhpmyadminMap
		DockerComposeMap["networks"] = DockerComposeNetWorksMap

		//写入yaml文件 跳转到新建网站
		DockerComposeYamlString, _ := class.MapToYaml(DockerComposeMap)
		class.WriteFile(basepath+"docker-compose.yaml", DockerComposeYamlString)

		//创建mysql my.cnf
		class.WriteFile(basepath+"config/mysql/my.cnf", Template.MysqlCnf())

		//启动docker-compose
		fmt.Println("服务正在启动中，预计需要10分钟，请您耐心等待......")
		class.ExecLinuxCommand("cd " + basepath + " && docker-compose up -d && docker-compose stop phpmyadmin")

		//创建拉取任务
		fmt.Println("创建后台拉取镜像任务成功")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/mysql:latest > /dev/null 2>&1 & ")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/php:v5.6-sec > /dev/null 2>&1 & ")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/php:v5.6 > /dev/null 2>&1 & ")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/php:v7.0-sec > /dev/null 2>&1 & ")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/php:v7.0 > /dev/null 2>&1 & ")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/php:v7.1-sec > /dev/null 2>&1 & ")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/php:v7.1 > /dev/null 2>&1 & ")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/php:v7.2-sec > /dev/null 2>&1 & ")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/php:v7.2 > /dev/null 2>&1 & ")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/php:v7.3-sec > /dev/null 2>&1 & ")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/php:v7.3 > /dev/null 2>&1 & ")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/nginx:v1 > /dev/null 2>&1 & ")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/memcached:1.6.6 > /dev/null 2>&1 & ")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/redis:5.0.9 > /dev/null 2>&1 & ")
		class.ExecLinuxCommand("nohup docker pull registry.cn-beijing.aliyuncs.com/jinlicode/phpmyadmin:5.0.2 > /dev/null 2>&1 & ")

		//回显数据库密码
		fmt.Println("\n=======================数据库ROOT信息===========================")
		fmt.Println("。数据库服务器地址:" + class.ReadMysqlHost(basepath))
		fmt.Println("。数据库用户名:root")
		fmt.Println("。数据库密码:" + mysqlRandPassword)
		fmt.Println("================================================================")

		return false
	}
	return false
}
