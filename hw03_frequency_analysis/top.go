package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var regexpValid = regexp.MustCompile(`^[0-9A-Za-zа-яА-Я_,-]+$`)

func Top10(text string) []string {
	countMap := map[string]int{}
	textArr := strings.Fields(text)
	outputArr := []string{}

	for _, word := range textArr {
		if !regexpValid.MatchString(word) {
			continue
		}
		if countMap[word] == 0 {
			outputArr = append(outputArr, word)
		}
		countMap[word]++
	}

	sort.Slice(outputArr, func(i, j int) bool {
		if countMap[outputArr[i]] == countMap[outputArr[j]] {
			return strings.Compare(outputArr[i], outputArr[j]) < 0
		}
		return countMap[outputArr[i]] > countMap[outputArr[j]]
	})

	max := 10
	if max > len(outputArr) {
		max = len(outputArr)
	}

	return outputArr[0:max]
}
