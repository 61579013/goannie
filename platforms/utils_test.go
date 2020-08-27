package platforms_test

import (
	"fmt"
	"gitee.com/rock_rabbit/goannie/platforms"
	"testing"
)

func TestGetVideoID(t *testing.T) {
	isVID,_ := platforms.AddVideoID("bilibili","12345611")
	fmt.Println(isVID)
}
