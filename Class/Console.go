package Class

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ghodss/yaml"
)

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

//mysql数据库查询
// func MysqlGetDatabases(Mysql []string) []string {

// }

//mysql数据库信息获取
func MysqlInfo(YamlFile string) map[string]interface{} {
	println(YamlFile)
	DockerComposeJson, _ := yaml.YAMLToJSON([]byte(YamlFile))
	var m map[string]interface{}
	json.Unmarshal([]byte(DockerComposeJson), &m)
	return m
}

// Convert json string to map
func JsonToMap(jsonStr string) (map[string]string, error) {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return nil, err
	}

	for k, v := range m {
		fmt.Printf("%v: %v\n", k, v)
	}

	return m, nil
}

// Convert map json string
func MapToJson(m map[string]string) (string, error) {
	jsonByte, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Marshal with error: %+v\n", err)
		return "", nil
	}
	return string(jsonByte), nil
}
