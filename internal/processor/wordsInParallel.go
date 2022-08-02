package processor

import (
	"rmertz.com/anagram/internal/model"
)

type processorParallel struct {
	wordsDictionary    []string
	notWordsDictionary []string
	words              []string
}

func NewProcessorParallel(wordsDictionary []string, notWordsDictionary []string, words []string) *processorParallel {
	return &processorParallel{
		wordsDictionary:    wordsDictionary,
		notWordsDictionary: notWordsDictionary,
		words:              words,
	}
}

func worker(wordsDictionary []string, notWordsDictionary []string, words []string, result chan<- model.Result) {
	result <- processPermutations(words, wordsDictionary, notWordsDictionary)
}

func (p *processorParallel) getWords() model.Result {
	chunks := getChunks(p.words)
	channelResult := make(chan model.Result)

	for _, chunk := range chunks {
		go worker(p.wordsDictionary, p.notWordsDictionary, chunk, channelResult)
	}
	return <-channelResult
}
