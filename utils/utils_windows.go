package utils

import (
	"fmt"
	"os"
	"strings"
)

// SetGoannieEnv 设置windows环境变量
func SetGoannieEnv(path string) error {
	if strings.Index(os.Getenv("PATH"), path) == -1 {
		err := os.Setenv("PATH", fmt.Sprintf("%s;%s", os.Getenv("PATH"), path))
		if err != nil {
			return err
		}
	}
	return nil
}
