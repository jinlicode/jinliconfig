package class

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
)

// 转换Yaml文件为Map
func YamlFileToMap(YamlFile string) map[string]interface{} {
	DockerComposeJson, _ := yaml.YAMLToJSON([]byte(YamlFile))
	var m map[string]interface{}
	json.Unmarshal([]byte(DockerComposeJson), &m)
	return m
}

// 转换Map为json文件
func MapToJson(m map[string]interface{}) (string, error) {
	jsonByte, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Marshal with error: %+v\n", err)
		return "", nil
	}
	return string(jsonByte), nil
}

func mapToSlice(m map[string]string) []string {
	s := make([]string, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}
