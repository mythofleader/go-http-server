package util

import (
	"testing"
)

func TestIsWildcardMatch(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		pathStr  string
		expected bool
	}{
		{"exact match", "abc", "abc", true},
		{"wildcard match", "a*", "abc", true},
		{"wildcard no match", "a*", "xbc", false},
		{"wildcard in middle", "a*c", "abc", true},
		{"wildcard middle no match", "a*c", "axyz", false},
		{"wildcard multiple", "a*b*c", "axyzbc", true},
		{"wildcard empty pattern", "", "abc", false},
		{"wildcard empty path", "a*", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isWildcardMatch(tt.pattern, tt.pathStr)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsParamPatternMatch(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		pathStr  string
		expected bool
	}{
		{"exact match", "abc/def", "abc/def", true},
		{"param match", "abc/:id", "abc/123", true},
		{"param mismatch", "abc/:id", "abc", false},
		{"multiple params match", "abc/:id/:name", "abc/123/john", true},
		{"multiple params no match", "abc/:id/:name", "abc/123", false},
		{"param with static mismatch", "abc/:id", "xyz/123", false},
		{"empty pattern", "", "abc/123", false},
		{"empty path", "abc/:id", "", false},
		{"mismatch length", "abc/def", "abc/def/xyz", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isParamPatternMatch(tt.pattern, tt.pathStr)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
