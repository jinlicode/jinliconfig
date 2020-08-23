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

//一维map转换切片 string类型
func mapToSlice(m map[string]string) []string {
	s := make([]string, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

// map 转yaml 类型
func MapToYaml(m map[string]interface{}) (string, error) {

	//先转成json

	jsonByte, err := json.Marshal(m)
	if err != nil {
		return "map to json error", nil
	}

	y, err := yaml.JSONToYAML(jsonByte)
	if err != nil {
		return "json to yaml error", nil
	}
	//最后转成yaml
	return string(y), nil
}

func JSONToYaml(json string) (string, error) {

	j := []byte(json)
	y, err := yaml.JSONToYAML(j)
	if err != nil {
		return "json to yaml error", nil
	}
	return string(y), nil
}
