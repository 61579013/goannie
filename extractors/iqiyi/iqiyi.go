package iqiyi

import (
	"errors"
	"fmt"

	"gitee.com/rock_rabbit/goannie/config"
	"gitee.com/rock_rabbit/goannie/extractors/types"
	"gitee.com/rock_rabbit/goannie/request"
	"gitee.com/rock_rabbit/goannie/utils"
	"gitee.com/rock_rabbit/goannie/utils/annie"
)

const (
	referer       = "https://www.iqiyi.com"
	defaultCookie = "QC005=a7cebba4b844a08d30045630804e00a3; QC006=fc41aa8d932d6fb52744202b52902a3a; QC173=0; __uuid=096abec1-d1c3-5071-9614-0d315b15f650; QP0030=1; T00404=ca8ed5c103e1547fea010cf2a30fb9b0; QP0013=; QC124=1%7C0; idx=%7B%22ip%22%3A%22101X21X43X90%22%2C%22geocode%22%3A%221156130500%22%2C%22exptime%22%3A1612967268%7D; QC021=%5B%7B%22key%22%3A%22%E6%B3%BD%E5%A1%94%E5%A5%A5%E7%89%B9%E6%9B%BC%EF%BC%9A%E6%B4%8B%E5%AD%90%E5%BD%BB%E5%BA%95%E9%BB%91%E5%8C%96%EF%BC%8C%E7%AB%9F%E8%BF%9E%E6%B3%BD%E5%A1%94%E9%83%BD%E4%B8%8B%E7%8B%A0%E6%89%8B%22%7D%2C%7B%22key%22%3A%2202%3A50%20%E6%B3%BD%E5%A1%94%E5%A5%A5%E7%89%B9%E6%9B%BC%EF%BC%9A%E6%B4%8B%E5%AD%90%E5%BD%BB%E5%BA%95%E9%BB%91%E5%8C%96%EF%BC%8C%E7%AB%9F%E8%BF%9E%E6%B3%BD%E5%A1%94%E9%83%BD%E4%B8%8B%E7%8B%A0%E6%89%8B%22%7D%2C%7B%22key%22%3A%22%E6%B3%BD%E5%A1%94%E5%A5%A5%E7%89%B9%E6%9B%BC%EF%BC%9A%E6%89%80%E6%9C%89%E4%BA%BA%E9%83%BD%E6%9D%A5%E6%8A%A2%E7%9A%84%E6%81%B6%E9%AD%94%E7%A2%8E%E7%89%87%EF%BC%8C%E8%B4%9D%E5%88%A9%E4%BA%9A%E7%BB%86%E8%83%9E%22%7D%2C%7B%22key%22%3A%22%E7%8E%8B%E7%89%8C%E6%8E%A8%E7%90%86%E9%A6%86%22%7D%5D; QC175=%7B%22upd%22%3Atrue%2C%22ct%22%3A%22%22%7D; PCAU=0; QC178=true; QC159=%7B%22color%22%3A%22FFFFFF%22%2C%22channelConfig%22%3A0%2C%22isOpen%22%3A1%2C%22speed%22%3A10%2C%22density%22%3A30%2C%22opacity%22%3A86%2C%22isFilterColorFont%22%3A1%2C%22proofShield%22%3A0%2C%22forcedFontSize%22%3A24%2C%22isFilterImage%22%3A1%2C%22hadTip%22%3A1%2C%22isFilterHongbao%22%3A0%2C%22hideRoleTip%22%3A1%7D; Hm_lvt_292c77bd6e6064e926d1d58f63241745=1612103265,1612534525,1612569616; Hm_lpvt_292c77bd6e6064e926d1d58f63241745=1612569616; QC176=%7B%22state%22%3A0%2C%22ct%22%3A1612569616783%7D; Hm_lvt_53b7374a63c37483e5dd97d78d9bb36e=1612103967,1612103977,1612536050,1612569621; QC007=DIRECT; QC008=1609156502.1612536049.1612569620.4; nu=0; CA0001=%7B%22code%22%3A%22b20f64bd520f631bA00000%22%2C%22type%22%3A132%2C%22detail%22%3A%7B%22fv%22%3A%22%22%2C%22text1%22%3A%22%E9%A6%96%E6%9C%88%E7%89%B9%E6%83%A0%22%2C%22autoRenew%22%3A%22true%22%2C%22type1%22%3A%7B%22vipType%22%3A%229c4e4c1c18827a41%22%2C%22type%22%3A%225%22%2C%22fc%22%3A%2290647efdbbc99688%22%7D%2C%22packageAmount%22%3A%221%22%7D%2C%22passportId%22%3A%22%22%2C%22locale%22%3A%22cn%22%7D; Hm_lpvt_53b7374a63c37483e5dd97d78d9bb36e=1612569904; T00700=EgcI9L-tIRABEgcIz7-tIRABEgcI67-tIRACEgcIkMDtIRABEgcIg8DtIRABEgcI0b-tIRABEgcI4b-tIRAB; QP0027=22; TQC002=type%3Djspfmc140109%26pla%3D11%26uid%3Da7cebba4b844a08d30045630804e00a3%26ppuid%3D%26brs%3DCHROME%26pgtype%3Dplay%26purl%3Dhttps%3A%252F%252Fwww.iqiyi.com%252Fv_b8m05imcdo.html%26cid%3D1%26tmplt%3D%26tm1%3D1328%2C0; __dfp=a0198984197931448bb1f6a14364744da0f7e7f26b188bfd0b6bbace10547de0dd@1613399249444@1612103250444; QC010=52942892; IMS=IggQAxj_-fqABiokCiA5MTE0NmQ3ZjliMzZjMDc0ZDA2NGEyMjQzYmVmNGYzYhAAciQKIDkxMTQ2ZDdmOWIzNmMwNzRkMDY0YTIyNDNiZWY0ZjNiEACCAQCKASQKIgogOTExNDZkN2Y5YjM2YzA3NGQwNjRhMjI0M2JlZjRmM2I; QP0022=CNC%7CHeBei-121.27.179.191%7Cvcdn_Rcache_taiyuan9_cnc"
)

type extractor struct{}

// Name 平台名称
func (e *extractor) Name() string {
	return "爱奇艺视频"
}

// Key 存储器使用的键
func (e *extractor) Key() string {
	return "iqiyi"
}

// DefCookie 默认cookie
func (e *extractor) DefCookie() string {
	return defaultCookie
}

// New 创建一个解析器
func New() types.Extractor {
	return &extractor{}
}

// Extract 运行解析器
func (e *extractor) Extract(url string, option types.Options) error {
	html, err := request.GetByte(url, referer, nil)
	if err != nil {
		return err
	}
	regexpVid := utils.MatchOneOfByte(html, `param['vid'] = "(.*?)"`, `"vid":"(.*?)"`)
	if len(regexpVid) < 2 {
		return errors.New("获取vid失败")
	}
	vid := string(regexpVid[1])
	if config.GetBool("app.isFiltrationID") && vid != "" && option.Verify.Check(e.Key(), vid) {
		return fmt.Errorf("%s 在过滤库中，修改配置文件isFiltrationID为false不再过滤。", vid)
	}

	if err := annie.Download(url, option.SavePath, option.Cookie); err != nil {
		return err
	}

	option.Verify.Add(e.Key(), vid)
	return nil
}
