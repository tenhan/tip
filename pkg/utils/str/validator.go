package str

// IsAlpha
func IsAlpha(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, v := range str {
		if !((v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z')) {
			return false
		}
	}
	return true
}
