package manage

import (
	"flag"
	"fmt"
	"jinliconfig/class"
	"os"
	"strings"
	"time"
)

//FlagBackupExec 命令行整体备份
func FlagBackupExec(basepath string) {

	backup := flag.String("backup", "", "--backup=db 备份数据库 --backup=site 备份网站 --backup=all 备份全部")
	flag.Parse()
	if class.CheckFileExist(basepath + "docker-compose.yaml") {
		if *backup == "db" {

			//删除7天谴的db目录
			class.ExecLinuxCommand(`find ` + basepath + `autobackup/database/ -type f -mtime +7 -exec rm -rf {} \;`)

			BackupMysqlPassword := class.ReadMysqlRootPassword(basepath)
			BackupReadMysqlHost := class.ReadMysqlHost(basepath)
			fileName := time.Now().Format("20060102150405") + ".sql.gz"
			backupName := basepath + "autobackup/database/" + fileName
			class.ExecLinuxCommand(`cd ` + basepath + ` && docker-compose exec -T mysql mysqldump -h` + BackupReadMysqlHost + ` --all-databases -uroot -p` + BackupMysqlPassword + ` |gzip >` + backupName)
			os.Exit(1)
		} else if *backup == "site" {

			//删除7天谴的site目录
			class.ExecLinuxCommand(`find ` + basepath + `autobackup/site/ -type f -mtime +7 -exec rm -rf {} \;`)

			fileName := time.Now().Format("20060102150405") + ".tar.gz"
			backupName := basepath + "autobackup/site/" + fileName
			class.ExecLinuxCommand(`tar -zcf ` + backupName + ` ` + basepath + `code/`)
			os.Exit(2)

		} else if *backup == "all" {

			//删除7天谴的db目录
			class.ExecLinuxCommand(`find ` + basepath + `autobackup/database/ -type f -mtime +7 -exec rm -rf {} \;`)
			//删除7天谴的site目录
			class.ExecLinuxCommand(`find ` + basepath + `autobackup/site/ -type f -mtime +7 -exec rm -rf {} \;`)

			fileName := ""
			backupName := ""
			BackupMysqlPassword := class.ReadMysqlRootPassword(basepath)
			BackupReadMysqlHost := class.ReadMysqlHost(basepath)
			fileName = time.Now().Format("20060102150405") + ".sql.gz"
			backupName = basepath + "autobackup/database/" + fileName
			class.ExecLinuxCommand(`cd ` + basepath + ` && docker-compose exec -T mysql mysqldump -h` + BackupReadMysqlHost + ` --all-databases -uroot -p` + BackupMysqlPassword + ` |gzip >` + backupName)

			fileName = time.Now().Format("20060102150405") + ".tar.gz"
			backupName = basepath + "autobackup/site/" + fileName
			class.ExecLinuxCommand(`tar -zcf ` + backupName + ` ` + basepath + `code/`)
			os.Exit(3)
		}
	}
}

// BackupSiteManage 网站管理
func BackupSiteManage(basepath string, ExistSiteSlice []string) bool {

	WebBuckupSelectOption := []string{}
	WebBuckupSelectOption = append(ExistSiteSlice, "返回上层")
	WebBuckupSelect := class.ConsoleOptionsSelect("请选择您需要备份的网站", WebBuckupSelectOption, "请输入选项")
	if WebBuckupSelect == "interrupt" {
		fmt.Println("您已强制退出")
		os.Exit(1)
	}

	switch WebBuckupSelect {
	case "返回上层":
		fmt.Println("返回上层")
		return false
	case WebBuckupSelect:

		WebSiteBuckupSelect := class.ConsoleOptionsSelect("请选择您需要备份选项", []string{WebBuckupSelect + "的数据库备份", WebBuckupSelect + "的网站备份", WebBuckupSelect + "的数据库+网站备份", "返回上层"}, "请输入选项")
		WebSiteBuckupSelectString := strings.Replace(WebBuckupSelect, ".", "_", -1)

		switch WebSiteBuckupSelect {
		case "返回上层":
			fmt.Println("返回上层")
			return false
		case WebBuckupSelect + "的数据库备份":

			MysqlPassword := class.ReadMysqlRootPassword(basepath)
			MysqlHost := class.ReadMysqlHost(basepath)

			fileName := WebSiteBuckupSelectString + "_" + time.Now().Format("20060102150405") + ".sql.gz"
			backupName := basepath + "backup/database/" + fileName
			class.ExecLinuxCommand(`cd ` + basepath + ` && docker-compose exec -T mysql mysqldump -h` + MysqlHost + ` -uroot -p` + MysqlPassword + ` ` + WebSiteBuckupSelectString + ` |gzip >` + backupName)
			class.PrintHr()
			fmt.Println("数据库备份成功，备份在" + backupName)
			class.PrintHr()
			return false
		case WebBuckupSelect + "的网站备份":

			fileName := WebSiteBuckupSelectString + "_" + time.Now().Format("20060102150405") + ".tar.gz"
			backupName := basepath + "backup/site/" + fileName
			class.ExecLinuxCommand(`tar -zcf ` + backupName + ` ` + basepath + `code/` + WebSiteBuckupSelectString)
			class.PrintHr()
			fmt.Println("网站备份成功，备份在" + backupName)
			class.PrintHr()
			return false

		case WebBuckupSelect + "的数据库+网站备份":
			fileName := ""
			backupName := ""
			MysqlPassword := class.ReadMysqlRootPassword(basepath)
			MysqlHost := class.ReadMysqlHost(basepath)

			fileName = WebSiteBuckupSelectString + "_" + time.Now().Format("20060102150405") + ".sql.gz"
			backupName = basepath + "backup/database/" + fileName
			class.ExecLinuxCommand(`cd ` + basepath + ` && docker-compose exec -T mysql mysqldump -h` + MysqlHost + ` -uroot -p` + MysqlPassword + ` ` + WebSiteBuckupSelectString + ` |gzip >` + backupName)
			class.PrintHr()
			fmt.Println("数据库备份成功，备份在" + backupName)
			class.PrintHr()

			fileName = WebSiteBuckupSelectString + "_" + time.Now().Format("20060102150405") + ".tar.gz"
			backupName = basepath + "backup/site/" + fileName
			class.ExecLinuxCommand(`tar -zcf ` + backupName + ` ` + basepath + `code/` + WebSiteBuckupSelectString)
			fmt.Println("网站备份成功，备份在" + backupName)
			class.PrintHr()
			return false
		}
	}
	return true
}
