package main

import (
	"fmt"
	"jinlicode/Class"
	"jinlicode/Template"
)

func main() {
	ddd := Class.UserInput("请输入您的域名")
	ccc := Class.OptionsSelect("run", []string{"全新安装", "重新安装（不删除数据）", "重新安装（删除数据）", "管理服务"}, "请输入选项")
	if ddd == "sss" {
		print("hhh")
	}
	fmt.Println(ccc)

	fmt.Println(Template.TemplatePhpWww())
	fmt.Println(Template.TemplateMysql())
	fmt.Println(Template.TemplatePhpIni())
	fmt.Println(Template.TemplatePhpFpm())
	fmt.Println(Template.TemplateCaddy())
}
