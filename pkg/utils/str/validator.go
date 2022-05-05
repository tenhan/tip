package str

import (
	"regexp"
	"strings"
)

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

func IsEnglish(str string) bool {
	str = strings.Trim(str, " \n\r")
	if len(str) == 0 {
		return false
	}
	// only contains letters, " ", "'", "-", ",", "?", "!"
	reg := `^[a-zA-Z][a-zA-Z-!?',\. ]+$`
	ep := regexp.MustCompile(reg)
	return ep.MatchString(str)
}
