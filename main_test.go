package main_test

import (
	pf "gitee.com/rock_rabbit/goannie/platforms"
	"testing"
)

func TestXgUserList(t *testing.T) {
	runType := pf.RunType{
		"https://www.ixigua.com/home/50680494408/?video_card_type=shortvideo",
		"G:\\新建文件夹",
		"",
	}
	err := pf.RunXgUserList(runType)
	if err != nil{
		t.Log(err)
	}
}

// 测试单西瓜下载
/*func TestXgOne(t *testing.T) {
	runType := pf.RunType{
		"https://www.ixigua.com/6832194590221533707",
		"G:\\新建文件夹",
		"",
	}
	err := pf.RunXgOne(runType)
	if err != nil{
		t.Log(err)
	}
	runType = pf.RunType{
		"https://m.ixigua.com/video/6729350683222344199",
		"G:\\新建文件夹",
		"",
	}
	err = pf.RunXgOne(runType)
	if err != nil{
		t.Log(err)
	}
}*/
