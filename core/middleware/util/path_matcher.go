package util

import (
	"path"
	"strings"
)

func isWildcardMatch(pattern, pathStr string) bool {
	matched, _ := path.Match(pattern, pathStr)
	return matched
}

func isParamPatternMatch(pattern, pathStr string) bool {
	p1 := strings.Split(pattern, "/")
	p2 := strings.Split(pathStr, "/")
	if len(p1) != len(p2) {
		return false
	}
	for i := range p1 {
		if strings.HasPrefix(p1[i], ":") {
			continue
		}
		if p1[i] != p2[i] {
			return false
		}
	}
	return true
}

func IsSkipPaths(pathStr string, skipPaths []string) bool {
	for _, ignorePath := range skipPaths {
		if pathStr == ignorePath || isWildcardMatch(ignorePath, pathStr) || isParamPatternMatch(ignorePath, pathStr) {
			return true
		}
	}

	return false
}
