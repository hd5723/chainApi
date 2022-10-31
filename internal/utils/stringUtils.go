package utils

import "strings"

func Substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	if end < length {
		length = end
	}
	var substring = ""
	for i := start; i < end; i++ {
		substring += string(r[i])
	}
	return substring
}

func ExistsSpecialLetters(source string) bool {
	if strings.Contains(source, "_") || strings.Contains(source, ";") || strings.Contains(source, "/") ||
		strings.Contains(source, "-") || strings.Contains(source, "~") || strings.Contains(source, "!") ||
		strings.Contains(source, "～") || strings.Contains(source, "！") || strings.Contains(source, "¥") ||
		strings.Contains(source, "@") || strings.Contains(source, "#") || strings.Contains(source, "$") ||
		strings.Contains(source, "%") || strings.Contains(source, "^") || strings.Contains(source, "&") ||
		strings.Contains(source, "*") || strings.Contains(source, "(") || strings.Contains(source, ")") ||
		strings.Contains(source, "（") || strings.Contains(source, "）") || strings.Contains(source, "=") ||
		strings.Contains(source, "+") || strings.Contains(source, " ") ||
		strings.Contains(source, "{") || strings.Contains(source, "}") ||
		strings.Contains(source, "[") || strings.Contains(source, "]") ||
		strings.Contains(source, "【") || strings.Contains(source, "】") ||
		strings.Contains(source, "？") || strings.Contains(source, "《") || strings.Contains(source, "》") ||
		strings.Contains(source, "?") || strings.Contains(source, "<") || strings.Contains(source, ">") {
		return true
	}
	return false
}
