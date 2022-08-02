package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var alreadyRead []string

const path_real_words_dictionary string = "realWordsDictionary.txt"
const path_unreal_words_dictionary string = "unrealWordsDictionary.txt"

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func openDictionaries() ([]string, []string, error) {
	realWordsDictionary, err := os.Open(path_real_words_dictionary)

	if err != nil {
		return nil, nil, err
	}

	defer realWordsDictionary.Close()

	realWords := []string{}

	var scanner = bufio.NewScanner(realWordsDictionary)

	for scanner.Scan() {
		realWords = append(realWords, scanner.Text())
	}

	unrealWordsDictionary, err := os.Open(path_unreal_words_dictionary)

	if err != nil {
		return nil, nil, err
	}

	defer unrealWordsDictionary.Close()

	unrealWords := []string{}

	scanner = bufio.NewScanner(unrealWordsDictionary)

	for scanner.Scan() {
		unrealWords = append(unrealWords, scanner.Text())
	}

	return realWords, unrealWords, nil

}

func existentWord(word string) bool {
	wordLowerCase := strings.ToLower(word)
	url := fmt.Sprintf("http://api.dicionario-aberto.net/word/%s/1", wordLowerCase)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	return len(sb) > 2

}

// func permute(str []string, startIndex int, endIndex int) {
// 	if startIndex == endIndex {
// 		word := strings.Join(str, "")
// 		if !contains(alreadyRead, word) {
// 			alreadyRead = append(alreadyRead, word)
// 			if existentWord(word) {
// 				fmt.Println(word)
// 			}
// 		}
// 	}

// 	for i := startIndex; i <= endIndex; i++ {
// 		str = swap(str, startIndex, i)
// 		permute(str, startIndex+1, endIndex)
// 		str = swap(str, startIndex, i)
// 	}
// }

func processPermutations(words []string, realWords []string, unrealWords []string) ([]string, []string, []string) {
	foundWords := []string{}
	newWords := []string{}
	newNotWords := []string{}
	for _, word := range words {
		if contains(unrealWords, word) {
			continue
		} else if contains(realWords, word) {
			foundWords = append(foundWords, word)
		} else if existentWord(word) {
			newWords = append(newWords, word)
			foundWords = append(foundWords, word)
		} else {
			newNotWords = append(newNotWords, word)
		}
	}
	return foundWords, newWords, newNotWords
}

func permute(str []string, startIndex int, endIndex int, words []string) []string {
	if startIndex == endIndex {
		word := strings.Join(str, "")
		if !contains(alreadyRead, word) {
			alreadyRead = append(alreadyRead, word)
			//if existentWord(word) {
			//fmt.Println(word)
			words = append(words, word)
			return words
			//}
		}
	}

	for i := startIndex; i <= endIndex; i++ {
		str = swap(str, startIndex, i)
		words = permute(str, startIndex+1, endIndex, words)
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

func findSubsets(subset [][]string, strIn []string, strOutput []string, index int) [][]string {
	if index == len(strIn) {
		subset = append(subset, strOutput)
		return subset
	}
	internalSubset := findSubsets(subset, strIn, strOutput, index+1)
	strOutput = append(strOutput, strIn[index])
	return findSubsets(internalSubset, strIn, strOutput, index+1)
}

func main() {
	realWordsDictionary, unrealWordsDictionary, err := openDictionaries()

	if err != nil {
		log.Fatalln(err)
	}

	str := []string{"V", "I", "O"} //, "E", "D"} //, "S"}
	//str := []string{"A", "I", "C", "R"}
	//lenStr := len(str)
	//permute(str, 0, lenStr-1)
	words := []string{}
	subSets := findSubsets([][]string{}, str, []string{}, 0)
	for _, subSet := range subSets {
		lenStr := len(subSet)
		if lenStr >= 3 {
			words = permute(subSet, 0, lenStr-1, words)
		}
	}
	foundWords, newWords, notWords := processPermutations(words, realWordsDictionary, unrealWordsDictionary)

	fmt.Println(foundWords)
	fmt.Println(newWords)
	fmt.Println(notWords)

}
