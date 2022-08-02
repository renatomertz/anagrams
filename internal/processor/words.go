package processor

import "rmertz.com/anagram/internal/model"

type wordsStrategy interface {
	getWords() model.Result
}
