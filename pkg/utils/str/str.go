package str

import "strings"

// Trim
func Trim(str string) string {
	return strings.Trim(str, "\n\r ")
}

// RemoveDuplicatedWhiteSpace
func RemoveDuplicatedWhiteSpace(str string) string {
	lines := strings.Split(str, "\n")
	var results []string
	for _, v := range lines {
		line := Trim(v)
		if line != "" {
			results = append(results, line)
		}
	}
	return strings.Join(results, "\n")
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
