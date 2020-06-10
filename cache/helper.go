package cache

import "unsafe"

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