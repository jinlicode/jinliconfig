package class

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func ExecNginxSslCreate() {

}

// ExecLinuxCommand 执行liunx命令 且回显,不支持取得返回值
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

//执行命令，并且返回字符串回显
func ExecLinuxCommandReturn(BashCommand string) string {
	cmd := exec.Command("bash", "-c", BashCommand)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(out))
	}
	return string(out)
}

//版本号对比version1系统版本，version2为最低版本 返回值：1为正常，-1为版本过低， 0为版本匹配
func compareVersion(version1 string, version2 string) int {
	versionA := strings.Split(version1, ".")
	versionB := strings.Split(version2, ".")

	for i := len(versionA); i < 4; i++ {
		versionA = append(versionA, "0")
	}
	for i := len(versionB); i < 4; i++ {
		versionB = append(versionB, "0")
	}
	for i := 0; i < 4; i++ {
		version1, _ := strconv.Atoi(versionA[i])
		version2, _ := strconv.Atoi(versionB[i])
		if version1 == version2 {
			continue
		} else if version1 > version2 {
			return 1
		} else {
			return -1
		}
	}
	return 0
}

//检查docker是否符合安装环境
func ChkDokcerInstall() bool {
	DockerReturen := ExecLinuxCommandReturn("docker -v")
	DockerReturenRe, _ := regexp.Compile(`\d+(?:\.\d+)+`)
	DockerVersion := DockerReturenRe.FindString(DockerReturen)
	DockerInstallResult := compareVersion(DockerVersion, "19.03.12")
	if DockerInstallResult != -1 {
		DockerComposeReturen := ExecLinuxCommandReturn("docker-compose -v")
		DockerComposeReturenRe, _ := regexp.Compile(`\d+(?:\.\d+)+`)
		DockerComposeVersion := DockerComposeReturenRe.FindString(DockerComposeReturen)
		DockerComposeInstallResult := compareVersion(DockerComposeVersion, "1.18.0")
		if DockerComposeInstallResult == -1 {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}

//卸载docker方法，返回命令提示
func ChkDokcerRemove() {
	ExecLinuxCommand("systemctl disable docker.service")
	ExecLinuxCommand("systemctl stop docker.service")
	ExecLinuxCommand("rm -rf /var/lib/docker/")
	ExecLinuxCommand("rm -rf /etc/docker")
	ExecLinuxCommand("yum remove -y docker*")
}

//ExecDockerInstall 执行安装docker操作
func ExecDockerInstall() {
	//安装源
	ExecLinuxCommand("yum install epel-release -y")
	//关闭防火墙
	ExecLinuxCommand("systemctl stop firewalld.service && systemctl disable firewalld.service && setenforce 0 && sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config")
	// step 1: 安装必要的一些系统工具
	ExecLinuxCommand("sudo yum install -y yum-utils device-mapper-persistent-data lvm2 git epel-*")
	// Step 2: 添加软件源信息
	ExecLinuxCommand("sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo")
	// Step 3: 更新并安装 Docker-CE
	ExecLinuxCommand("sudo yum makecache fast && sudo yum -y install docker-ce-19.03.12 docker-compose-1.18.0")
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
