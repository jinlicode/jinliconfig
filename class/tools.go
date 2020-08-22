package class

import (
	"math/rand"
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
