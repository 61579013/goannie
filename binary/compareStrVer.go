package binary

import (
	"strings"
)

/**
* 对比字符串版本号
* Created by 908204694@qq.com
 */

// VersionBig 版本号大
const VersionBig = 1

// VersionSmall 版本号小
const VersionSmall = 2

// VersionEqual 版本号相等
const VersionEqual = 0

func compareStrVer(verA, verB string) int {

	verStrArrA := spliteStrByNet(verA)
	verStrArrB := spliteStrByNet(verB)

	lenStrA := len(verStrArrA)
	lenStrB := len(verStrArrB)

	if lenStrA != lenStrB {
		panic("版本号格式不一致")
	}
	return compareArrStrVers(verStrArrA, verStrArrB)
}

// 比较版本号字符串数组
func compareArrStrVers(verA, verB []string) int {

	for index, _ := range verA {

		littleResult := compareLittleVer(verA[index], verB[index])

		if littleResult != VersionEqual {
			return littleResult
		}
	}

	return VersionEqual
}

//
// 比较小版本号字符串
//
func compareLittleVer(verA, verB string) int {

	bytesA := []byte(verA)
	bytesB := []byte(verB)

	lenA := len(bytesA)
	lenB := len(bytesB)
	if lenA > lenB {
		return VersionBig
	}

	if lenA < lenB {
		return VersionSmall
	}

	//如果长度相等则按byte位进行比较

	return compareByBytes(bytesA, bytesB)
}

// 按byte位进行比较小版本号
func compareByBytes(verA, verB []byte) int {

	for index, _ := range verA {
		if verA[index] > verB[index] {
			return VersionBig
		}
		if verA[index] < verB[index] {
			return VersionSmall
		}

	}

	return VersionEqual
}

// 按“.”分割版本号为小版本号的字符串数组
func spliteStrByNet(strV string) []string {

	return strings.Split(strV, ".")
}
