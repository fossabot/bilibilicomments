package main

import (
	"fmt"
	"strconv"
)

// crack 反推bilibili midhash, 返回-1就是没有找到
func crack(input string) string {
	// Initialize the CRC32 table
	var crctable [256]uint32
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

	// Function to calculate CRC32 hash
	crc32 := func(s string) uint32 {
		crcstart := uint32(0xFFFFFFFF)
		for i := 0; i < len(s); i++ {
			index := (crcstart ^ uint32(s[i])) & 255
			crcstart = (crcstart >> 8) ^ crctable[index]
		}
		return crcstart
	}

	// Function to get the last index of CRC32
	crc32LastIndex := func(s string) int {
		crcstart := uint32(0xFFFFFFFF)
		var index uint32
		for i := 0; i < len(s); i++ {
			index = (crcstart ^ uint32(s[i])) & 255
			crcstart = (crcstart >> 8) ^ crctable[index]
		}
		return int(index)
	}

	// Function to find CRC index
	getCRCIndex := func(t int) int {
		for i := 0; i < 256; i++ {
			if crctable[i]>>24 == uint32(t) {
				return i
			}
		}
		return -1
	}

	// Function to perform deep check
	deepCheck := func(i int, index []int) (bool, string) {
		hashcode := crc32(strconv.Itoa(i))
		var result string
		for j := 2; j >= 0; j-- {
			tc := int(hashcode&0xff) ^ index[j]
			if tc < 48 || tc > 57 {
				return false, ""
			}
			result = string(rune(tc-48+'0')) + result
			hashcode = crctable[index[j]] ^ (hashcode >> 8)
		}
		return true, result
	}

	var index [4]int
	ht, _ := strconv.ParseInt(input, 16, 64)
	ht ^= 0xFFFFFFFF

	for i := 3; i >= 0; i-- {
		index[3-i] = getCRCIndex(int(ht >> (i * 8)))
		snum := crctable[index[3-i]]
		ht ^= int64(snum >> ((3 - i) * 8))
	}

	for i := 0; i < 100000000; i++ {
		if crc32LastIndex(strconv.Itoa(i)) == index[3] {
			if valid, result := deepCheck(i, index[:]); valid {
				return fmt.Sprintf("%d%s", i, result)
			}
		}
	}

	return "-1"
}
