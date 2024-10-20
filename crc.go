package main

import (
	"fmt"
	"strconv"
)

// 生成 CRC32 表
var crctable [256]uint32

func crcinit() {
	for i := 0; i < 256; i++ {
		crcreg := uint32(i)
		for j := 0; j < 8; j++ {
			if crcreg&1 != 0 {
				crcreg = 0xEDB88320 ^ (crcreg >> 1)
			} else {
				crcreg >>= 1
			}
		}
		crctable[i] = crcreg
	}
}

func crack(input string) string {
	// CRC32 计算函数
	crc32 := func(s string) uint32 {
		crcstart := uint32(0xFFFFFFFF)
		for i := 0; i < len(s); i++ {
			index := (crcstart ^ uint32(s[i])) & 255
			crcstart = (crcstart >> 8) ^ crctable[index]
		}
		return crcstart
	}

	// 获取最后使用的索引
	crc32LastIndex := func(s string) int {
		var index uint32
		crcstart := uint32(0xFFFFFFFF)
		for i := 0; i < len(s); i++ {
			index = (crcstart ^ uint32(s[i])) & 255
			crcstart = (crcstart >> 8) ^ crctable[index]
		}
		return int(index)
	}

	// 获取 CRC 索引
	getCRCIndex := func(t int) int {
		for i := 0; i < 256; i++ {
			if crctable[i]>>24 == uint32(t) {
				return i
			}
		}
		return -1
	}

	// 深度检查
	deepCheck := func(i int, index []int) (bool, string) {
		hashcode := crc32(strconv.Itoa(i))
		tc := int(hashcode&0xff) ^ index[2]
		if tc < 48 || tc > 57 {
			return false, ""
		}

		result := string(rune(tc - 48 + '0'))
		hashcode = crctable[index[2]] ^ (hashcode >> 8)

		tc = int(hashcode&0xff) ^ index[1]
		if tc < 48 || tc > 57 {
			return false, ""
		}
		result += string(rune(tc - 48 + '0'))
		hashcode = crctable[index[1]] ^ (hashcode >> 8)

		tc = int(hashcode&0xff) ^ index[0]
		if tc < 48 || tc > 57 {
			return false, ""
		}
		result += string(rune(tc - 48 + '0'))
		return true, result
	}

	var index = make([]int, 4)
	ht, _ := strconv.ParseInt(input, 16, 64)
	ht ^= 0xFFFFFFFF
	for i := 3; i >= 0; i-- {
		index[3-i] = getCRCIndex(int(ht >> (i * 8)))
		snum := crctable[index[3-i]]
		ht ^= int64(snum >> ((3 - i) * 8))
	}

	for i := 0; i < 100000000; i++ {
		if crc32LastIndex(strconv.Itoa(i)) == index[3] {
			if valid, result := deepCheck(i, index); valid {
				return fmt.Sprintf("%d%s", i, result)
			}
		}
	}

	return "-1"
}
