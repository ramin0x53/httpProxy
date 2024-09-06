package utility

import "strings"

// Returns include check and exist
func IncludeCheck(txt, subStr string) (bool, bool) {
	return strings.Contains(txt, subStr), subStr != ""
}

// Returns exclude check and exist
func ExcludeCheck(txt, subStr string) (bool, bool) {
	checkResult, exist := IncludeCheck(txt, subStr)
	return !checkResult, exist
}
