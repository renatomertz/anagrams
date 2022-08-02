package service

import (
	"fmt"
	"log"

	"rmertz.com/anagram/internal/file"
	"rmertz.com/anagram/internal/processor"
)

const path_words_dictionary string = "resources/wordsDictionary.txt"
const path_not_words_dictionary string = "resources/notWordsDictionary.txt"

func ProcessParallelVersion(str []string) {
	wordsDictionary, err := file.ReadContentFile(path_words_dictionary)
	if err != nil {
		log.Fatalln(err)
	}

	notWordsDictionary, err := file.ReadContentFile(path_not_words_dictionary)
	if err != nil {
		log.Fatalln(err)
	}

	groupOfLetterCombinations := processor.GenerateGroupOfLetterCombinations([][]string{}, str, []string{}, 0)
	possibleWords := processor.GetPossibleWords(groupOfLetterCombinations)

	fmt.Printf("Permutations to test: %d\n", len(possibleWords))

	wordsInParallel := processor.NewProcessorParallel(wordsDictionary, notWordsDictionary, possibleWords)

	searchWords := processor.NewStrategy()
	searchWords.AddStrategy(wordsInParallel)
	//result := searchWords.FindWords()
	channel := searchWords.FindWords()

	//channel := processor.ParallelProcessingPermutations(wordsDictionary, notWordsDictionary, possibleWords)
	for result := range channel {
		fmt.Println(result.FoundWords)
		errW := file.WriteContentFile(result.NewWords, path_words_dictionary)
		if errW != nil {
			log.Fatalf("Erro:", errW)
		}
		errR := file.WriteContentFile(result.NewNotWords, path_not_words_dictionary)
		if err != nil {
			log.Fatalf("Erro:", errR)
		}
	}
}
