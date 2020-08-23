package manage

import (
	"fmt"
	"jinliconfig/Template"
	"jinliconfig/class"
	"os"
	"strings"
)

// InstallJinliCode 首次安装jinlicode
func InstallJinliCode(basepath string) bool {
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
		class.ExecLinuxCommand("mkdir " + basepath)

		//创建代码目录
		class.ExecLinuxCommand("mkdir " + basepath + "code/")

		//创建各配置项目录
		class.ExecLinuxCommand("mkdir " + basepath + "config/ && mkdir " + basepath + "config/cert/ && mkdir " + basepath + "config/mysql/ && mkdir " + basepath + "config/nginx/ && mkdir " + basepath + "config/php/")

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
		class.WriteFile(basepath+"docker-compose.yaml", DockerComposeYamlString)

		//创建mysql my.cnf
		class.WriteFile(basepath+"config/mysql/my.cnf", Template.MysqlCnf())

		//启动docker-compose
		fmt.Println("服务正在启动中，预计需要10分钟，请您耐心等待......")
		class.ExecLinuxCommand("cd " + basepath + " && docker-compose up -d")

		//回显数据库密码
		fmt.Println("\n=======================数据库ROOT信息===========================")
		fmt.Println("。数据库服务器地址：" + class.ReadMysqlHost(basepath))
		fmt.Println("。数据库用户名：root")
		fmt.Println("。数据库密码：" + mysqlRandPassword)
		fmt.Println("================================================================")

		return false
	}
	return false
}
