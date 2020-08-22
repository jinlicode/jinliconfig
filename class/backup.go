package class

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

//FlagBackup 命令行整体备份
func FlagBackupExec(basepath string) {

	backup := flag.String("backup", "", "--backup=db 备份数据库 --backup=site 备份网站 --backup=all 备份全部")
	flag.Parse()
	if CheckFileExist(basepath + "docker-compose.yaml") {
		if *backup == "db" {

			//删除7天谴的db目录
			ExecLinuxCommand(`find ` + basepath + `autobackup/database/ -type f -mtime +7 -exec rm -rf {} \;`)

			BackupMysqlPassword := ReadMysqlRootPassword(basepath)
			fileName := time.Now().Format("20060102150405") + ".sql.gz"
			backupName := basepath + "autobackup/database/" + fileName
			ExecLinuxCommand(`cd ` + basepath + ` && docker-compose exec -T mysql mysqldump --all-databases -uroot -p` + BackupMysqlPassword + ` |gzip >` + backupName)
			os.Exit(1)
		} else if *backup == "site" {

			//删除7天谴的site目录
			ExecLinuxCommand(`find ` + basepath + `autobackup/site/ -type f -mtime +7 -exec rm -rf {} \;`)

			fileName := time.Now().Format("20060102150405") + ".tar.gz"
			backupName := basepath + "autobackup/site/" + fileName
			ExecLinuxCommand(`tar -zcf ` + backupName + ` ` + basepath + `code/`)
			os.Exit(2)

		} else if *backup == "all" {

			//删除7天谴的db目录
			ExecLinuxCommand(`find ` + basepath + `autobackup/database/ -type f -mtime +7 -exec rm -rf {} \;`)
			//删除7天谴的site目录
			ExecLinuxCommand(`find ` + basepath + `autobackup/site/ -type f -mtime +7 -exec rm -rf {} \;`)

			fileName := ""
			backupName := ""
			BackupMysqlPassword := ReadMysqlRootPassword(basepath)
			fileName = time.Now().Format("20060102150405") + ".sql.gz"
			backupName = basepath + "autobackup/database/" + fileName
			ExecLinuxCommand(`cd ` + basepath + ` && docker-compose exec -T mysql mysqldump --all-databases -uroot -p` + BackupMysqlPassword + ` |gzip >` + backupName)

			fileName = time.Now().Format("20060102150405") + ".tar.gz"
			backupName = basepath + "autobackup/site/" + fileName
			ExecLinuxCommand(`tar -zcf ` + backupName + ` ` + basepath + `code/`)
			os.Exit(3)
		}
	}
}

// BackupSiteManage 网站管理
func BackupSiteManage(basepath string, ExistSiteSlice []string) bool {

	WebBuckupSelectOption := []string{}
	WebBuckupSelectOption = append(ExistSiteSlice, "返回上层")
	WebBuckupSelect := ConsoleOptionsSelect("请选择您需要备份的网站", WebBuckupSelectOption, "请输入选项")

	switch WebBuckupSelect {
	case "返回上层":
		fmt.Println("返回上层")
		return false
	case WebBuckupSelect:

		WebSiteBuckupSelect := ConsoleOptionsSelect("请选择您需要备份选项", []string{WebBuckupSelect + "的数据库备份", WebBuckupSelect + "的网站备份", WebBuckupSelect + "的数据库+网站备份", "返回上层"}, "请输入选项")
		WebSiteBuckupSelectString := strings.Replace(WebBuckupSelect, ".", "_", -1)

		switch WebSiteBuckupSelect {
		case "返回上层":
			fmt.Println("返回上层")
			return false
		case WebBuckupSelect + "的数据库备份":

			MysqlPassword := ReadMysqlRootPassword(basepath)
			fileName := WebSiteBuckupSelectString + "_" + time.Now().Format("20060102150405") + ".sql.gz"
			backupName := basepath + "backup/database/" + fileName
			ExecLinuxCommand(`cd ` + basepath + ` && docker-compose exec -T mysql mysqldump -uroot -p` + MysqlPassword + ` ` + WebSiteBuckupSelectString + ` |gzip >` + backupName)
			fmt.Println("数据库备份成功，备份在" + backupName)
			return false
		case WebBuckupSelect + "的网站备份":

			fileName := WebSiteBuckupSelectString + "_" + time.Now().Format("20060102150405") + ".tar.gz"
			backupName := basepath + "backup/site/" + fileName
			ExecLinuxCommand(`tar -zcf ` + backupName + ` ` + basepath + `code/` + WebSiteBuckupSelectString)
			fmt.Println("网站备份成功，备份在" + backupName)
			return false

		case WebBuckupSelect + "的数据库+网站备份":
			fileName := ""
			backupName := ""
			MysqlPassword := ReadMysqlRootPassword(basepath)
			fileName = WebSiteBuckupSelectString + "_" + time.Now().Format("20060102150405") + ".sql.gz"
			backupName = basepath + "backup/database/" + fileName
			ExecLinuxCommand(`cd ` + basepath + ` && docker-compose exec -T mysql mysqldump -uroot -p` + MysqlPassword + ` ` + WebSiteBuckupSelectString + ` |gzip >` + backupName)
			fmt.Println("数据库备份成功，备份在" + backupName)

			fileName = WebSiteBuckupSelectString + "_" + time.Now().Format("20060102150405") + ".tar.gz"
			backupName = basepath + "backup/site/" + fileName
			ExecLinuxCommand(`tar -zcf ` + backupName + ` ` + basepath + `code/` + WebSiteBuckupSelectString)
			fmt.Println("网站备份成功，备份在" + backupName)
			return false
		}
	}
	return true
}
