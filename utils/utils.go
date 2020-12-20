package utils

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// GetStrInput 获取字符串控制台输入
func GetStrInput(info string, outStr *string) error {
	color.New(color.FgGreen).Printf("$ %s：", info)
	reader := bufio.NewReader(os.Stdin)
	data, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	*outStr = strings.Replace(string(data), "\n", "", -1)
	return nil
}

// GetIntInput 获取int控制台输入
func GetIntInput(info string, outInt *int) error {
	var (
		str string
		err error
	)
	if err = GetStrInput(info, &str); err != nil {
		return err
	}
	if *outInt, err = strconv.Atoi(str); err != nil {
		return err
	}
	return nil
}

// GetUint64Input 获取Uint64控制台输入
func GetUint64Input(info string, outUint64 *uint64) error {
	var (
		str string
		err error
	)
	if err = GetStrInput(info, &str); err != nil {
		return err
	}
	if *outUint64, err = strconv.ParseUint(str, 10, 0); err != nil {
		return err
	}
	return nil
}

// GetDirInput 获取目录控制台输入
func GetDirInput(info string, outDir *string) error {
	var str string
	if err := GetStrInput(info, &str); err != nil {
		return err
	}
	_, err := os.Stat(str)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return err
	}
	*outDir = str
	return nil
}

// ExitInfo 结束程序
func ExitInfo() {
	var exit string
	GetStrInput("回车退出", &exit)
	os.Exit(1)
}

// ErrInfo 打印错误信息
func ErrInfo(errInfo string) {
	color.New(color.FgRed).Printf("错误信息：%s\n", errInfo)
}

// Infoln 打印信息
func Infoln(a ...interface{}) {
	color.New(color.FgBlue).Println(a...)
}

// Infof 打印信息
func Infof(format string, a ...interface{}) {
	color.New(color.FgBlue).Printf(format, a...)
}

// Info 打印信息
func Info(a ...interface{}) {
	color.New(color.FgBlue).Print(a...)
}

// GetCurrentPath 获取程序所在目录
func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\". `)
	}
	return string(path[0 : i+1]), nil
}

// IsExist 文件夹或文件是否存在
func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Domain get the domain of given URL
func Domain(url string) string {
	domainPattern := `([a-z0-9][-a-z0-9]{0,62})\.` +
		`(com\.cn|com\.hk|` +
		`cn|com|net|edu|gov|biz|org|info|pro|name|xxx|xyz|be|` +
		`me|top|cc|tv|tt)`
	domain := MatchOneOf(url, domainPattern)
	if domain != nil {
		return domain[1]
	}
	return ""
}

// MatchOneOf match one of the patterns
func MatchOneOf(text string, patterns ...string) []string {
	var (
		re    *regexp.Regexp
		value []string
	)
	for _, pattern := range patterns {
		// (?flags): set flags within current group; non-capturing
		// s: let . match \n (default false)
		// https://github.com/google/re2/wiki/Syntax
		re = regexp.MustCompile(pattern)
		value = re.FindStringSubmatch(text)
		if len(value) > 0 {
			return value
		}
	}
	return nil
}
