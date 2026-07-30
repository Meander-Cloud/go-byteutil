package zerocopy

import (
	"unsafe"
)

func ByteSliceToString(bs []byte) string {
	bsLen := len(bs)
	if bsLen == 0 {
		return ""
	}

	return unsafe.String(unsafe.SliceData(bs), bsLen)
}

func StringToByteSlice(s string) []byte {
	sLen := len(s)
	if sLen == 0 {
		return nil
	}

	return unsafe.Slice(unsafe.StringData(s), sLen)
}
