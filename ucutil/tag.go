package ucutil

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// 从字符串中解析标签，符合 #xxx 的属于标签
func ParseTags(str string) []string {
	const PATTERN string = "(#\\S+\\s)|(#\\S+$)"
	rgx, err := regexp.Compile(PATTERN)
	if err != nil {
		fmt.Printf("compile regex failed: %s\n", err.Error())
		return make([]string, 0)
	}
	indexArr := rgx.FindAllStringIndex(str, -1)
	resultArr := make([]string, 0)
	for _, idxElem := range indexArr {
		startIndex := getStartIndex(str, idxElem)
		if startIndex < 0 {
			continue
		}
		tagName := TrimTags(str[startIndex:idxElem[1]])
		if tagName == "" {
			continue
		}
		resultArr = append(resultArr, tagName)
	}
	return resultArr
}

func getStartIndex(str string, idxElem []int) int {
	if len(idxElem) < 2 {
		return -1
	}
	if idxElem[0] != 0 {
		lastCh, _ := utf8.DecodeLastRuneInString(str[:idxElem[0]])
		if !unicode.IsSpace(lastCh) {
			return -1
		}
	}

	segment := strings.TrimSpace(str[idxElem[0]:idxElem[1]])
	lastSharpIndex := -1
	for i, ch := range segment {
		if ch == '#' {
			lastSharpIndex = i
		}
	}
	if lastSharpIndex < 0 || lastSharpIndex >= len(segment)-1 {
		return -1
	}
	return idxElem[0] + lastSharpIndex
}

func TrimTags(str string) string {
	tagName := strings.TrimSpace(str)
	if len(tagName) <= 0 {
		return ""
	}
	return strings.TrimLeft(tagName, "#")
}
