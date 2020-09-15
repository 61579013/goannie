package binary_test

import (
	"fmt"
	"testing"

	"gitee.com/rock_rabbit/goannie/binary"
)

func TestGetVersionJSON(t *testing.T) {
	resData, err := binary.GetVersionJSON()
	if err != nil {
		t.Log(err)
		return
	}
	fmt.Println(resData)
}

func TestGetNetworkJSONFile(t *testing.T) {
	resData, err := binary.GetNetworkJSONFile()
	if err != nil {
		t.Log(err)
		return
	}
	fmt.Println(resData)
}

func TestIsUpdate(t *testing.T) {
	resData, err := binary.IsUpdate()
	if err != nil {
		t.Log(err)
		return
	}
	fmt.Println(resData)
}

func TestExecUpdate(t *testing.T) {
	updateJSON, err := binary.IsUpdate()
	if err != nil {
		t.Log(err)
		return
	}
	if err = binary.ExecUpdate(updateJSON); err != nil {
		t.Log(err)
		return
	}

}
