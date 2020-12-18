package platforms

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestT(t *testing.T) {
	data, err := xgGetSsrHydratedData("https://www.ixigua.com/6885238748574351880?logTag=1Tstxc-bne_8ehlxkGXjr")
	if err != nil {
		t.Log(err)
		return
	}

	fmt.Println(data)
	decoded, _ := base64.StdEncoding.DecodeString(xgGetNewDownloadUrl(data))
	fmt.Println(string(decoded))

}
