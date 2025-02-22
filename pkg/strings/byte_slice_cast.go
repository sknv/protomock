package strings

import "unsafe"

func StringToByteSlice(str string) []byte {
	if str == "" {
		return nil
	}

	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func ByteSliceToString(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}

	return unsafe.String(unsafe.SliceData(bs), len(bs))
}
