package class

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//mysql数据库信息获取
func MysqlInfo(YamlFile string) []string {
	//yaml转换成map
	YamlMap := YamlFileToMap(YamlFile)
	//map获取数据库密码
	MysqlPassword := YamlMap["services"].(map[string]interface{})["mysql"].(map[string]interface{})["environment"].(map[string]interface{})["MYSQL_ROOT_PASSWORD"]
	MysqlHost := YamlMap["services"].(map[string]interface{})["mysql"].(map[string]interface{})["networks"].(map[string]interface{})["discuz"].(map[string]interface{})["ipv4_address"]
	MysqlUser := "root"
	//获取数据库内容
	databases := MysqlGetDatabases(MysqlHost.(string), MysqlUser, MysqlPassword.(string))
	return databases
}

//在数据库执行sql语句
func MysqlQuery(MysqlHost string, MysqlUser string, MysqlPassword string, DatabaseName string, QuerySql string) {
	//数据库配置

	port := "3306"
	//Db数据库连接池
	var DB *sql.DB

	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{MysqlUser, ":", MysqlPassword, "@tcp(", MysqlHost, ":", port, ")/", DatabaseName, "?charset=utf8mb4"}, "")

	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	//验证连接
	if err := DB.Ping(); err != nil {
		fmt.Println("opon database fail")
	}
	//执行查询
	Databases, _ := DB.Query(QuerySql)
	defer Databases.Close()
}

//读取数据库内容
func MysqlGetDatabases(MysqlHost string, MysqlUser string, MysqlPassword string) []string {
	//数据库配置

	port := "3306"
	dbName := "discuz"
	//Db数据库连接池
	var DB *sql.DB

	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{MysqlUser, ":", MysqlPassword, "@tcp(", MysqlHost, ":", port, ")/", dbName, "?charset=utf8mb4"}, "")

	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	//验证连接
	if err := DB.Ping(); err != nil {
		fmt.Println("opon database fail")
	}
	//执行查询
	Databases, err := DB.Query("show databases;")
	if err != nil {
		fmt.Println("查询失败，请检查docker-compose.yaml")
	}
	//循环出所有数据库名称，压入map
	m2 := make(map[string]string)
	var table string
	for Databases.Next() {
		Databases.Scan(&table)
		m2[table] = table
	}
	//从map里删除数据库默认数据库
	delete(m2, "information_schema")
	delete(m2, "performance_schema")
	delete(m2, "sys")
	delete(m2, "mysql")
	delete(m2, "test")
	//转换成切片
	m3 := mapToSlice(m2)
	return m3
}

//CreateDatabase 创建数据库并创建对应的用户
func CreateDatabase(basepath string, rootPass string, username string, dataName string, dataPwd string) {
	mysqlCommand := ""

	//创建数据库
	mysqlCommand = `create database ` + dataName + ` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci`
	ExecLinuxCommand(`cd ` + basepath + ` && docker-compose exec mysql bash -c "mysql -uroot -p` + rootPass + ` -e '` + mysqlCommand + `'"`)

	//先删用户
	mysqlCommand = `drop user \"` + username + `\"@\"%\"`
	ExecLinuxCommand(`cd ` + basepath + ` && docker-compose exec mysql bash -c "mysql -uroot -p` + rootPass + ` -e '` + mysqlCommand + `'"`)

	//创建用户
	mysqlCommand = `grant all privileges on ` + dataName + `.* to \"` + username + `\"@\"%\" identified by \"` + dataPwd + `\"`
	ExecLinuxCommand(`cd ` + basepath + ` && docker-compose exec mysql bash -c "mysql -uroot -p` + rootPass + ` -e '` + mysqlCommand + `'"`)

	ExecLinuxCommand(`cd ` + basepath + ` && docker-compose exec mysql bash -c "mysql -uroot -p` + rootPass + ` -e 'flush privileges'"`)
}
