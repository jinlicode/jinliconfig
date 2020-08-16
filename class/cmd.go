package class

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func ExecNginxSslCreate() {

}

// ExecLinuxCommand 执行liunx命令 且回显
func ExecLinuxCommand(CommandString string) {

	cmd := exec.Command("bash", "-c", CommandString)

	//显示运行的命令
	stdout, err := cmd.StdoutPipe()

	//直接错误 直接断下
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	cmd.Start()
	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		//循环打印每行代码
		fmt.Printf(line)
	}

	cmd.Wait()
}

//ExecDockerInstall 执行安装docker操作
func ExecDockerInstall() {
	//关闭防火墙
	ExecLinuxCommand("systemctl stop firewalld.service && systemctl disable firewalld.service && setenforce 0 && sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config")
	// step 1: 安装必要的一些系统工具
	ExecLinuxCommand("sudo yum install -y yum-utils device-mapper-persistent-data lvm2 git epel-*")
	// Step 2: 添加软件源信息
	ExecLinuxCommand("sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo")
	// Step 3: 更新并安装 Docker-CE
	ExecLinuxCommand("sudo yum makecache fast && sudo yum -y install docker-ce docker-compose")
	// Step 4: 开启Docker服务
	ExecLinuxCommand("sudo systemctl start docker")
	// step 5: 设置开机启动
	ExecLinuxCommand("sudo systemctl enable docker")

	//设置docker源
	ExecLinuxCommand("mkdir -p /etc/docker")
	WriteFile("/etc/docker/daemon.json", `{"registry-mirrors":["https://docker.mirrors.ustc.edu.cn"],"log-driver":"json-file","log-opts":{"max-size":"1m","max-file":"1"}}`)

	//重载docker
	ExecLinuxCommand("sudo systemctl daemon-reload && sudo systemctl restart docker")
}
