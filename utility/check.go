package utility

import "strings"

func IncludeCheck(txt, subStr string) bool {
	return subStr != "" && strings.Contains(txt, subStr)
}

func ExcludeCheck(txt, subStr string) bool {
	return !IncludeCheck(txt, subStr)
}
