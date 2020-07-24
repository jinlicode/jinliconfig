package class

import "github.com/AlecAivazis/survey/v2"

/*
抽象出用户选择器需要输入三个函数,
msg:字符串，用来提示信息,
Opt：输入数组用来选择,
help：字符串，用来写帮助信息。
*/
func ConsoleOptionsSelect(msg string, Opt []string, help string) string {
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
func ConsoleUserInput(msg string) string {
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

//要求用户确认,返回布尔值
func ConsoleUserConfirm(msg string) bool {
	name := false
	prompt := &survey.Confirm{
		Message: msg,
	}
	survey.AskOne(prompt, &name)
	return name
}

//多行文本输入，一般是输入证书时候使用，返回string类型
func ConsoleUserText(msg string) string {
	text := ""
	prompt := &survey.Multiline{
		Message: msg,
	}
	survey.AskOne(prompt, &text)
	return text
}
