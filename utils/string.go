package utils

import "unsafe"

func BytesToString(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}
