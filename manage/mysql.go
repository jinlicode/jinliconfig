package manage

import (
	"fmt"
	"jinliconfig/class"
)

// PhpMyAdminManage 管理phpmyadmin面板
func PhpMyAdminManage(basepath string) bool {
reSelectPhpMyAdmin:
	PhpMyadminSelect := class.ConsoleOptionsSelect("phpmyadmin面板", []string{"开启", "查看配置", "退出", "返回上一层"}, "请输入选项")
	switch PhpMyadminSelect {
	case "开启":
		class.ExecLinuxCommand("cd " + basepath + " && docker-compose up -d phpmyadmin")
		fmt.Println("phpmyadmin管理地址：您的服务器IP地址:8080")
		fmt.Println("mysql用户名：root")
		fmt.Println("mysql密码：" + class.ReadMysqlRootPassword(basepath))
		goto reSelectPhpMyAdmin
	case "查看配置":
		fmt.Println("phpmyadmin管理地址：您的服务器IP地址:8080")
		fmt.Println("mysql用户名：root")
		fmt.Println("mysql密码：" + class.ReadMysqlRootPassword(basepath))
		goto reSelectPhpMyAdmin
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
