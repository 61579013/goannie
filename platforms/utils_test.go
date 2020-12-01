package platforms

import (
	"testing"
)

func TestT(t *testing.T) {
	err := RunTxDetail(RunType{
		URL:           "https://v.qq.com/detail/m/mzc002001kt6n30.html",
		SavePath:      "C:/Users/Administrator/Desktop/新建文件夹 (5)",
		CookieFile:    "",
		DefaultCookie: "",
		IsDeWeight:    false,
		RedisConn:     nil,
	}, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
}
