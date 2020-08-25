package manage

import (
	"fmt"
	"jinliconfig/class"
	"regexp"
)

// PhpMyAdminManage 管理phpmyadmin面板
func PhpMyAdminManage(basepath string) bool {
reSelectPhpMyAdmin:
	PhpMyadminSelect := class.ConsoleOptionsSelect("phpmyadmin面板", []string{"开启", "查看配置", "重置root密码", "退出", "返回上一层"}, "请输入选项")
	switch PhpMyadminSelect {
	case "开启":
		class.ExecLinuxCommand("cd " + basepath + " && docker-compose up -d phpmyadmin")
		class.PrintHr()
		fmt.Println("phpmyadmin管理地址：您的服务器IP地址:8080")
		fmt.Println("mysql用户名：root")
		fmt.Println("mysql密码：" + class.ReadMysqlRootPassword(basepath))
		class.PrintHr()
		goto reSelectPhpMyAdmin
	case "查看配置":
		class.PrintHr()
		fmt.Println("phpmyadmin管理地址：您的服务器IP地址:8080")
		fmt.Println("mysql用户名：root")
		fmt.Println("mysql密码：" + class.ReadMysqlRootPassword(basepath))
		class.PrintHr()
		goto reSelectPhpMyAdmin
	case "重置root密码":

		ReConfirm := class.ConsoleUserConfirm("确定重置root密码吗？")

		if ReConfirm != true {
			fmt.Println("已取消操作")
			return false
		}
		
		newPass := MysqlRootEditPass(basepath)
		if newPass != "" {
			class.PrintHr()
			fmt.Println("root新密码密码：" + newPass)
			class.PrintHr()
		}
		return false

	case "退出":
		class.ExecLinuxCommand("cd " + basepath + " && docker-compose stop phpmyadmin")
	case "返回上一层":
		class.ExecLinuxCommand("cd " + basepath + " && docker-compose stop phpmyadmin")
		return false
	default:
		class.ExecLinuxCommand("cd " + basepath + " && docker-compose stop phpmyadmin")
	}

	return true
}

// MysqlSiteEditPass 网站mysql密码服务
func MysqlSiteEditPass(basepath string, newDomain string) string {

	if class.CheckFileExist(basepath + "docker-compose.yaml") {

		//读取docker-compose配置文件
		DockerComposeYamlRead := class.ReadFile(basepath + "docker-compose.yaml")
		DockerComposeYamlMap := class.YamlFileToMap(DockerComposeYamlRead)
		newDomainMap := DockerComposeYamlMap["services"].(map[string]interface{})[newDomain]

		//获取mysqlroot信息
		MysqlRootPassword := class.ReadMysqlRootPassword(basepath)
		MysqlRootPasswordString := MysqlRootPassword
		MysqlRootHOST := class.ReadMysqlHost(basepath)

		//获取随机密码
		mysqlSiteRandPassword := class.RandomString(16)

		//转成json 再做替换
		newDomainMapJSONString, _ := class.MapToJson(newDomainMap.(map[string]interface{}))
		// SitePhpVersionCompose = strings.Replace(SitePhpVersionCompose, "MYSQL_PASS=MYSQL_PASS", "MYSQL_PASS="+mysqlSiteRandPassword, 1)

		parttern := `(?U)"MYSQL_PASS=(.*)"`
		re, _ := regexp.Compile(parttern)
		//将匹配到的部分替换为"##.#"
		newDomainMapJSONString = re.ReplaceAllString(newDomainMapJSONString, `"MYSQL_PASS=`+mysqlSiteRandPassword+`"`)

		class.MysqlQuery(MysqlRootHOST, "root", MysqlRootPasswordString, "mysql", `set password for '`+newDomain+`'@'%' = password('`+mysqlSiteRandPassword+`');`)
		class.MysqlQuery(MysqlRootHOST, "root", MysqlRootPasswordString, "mysql", "flush privileges")

		//写入yaml
		NewSitePhpVersionComposeYaml, _ := class.JSONToYaml(newDomainMapJSONString)

		//yam to map
		NewSitePhpVersionComposeMap := class.YamlFileToMap(NewSitePhpVersionComposeYaml)

		//插入到总的yaml文件中
		DockerComposeYamlMap["services"].(map[string]interface{})[newDomain] = NewSitePhpVersionComposeMap

		NewDockerComposeYamlString, _ := class.MapToYaml(DockerComposeYamlMap)
		//写入新的yaml文件
		class.WriteFile(basepath+"docker-compose.yaml", NewDockerComposeYamlString)

		return mysqlSiteRandPassword
	}
	return ""
}

// MysqlRootEditPass mysql密码服务
func MysqlRootEditPass(basepath string) string {

	if class.CheckFileExist(basepath + "docker-compose.yaml") {

		//读取docker-compose配置文件
		DockerComposeYamlRead := class.ReadFile(basepath + "docker-compose.yaml")
		DockerComposeYamlMap := class.YamlFileToMap(DockerComposeYamlRead)

		//获取mysqlroot信息
		MysqlRootPassword := class.ReadMysqlRootPassword(basepath)
		MysqlRootPasswordString := MysqlRootPassword
		MysqlRootHOST := class.ReadMysqlHost(basepath)

		//获取随机密码
		mysqlRandPassword := class.RandomString(16)

		class.MysqlQuery(MysqlRootHOST, "root", MysqlRootPasswordString, "mysql", `set password for 'root'@'%' = password('`+mysqlRandPassword+`');`)
		class.MysqlQuery(MysqlRootHOST, "root", mysqlRandPassword, "mysql", "flush privileges")

		DockerComposeYamlMap["services"].(map[string]interface{})["mysql"].(map[string]interface{})["environment"].(map[string]interface{})["MYSQL_ROOT_PASSWORD"] = mysqlRandPassword

		NewDockerComposeYamlString, _ := class.MapToYaml(DockerComposeYamlMap)
		//写入新的yaml文件
		class.WriteFile(basepath+"docker-compose.yaml", NewDockerComposeYamlString)

		return mysqlRandPassword

	}
	return ""
}
