package downloader

import (
	"fmt"
	"mime"
	"regexp"
	"strings"
	"time"
)

// GetFilename 生成文件名称
func GetFilename(url, contentDisposition, contentType string) string {
	_, params, _ := mime.ParseMediaType(contentDisposition)
	if filename, ok := params["filename"]; ok && filename != "" {
		return filename
	}
	var name string
	namelist := strings.Split(url, "/")
	if len(namelist) != 0 {
		name = namelist[len(namelist)-1]
		splistname := strings.Split(name, "?")
		if len(splistname) != 0 {
			name = splistname[0]
		}
	}
	name = GetFiltrationFilename(name)
	if name == "" {
		var ext string
		extlist, _ := mime.ExtensionsByType(contentType)
		if len(extlist) != 0 {
			ext = extlist[0]
		}
		name = fmt.Sprintf("file_%d%s", time.Now().UnixNano(), ext)
	}
	return name
}

// GetFiltrationFilename 返回过滤后的文件名
func GetFiltrationFilename(name string) string {
	if name == "" {
		return ""
	}
	return strings.Join(regexp.MustCompile(`[^?\\/:*<>|]+`).FindAllString(name, -1), "_")
}
