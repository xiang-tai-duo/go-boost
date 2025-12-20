//go:build linux

package winbase

func CString(s string) *uint16 {
	return nil
}

func FreeCString(ptr *uint16) {
}

func GoString(lpwsz interface{}) string {
	return ""
}

func IsBadStringPtr(p *uint16) bool {
	return false
}
