package main

import (
	"sort"
	"strings"
)

// levenshteinDistance calculates the Levenshtein distance between two strings
func levenshteinDistance(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	matrix := make([][]int, len(a)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(b)+1)
		matrix[i][0] = i
	}
	for j := 0; j <= len(b); j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}
	return matrix[len(a)][len(b)]
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

// findBestMatches finds the best matching strings using Levenshtein distance
func findBestMatches(input string, options []string, maxResults int) []string {
	type match struct {
		value    string
		distance int
	}

	var matches []match
	inputLower := strings.ToLower(input)

	for _, option := range options {
		optionLower := strings.ToLower(option)

		if inputLower == optionLower {
			return []string{option}
		}

		if strings.Contains(optionLower, inputLower) {
			matches = append(matches, match{option, 0})
			continue
		}

		distance := levenshteinDistance(inputLower, optionLower)

		maxLen := len(input)
		if len(option) > maxLen {
			maxLen = len(option)
		}

		if distance <= maxLen/2 {
			matches = append(matches, match{option, distance})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].distance < matches[j].distance
	})

	var results []string
	limit := maxResults
	if len(matches) < limit {
		limit = len(matches)
	}

	for i := 0; i < limit; i++ {
		results = append(results, matches[i].value)
	}

	return results
}
