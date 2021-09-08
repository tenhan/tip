package str

import (
	"bytes"
	"strings"
)

// Trim
func Trim(str string) string {
	return strings.Trim(str, "\n\r ")
}

// RemoveDuplicatedWhiteSpace
func RemoveDuplicatedWhiteSpace(str string) string {
	if len(str) <= 1{
		return str
	}
	buffer := bytes.NewBufferString("")
	s := []rune(str)
	buffer.WriteRune(s[0])
	for i:=1;i<len(s);i++{
		if ! (s[i] == s[i-1] && s[i] == ' '){
			buffer.WriteRune(s[i])
		}
	}
	return buffer.String()
}

// StartWith
func StartWith(s, subStr string) bool {
	if strings.Index(s, subStr) == 0 {
		return true
	}
	return false
}

// EndWith
func EndWith(s, subStr string) bool {
	i := strings.Index(s, subStr)
	if i == -1 {
		return false
	}
	if strings.Index(s, subStr) == len(s)-len(subStr) {
		return true
	}
	return false
}
