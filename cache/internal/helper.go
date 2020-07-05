package internal

import (
	"crypto/md5"
	"encoding/hex"
	"unsafe"
)

func BytesToString(b []byte) string {
	return string(b)
}

func StringToBytes(s string) []byte {
	return []byte(s)
}

func BytesToStringUnsafe(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToBytesUnsafe(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
