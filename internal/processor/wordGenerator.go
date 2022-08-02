package processor

import (
	"strings"

	"golang.org/x/exp/slices"
)

var alreadyRead []string

func GetPossibleWords(groupOfLetterCombinations [][]string) []string {
	words := []string{}
	for _, groupOfLetter := range groupOfLetterCombinations {
		lenStr := len(groupOfLetter)
		if lenStr >= 3 {
			words = generateWords(groupOfLetter, 0, lenStr-1, words)
		}
	}
	return words
}

func generateWords(str []string, startIndex int, endIndex int, words []string) []string {
	if startIndex == endIndex {
		word := strings.Join(str, "")
		if !slices.Contains(alreadyRead, word) {
			alreadyRead = append(alreadyRead, word)
			words = append(words, word)
			return words
		}
	}

	for i := startIndex; i <= endIndex; i++ {
		str = swap(str, startIndex, i)
		words = generateWords(str, startIndex+1, endIndex, words)
		str = swap(str, startIndex, i)
	}
	return words
}

func swap(str []string, index1 int, index2 int) []string {
	temp := str[index1]
	str[index1] = str[index2]
	str[index2] = temp
	return str
}
