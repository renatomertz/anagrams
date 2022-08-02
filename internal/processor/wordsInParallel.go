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

func worker(wordsDictionary []string, notWordsDictionary []string, work <-chan []string, result chan<- model.Result) {
	result <- processPermutations(<-work, wordsDictionary, notWordsDictionary)
}

func (p *processorParallel) getWords() chan model.Result {
	chunks := getChunks(p.words)
	result := make(chan model.Result)
	jobs := make(chan []string)

	for _, chunk := range chunks {
		go worker(p.wordsDictionary, p.notWordsDictionary, jobs, result)
		jobs <- chunk
	}

	close(jobs)

	return result
}
