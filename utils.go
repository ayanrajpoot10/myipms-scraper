package main

import (
	"sort"
	"strings"
)

// levenshteinDistance calculates the Levenshtein distance between two strings.
func levenshteinDistance(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	// Use only two rows to save memory
	prev, curr := make([]int, lb+1), make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		prev[j] = j
	}

	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min3(
				prev[j]+1,    // deletion
				curr[j-1]+1,  // insertion
				prev[j-1]+cost, // substitution
			)
		}
		prev, curr = curr, prev
	}
	return prev[lb]
}

// min3 returns the minimum of three integers.
func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// findBestMatches finds the best matching strings using Levenshtein distance.
func findBestMatches(input string, options []string, maxResults int) []string {
	inputLower := strings.ToLower(input)
	type match struct {
		value    string
		distance int
	}
	matches := make([]match, 0, len(options))

	for _, option := range options {
		optionLower := strings.ToLower(option)

		if inputLower == optionLower {
			return []string{option} // exact match shortcut
		}
		if strings.Contains(optionLower, inputLower) {
			matches = append(matches, match{option, 0})
			continue
		}

		dist := levenshteinDistance(inputLower, optionLower)
		if dist <= max(len(input), len(option))/2 {
			matches = append(matches, match{option, dist})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].distance < matches[j].distance
	})

	if len(matches) > maxResults {
		matches = matches[:maxResults]
	}

	results := make([]string, len(matches))
	for i, m := range matches {
		results[i] = m.value
	}
	return results
}

// max returns the larger of two ints.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
