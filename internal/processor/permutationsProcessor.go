package processor

import (
	"golang.org/x/exp/slices"
	"rmertz.com/anagram/internal/model"
)

func processPermutations(words []string, realWords []string, unrealWords []string) model.Result {
	foundWords := []string{}
	newWords := []string{}
	newNotWords := []string{}

	for _, word := range words {
		if slices.Contains(unrealWords, word) {
			continue
		} else if slices.Contains(realWords, word) {
			foundWords = append(foundWords, word)
		} else {
			exists := isValidWord(word)
			if exists {
				newWords = append(newWords, word)
				foundWords = append(foundWords, word)
			} else {
				newNotWords = append(newNotWords, word)
			}
		}
	}
	finalResult := model.Result{FoundWords: foundWords, NewWords: newWords, NewNotWords: newNotWords}
	return finalResult
}
