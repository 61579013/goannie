package platforms_test

import (
	"fmt"
	"testing"

	"gitee.com/rock_rabbit/goannie/platforms"
)

func TestT(t *testing.T) {
	fmt.Println(platforms.YkGetUID("http://i.youku.com/u/UNzA0Nzk1Njg3Mg==?spm=a2h0c.8166622.PhoneSokuUgc_13.3"))
	fmt.Println(platforms.YkGetUID("http://i.youku.com/u/UNzA0Nzk1Njg3Mg==/?"))
	fmt.Println(platforms.YkGetUID("http://i.youku.com/u/UNzA0Nzk1Njg3Mg=="))
	fmt.Println(platforms.YkGetUID("http://i.youku.com/i/UNzA0Nzk1Njg3Mg==?spm=a2h0c.8166622.PhoneSokuUgc_13.3"))
	fmt.Println(platforms.YkGetUID("http://i.youku.com/i/UNzA0Nzk1Njg3Mg==/?"))
	fmt.Println(platforms.YkGetUID("http://i.youku.com/i/UNzA0Nzk1Njg3Mg=="))
	fmt.Println(platforms.YkGetUID("https://i.youku.com/yijialianmeng?spm=a2hzp.8244740.0.0"))
}

// func TestDouyin(t *testing.T) {
// 	// page, count, err := platforms.DyGetUserListPage("MS4wLjABAAAArrTGJuyQzorGDmgxrRBGNme7R87ahhFijRqxR_6Ubf29Gj0j34n1YuS6DZXbYGfa")
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	// fmt.Println(page, count)

// 	maxCursor := 0
// 	videoCount := 0
// 	for {
// 		res, err := platforms.DyGetUserlistAPI("MS4wLjABAAAArrTGJuyQzorGDmgxrRBGNme7R87ahhFijRqxR_6Ubf29Gj0j34n1YuS6DZXbYGfa", maxCursor, 0)
// 		if err != nil {
// 			fmt.Println(err)
// 			break
// 		}
// 		if res.HasMore == 0 {
// 			break
// 		}
// 		maxCursor = res.MaxCursor
// 		p := 0
// 		for _, _ = range res.AwemeList {
// 			p++
// 			videoCount++
// 		}
// 		fmt.Printf("\r数量：%d 频率：%d", videoCount, p)
// 	}
// 	fmt.Println("\n数量：", videoCount)
// }

// func TestDouyinSignature(t *testing.T) {
// 	fmt.Println(platforms.DyGetSignature("123456"))
// }
