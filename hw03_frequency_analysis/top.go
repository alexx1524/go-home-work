package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var (
	splitPattern   = regexp.MustCompile(`\S+`)
	replacePattern = regexp.MustCompile(`[!?,.\'\"]`)
)

func Top10(input string) []string {
	dict := make(map[string]int)

	matches := splitPattern.FindAllString(input, -1)
	for _, element := range matches {
		key := replacePattern.ReplaceAllString(strings.ToLower(element), "")
		if key == "-" {
			continue
		}

		dict[key]++
	}

	words := make([]string, 0, len(dict))

	for word := range dict {
		words = append(words, word)
	}

	sort.Slice(words, func(i, j int) bool {
		if dict[words[i]] == dict[words[j]] {
			return words[i] < words[j]
		}
		return dict[words[i]] > dict[words[j]]
	})

	if len(words) > 10 {
		return words[:10]
	}

	return words
}
