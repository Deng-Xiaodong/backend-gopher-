package util

import (
	"crypto/sha256"
)

func hash64(x int64, seed int64) int64 {
	x += seed
	x = (x ^ (x >> 33)) * 0xff51afd7ed558cc
	x = (x ^ (x >> 33)) * 0xc4ceb9fe1a85ec5
	x = x ^ (x >> 33)
	return x

}
func BytesToInt64(bytes []byte) int64 {
	var num int64
	for i := 0; i < 8; i++ {
		num <<= 8
		num |= int64(bytes[i])
	}
	return num
}
func Int64ToBytes(num int64) []byte {
	bytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		offset := 64 - 8*(i+1)
		bytes[i] = byte(num >> offset)
	}
	return bytes
}

func GetSha256(in []byte) []byte {
	bytes := sha256.Sum256(in)
	return append([]byte(nil), bytes[:]...)
}

func Xor(a, b []byte) []byte {
	min := 0
	if len(a) > len(b) {
		min = len(b)
	} else {
		min = len(a)
	}
	ret := make([]byte, min)
	for i := 0; i < min; i++ {
		ret[i] = a[i] ^ b[i]
	}
	return ret
}
