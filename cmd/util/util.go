package util

import (
	"bytes"
)

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}


func DebugError() {

}
