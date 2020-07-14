package main

import (
	"Class",
	"github.com/AlecAivazis/survey/v2"
)

//抽象出用户选择器需要输入三个函数
//msg:字符串，用来提示信息
//Opt：输入数组用来选择
//help：字符串，用来写帮助信息
func OptionsSelect(msg string, Opt []string, help string) string {
	color := ""
	prompt := &survey.Select{
		Message: msg,
		Options: Opt,
		Help:    help,
	}
	aaa := survey.AskOne(prompt, &color)
	if aaa != nil {
		return aaa.Error()
	} else {
		return color
	}
}

//抽象出用户输入，需要传入一个string字符串，用来提示用户
func UserInput(msg string) string {
	name := ""
	prompt := &survey.Input{
		Message: msg,
	}
	bbb := survey.AskOne(prompt, &name)
	if bbb != nil {
		return bbb.Error()
	} else {
		return name
	}
}

func main() {
	ddd := UserInput("请输入您的域名")
	ccc := OptionsSelect("run", []string{"全新安装", "重新安装（不删除数据）", "重新安装（删除数据）", "管理服务"}, "请输入选项")
	if ddd == "sss" {
		print("hhh")
	}
	print(ccc)
}
