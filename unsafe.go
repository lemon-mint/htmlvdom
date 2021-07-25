// +build !appengine

package htmlvdom

func StringToBytes(s string) []byte {
	/*
		sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
		bh := reflect.SliceHeader{
			Data: sh.Data,
			Len:  sh.Len,
			Cap:  sh.Len,
		}
		return *(*[]byte)(unsafe.Pointer(&bh))
	*/
	return []byte(s)
}

func BytesToString(b []byte) string {
	/*
		bh := *(*reflect.SliceHeader)(unsafe.Pointer(&b))
		sh := reflect.StringHeader{
			Data: bh.Data,
			Len:  bh.Len,
		}
		return *(*string)(unsafe.Pointer(&sh))
	*/
	return string(b)
}
