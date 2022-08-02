package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

type result struct {
	foundWords  []string
	newWords    []string
	newNotWords []string
}

var alreadyRead []string

const path_words_dictionary string = "wordsDictionary.txt"
const path_not_words_dictionary string = "notWordsDictionary.txt"

// func contains(slice []string, str string) bool {
// 	for _, v := range slice {
// 		if v == str {
// 			return true
// 		}
// 	}
// 	return false
// }

func getFileContent(file *os.File) []string {
	contentFile := []string{}

	var scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		contentFile = append(contentFile, scanner.Text())
	}

	return contentFile
}

func openFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		return nil, err
	}
	return file, nil

}

func readContentFile(fileName string) ([]string, error) {
	file, err := openFile(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content := getFileContent(file)
	return content, nil
}

func writeContentFile(content []string, fileName string) error {
	file, err := openFile(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	// Cria um writer responsavel por escrever cada linha do slice no arquivo de texto
	writer := bufio.NewWriter(file)
	for _, line := range content {
		fmt.Fprintln(writer, line)
	}

	// Caso a funcao flush retorne um erro ele sera retornado aqui tambem
	return writer.Flush()
}

// func searchIfWordExists(word string) bool {
// 	go func() {

// 	}()
// 	wordLowerCase := strings.ToLower(word)
// 	//url := fmt.Sprintf("http://api.dicionario-aberto.net/word/%s/1", wordLowerCase)
// 	url := fmt.Sprintf("https://www.dicio.com.br/%s", wordLowerCase)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	//Convert the body to type string
// 	sb := string(body)
// 	return !strings.Contains(sb, "Ops, a palavra que procura nao foi encontrada no Dicio.")
// 	//return len(sb) > 2
// }

func searchIfWordExists(word string, channel chan bool) {
	go func() {
		wordLowerCase := strings.ToLower(word)
		//url := fmt.Sprintf("http://api.dicionario-aberto.net/word/%s/1", wordLowerCase)
		url := fmt.Sprintf("https://www.dicio.com.br/%s", wordLowerCase)
		start := time.Now()
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		elapsed := time.Since(start)
		fmt.Printf("Time to call API: %s\n", elapsed)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		//Convert the body to type string
		sb := string(body)
		channel <- !strings.Contains(sb, "Ops, a palavra que procura nao foi encontrada no Dicio.")
		defer close(channel)
	}()
	//return len(sb) > 2
}

// func processPermutations(words []string, realWords []string, unrealWords []string) <-chan result {
// 	channel := make(chan result)
// 	foundWords := []string{}
// 	newWords := []string{}
// 	newNotWords := []string{}

// 	go func() {
// 		for _, word := range words {
// 			if slices.Contains(unrealWords, word) {
// 				continue
// 			} else if slices.Contains(realWords, word) {
// 				foundWords = append(foundWords, word)
// 			} else if searchIfWordExists(word) {
// 				newWords = append(newWords, word)
// 				foundWords = append(foundWords, word)
// 			} else {
// 				newNotWords = append(newNotWords, word)
// 			}
// 		}
// 		finalResult := result{foundWords, newWords, newNotWords}
// 		channel <- finalResult
// 		defer close(channel)
// 	}()

// 	return channel
// }

func processPermutations(words []string, realWords []string, unrealWords []string) <-chan result {
	channelResult := make(chan result)
	foundWords := []string{}
	newWords := []string{}
	newNotWords := []string{}

	go func() {
		for _, word := range words {
			if slices.Contains(unrealWords, word) {
				continue
			} else if slices.Contains(realWords, word) {
				foundWords = append(foundWords, word)
			} else {
				exists := make(chan bool)
				go searchIfWordExists(word, exists)
				result := <-exists
				defer close(exists)
				if result {
					newWords = append(newWords, word)
					foundWords = append(foundWords, word)
				} else {
					newNotWords = append(newNotWords, word)
				}
			}

		}
		finalResult := result{foundWords, newWords, newNotWords}
		channelResult <- finalResult
		defer close(channelResult)
	}()

	return channelResult
}

func permute(str []string, startIndex int, endIndex int, words []string) []string {
	if startIndex == endIndex {
		word := strings.Join(str, "")
		if !slices.Contains(alreadyRead, word) {
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

// func permute2(str []string, startIndex int, endIndex int, channel chan []string) {
// 	if startIndex == endIndex {
// 		word := strings.Join(str, "")
// 		if !slices.Contains(alreadyRead, word) {
// 			alreadyRead = append(alreadyRead, word)
// 			//if existentWord(word) {
// 			//fmt.Println(word)
// 			words := <-channel
// 			fmt.Println(words)
// 			words = append(words, word)
// 			channel <- words
// 			//return
// 			//}
// 		}
// 	}

// 	for i := startIndex; i <= endIndex; i++ {
// 		str = swap(str, startIndex, i)
// 		//channel <- words
// 		go permute2(str, startIndex+1, endIndex, channel)
// 		str = swap(str, startIndex, i)
// 	}
// 	//return
// }

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
	start := time.Now()

	// teste := searchIfWordExists("cdcdc")
	// fmt.Println(teste)

	wordsDictionary, err := readContentFile(path_words_dictionary)
	if err != nil {
		log.Fatalln(err)
	}

	notWordsDictionary, err := readContentFile(path_not_words_dictionary)
	if err != nil {
		log.Fatalln(err)
	}

	// fmt.Println(wordsDictionary)
	// fmt.Println(notWordsDictionary)

	str := []string{"N", "D", "E", "O", "C", "B"}
	//{"A", "U", "B", "B", "M", "A"}
	//{"A", "E", "M", "D", "A", "G"}
	//{"C", "F", "A", "H", "I", "E"}
	//{"C", "G", "O", "O", "T", "I", "N"}
	//str := []string{"A", "I", "C", "R"}
	//words := []string{}
	subSets := findSubsets([][]string{}, str, []string{}, 0)
	//channel := make(chan []string)
	//defer close(channel)
	words := []string{}
	for _, subSet := range subSets {
		lenStr := len(subSet)
		if lenStr >= 3 {
			words = permute(subSet, 0, lenStr-1, words)
		}
	}
	// words := <-channel
	// fmt.Println(words)

	fmt.Printf("Permutations to test: %d\n", len(words))
	channel := processPermutations(words, wordsDictionary, notWordsDictionary)

	// result := <-channel
	// fmt.Println(result)
	for result := range channel {
		fmt.Println(result.foundWords)
		err := writeContentFile(result.newWords, path_words_dictionary)
		if err != nil {
			log.Fatalf("Erro:", err)
		}
		err2 := writeContentFile(result.newNotWords, path_not_words_dictionary)
		if err != nil {
			log.Fatalf("Erro:", err2)
		}
	}
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
