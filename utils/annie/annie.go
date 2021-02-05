package annie

import (
	"fmt"
	"os"
	"os/exec"

	"gitee.com/rock_rabbit/goannie/config"
)

// Download 使用annie下载
func Download(u, savePath, cookie string) error {
	stream := config.GetString("app.stream")
	retryTimes := fmt.Sprint(config.GetInt("app.retryTimes"))
	arg := []string{}
	// cookie 设置
	arg = append(arg, "-c", cookie)
	// 下载项设置
	if stream != "default" && stream != "" {
		arg = append(arg, "-f", stream)
	}
	// 超时次数、保存地址、下载url
	arg = append(arg, []string{"-retry", retryTimes, "-o", savePath, u}...)
	cmd := exec.Command("annie", arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}
