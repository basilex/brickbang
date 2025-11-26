package utility

import "strings"

// MatchPattern compares resource and pattern using "*" as a single-segment wildcard.
//
// Rules:
//   - — matches exactly one segment
//     segment count must match, no variable-length patterns
//
// Examples:
//
//	MatchPattern("sys:user:delete", "sys:*:*") → true
//	MatchPattern("sys:user:delete", "sys:user:*") → true
//	MatchPattern("tenant:product:update", "tenant:*:update") → true
//	MatchPattern("tenant:order:update",  "tenant:product:*") → false
//	MatchPattern("tenant:product:read",  "tenant:*:*") → true
//	MatchPattern("tenant:product:read",  "*:product:read") → true
func MatchPattern(resource, pattern string) bool {
	if resource == pattern {
		return true
	}

	resSeg := strings.Split(resource, ":")
	patSeg := strings.Split(pattern, ":")

	if len(resSeg) != len(patSeg) {
		return false
	}

	for i := range resSeg {
		if patSeg[i] == "*" {
			continue
		}

		if patSeg[i] != resSeg[i] {
			return false
		}
	}

	return true
}
