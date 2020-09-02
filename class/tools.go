package class

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// RandomString 返回随机字符串
func RandomString(n int, allowedChars ...[]rune) string {
	var letters []rune
	defaultLetters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// ReadMysqlRootPassword  读取compose文件中的mysql root密码
func ReadMysqlRootPassword(basepath string) string {
	//先检测是否存在yaml文件
	if CheckFileExist(basepath + "docker-compose.yaml") {
		DockerComposeYamlRead := ReadFile(basepath + "docker-compose.yaml")
		DockerComposeYamlMap := YamlFileToMap(DockerComposeYamlRead)
		MysqlPassword := DockerComposeYamlMap["services"].(map[string]interface{})["mysql"].(map[string]interface{})["environment"].(map[string]interface{})["MYSQL_ROOT_PASSWORD"]

		return MysqlPassword.(string)
	}
	return ""
}

// ReadMysqlHost  读取compose文件中的mysql host
func ReadMysqlHost(basepath string) string {
	//先检测是否存在yaml文件
	if CheckFileExist(basepath + "docker-compose.yaml") {
		DockerComposeYamlRead := ReadFile(basepath + "docker-compose.yaml")
		DockerComposeYamlMap := YamlFileToMap(DockerComposeYamlRead)
		MysqlHost := DockerComposeYamlMap["services"].(map[string]interface{})["mysql"].(map[string]interface{})["networks"].(map[string]interface{})["jinli_net"].(map[string]interface{})["ipv4_address"]

		return MysqlHost.(string)
	}
	return ""
}

// ReadSiteMysqlInfo  读取compose文件中的网站 mysql 信息
func ReadSiteMysqlInfo(basepath string, dockerName string, readType string) string {
	//先检测是否存在yaml文件
	if CheckFileExist(basepath + "docker-compose.yaml") {
		DockerComposeYamlRead := ReadFile(basepath + "docker-compose.yaml")
		DockerComposeYamlMap := YamlFileToMap(DockerComposeYamlRead)
		bb := DockerComposeYamlMap["services"].(map[string]interface{})[dockerName].(map[string]interface{})
		mapJSON, err := MapToJson(bb)
		if err == nil {
			if readType == "host" {
				parttern := `(?U)"MYSQL_HOST=(.*)"`
				r := regexp.MustCompile(parttern)
				matchs := r.FindStringSubmatch(mapJSON)
				if len(matchs) == 2 {
					return matchs[1]
				}
			} else if readType == "pass" {
				parttern := `(?U)"MYSQL_PASS=(.*)"`
				r := regexp.MustCompile(parttern)
				matchs := r.FindStringSubmatch(mapJSON)
				if len(matchs) == 2 {
					return matchs[1]
				}
			} else if readType == "user" {
				parttern := `(?U)"MYSQL_USER=(.*)"`
				r := regexp.MustCompile(parttern)
				matchs := r.FindStringSubmatch(mapJSON)
				if len(matchs) == 2 {
					return matchs[1]
				}
			}

			return ""
		}
	}
	return ""
}

// PrintHr 打印一行等号
func PrintHr() {
	fmt.Println("\n====================================")
}

// GetPathFiles 获取木下的所有文件切片
func GetPathFiles(path string, isDir bool) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	fs := []string{}
	for _, f := range files {
		if isDir == true {
			if f.IsDir() == true {
				fs = append(fs, f.Name())
			}
		} else {
			fs = append(fs, f.Name())
		}
	}
	return fs
}

//GetComposeSiteNetMap 获取网站所有内网最后ip
func GetComposeSiteNetMap(basepath string) map[int]int {
	DockerComposeYamlRead := ReadFile(basepath + "docker-compose.yaml")
	DockerComposeYamlMap := YamlFileToMap(DockerComposeYamlRead)

	netMap := make(map[int]int)
	for k, v := range DockerComposeYamlMap["services"].(map[string]interface{}) {
		if k != "nginx" && k != "mysql" && k != "phpmyadmin" {
			//获取内网数字
			netString := v.(map[string]interface{})["networks"].(map[string]interface{})["jinli_net"].(map[string]interface{})["ipv4_address"].(string)
			if strings.Index(netString, "10.99.2") != -1 {
				netNumSlice := strings.Split(netString, ".")
				maxNumString := netNumSlice[3]
				netNum, err := strconv.Atoi(maxNumString)
				if err == nil {
					netMap[netNum] = netNum
				}
			}
		}
	}
	return netMap
}

// SiteDotToUnderline 点转下划线
func SiteDotToUnderline(siteName string) string {
	return strings.Replace(siteName, ".", "_", -1)
}

//CheckDockerMapServiceExist 判断是否存在
func CheckDockerMapServiceExist(basepath string, siteName string) bool {

	DockerComposeYamlRead := ReadFile(basepath + "docker-compose.yaml")
	DockerComposeYamlMap := YamlFileToMap(DockerComposeYamlRead)

	_, err := DockerComposeYamlMap["services"].(map[string]interface{})[siteName]

	return err
}

//获取当前服务的内网ip
func GetComposeServerNetString(basepath string, serverName string, fullName bool) string {
	DockerComposeYamlRead := ReadFile(basepath + "docker-compose.yaml")
	DockerComposeYamlMap := YamlFileToMap(DockerComposeYamlRead)

	//获取内网数字
	netString := DockerComposeYamlMap["services"].(map[string]interface{})[serverName].(map[string]interface{})["networks"].(map[string]interface{})["jinli_net"].(map[string]interface{})["ipv4_address"].(string)
	if fullName == true {
		return netString

	}
	netNumSlice := strings.Split(netString, ".")
	return netNumSlice[3]
}

//PHPChooseVersion php版本选择
func PHPChooseVersion(back bool) string {

	phpVersion := []string{
		"5.6",
		"5.6-sec(安全版本)",
		"7.0",
		"7.0-sec(安全版本)",
		"7.1",
		"7.1-sec(安全版本)",
		"7.2",
		"7.2-sec(安全版本)",
		"7.3",
		"7.3-sec(安全版本)",
	}
	if back == true {
		phpVersion = append(phpVersion, "返回上层")
	}
	PhpVersion := ConsoleOptionsSelect("请选择您需要的php版本, sec版本为安全版本", phpVersion, "请输入选项")

	return strings.Replace(PhpVersion, "(安全版本)", "", -1)
}
