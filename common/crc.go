package common

import (
	"fmt"
	"strconv"
)

var table [256]uint32

func init() {
	for i := 0; i < 256; i++ {
		reg := uint32(i)
		for j := 0; j < 8; j++ {
			if reg&1 != 0 {
				reg = 0xEDB88320 ^ (reg >> 1)
			} else {
				reg >>= 1
			}
		}
		table[i] = reg
	}
}

func Crack(input string) string {
	crc32 := func(s string) uint32 {
		start := uint32(0xFFFFFFFF)
		for i := 0; i < len(s); i++ {
			index := (start ^ uint32(s[i])) & 255
			start = (start >> 8) ^ table[index]
		}
		return start
	}

	crc32LastIndex := func(s string) int {
		var index uint32
		start := uint32(0xFFFFFFFF)
		for i := 0; i < len(s); i++ {
			index = (start ^ uint32(s[i])) & 255
			start = (start >> 8) ^ table[index]
		}
		return int(index)
	}

	getCRCIndex := func(t int) int {
		for i := 0; i < 256; i++ {
			if table[i]>>24 == uint32(t) {
				return i
			}
		}
		return -1
	}

	deepCheck := func(i int, index []int) (bool, string) {
		hashcode := crc32(strconv.Itoa(i))
		tc := int(hashcode&0xff) ^ index[2]
		if tc < 48 || tc > 57 {
			return false, ""
		}

		result := string(rune(tc - 48 + '0'))
		hashcode = table[index[2]] ^ (hashcode >> 8)

		tc = int(hashcode&0xff) ^ index[1]
		if tc < 48 || tc > 57 {
			return false, ""
		}
		result += string(rune(tc - 48 + '0'))
		hashcode = table[index[1]] ^ (hashcode >> 8)

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
		snum := table[index[3-i]]
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
