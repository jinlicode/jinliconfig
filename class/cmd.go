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
