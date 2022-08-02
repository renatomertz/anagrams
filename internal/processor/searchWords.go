package processor

import "rmertz.com/anagram/internal/model"

type searchStrategy struct {
	wordsStrategy wordsStrategy
}

func NewStrategy() *searchStrategy {
	return &searchStrategy{}
}

func (ss *searchStrategy) AddStrategy(ws wordsStrategy) {
	ss.wordsStrategy = ws
}

func (ss *searchStrategy) FindWords() chan model.Result {
	return ss.wordsStrategy.getWords()
}
