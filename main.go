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
	"rmertz.com/anagram/internal/service"
)

type result struct {
	foundWords  []string
	newWords    []string
	newNotWords []string
}

var alreadyRead []string

const path_words_dictionary string = "wordsDictionary.txt"
const path_not_words_dictionary string = "notWordsDictionary.txt"
const number_of_chunks int = 5

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
	writer := bufio.NewWriter(file)
	for _, line := range content {
		fmt.Fprintln(writer, line)
	}

	return writer.Flush()
}

func searchIfWordExists(word string) bool {
	wordLowerCase := strings.ToLower(word)
	//url := fmt.Sprintf("http://api.dicionario-aberto.net/word/%s/1", wordLowerCase)
	url := fmt.Sprintf("https://www.dicio.com.br/%s", wordLowerCase)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	return !strings.Contains(sb, "Ops, a palavra que procura não foi encontrada no Dicio.")
	//return len(sb) > 2
}

func processPermutations(words []string, realWords []string, unrealWords []string) result {
	foundWords := []string{}
	newWords := []string{}
	newNotWords := []string{}

	for _, word := range words {
		if slices.Contains(unrealWords, word) {
			continue
		} else if slices.Contains(realWords, word) {
			foundWords = append(foundWords, word)
		} else {
			exists := searchIfWordExists(word)
			if exists {
				newWords = append(newWords, word)
				foundWords = append(foundWords, word)
			} else {
				newNotWords = append(newNotWords, word)
			}
		}
	}
	finalResult := result{foundWords, newWords, newNotWords}
	return finalResult
}

func permute(str []string, startIndex int, endIndex int, words []string) []string {
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

func generate_combinations(subset [][]string, strIn []string, strOutput []string, index int) [][]string {
	if index == len(strIn) {
		subset = append(subset, strOutput)
		return subset
	}
	internalSubset := generate_combinations(subset, strIn, strOutput, index+1)
	strOutput = append(strOutput, strIn[index])
	return generate_combinations(internalSubset, strIn, strOutput, index+1)
}

func getChunks(words []string) [][]string {
	var divided [][]string
	wordsLen := len(words)
	chunkSize := (wordsLen + number_of_chunks - 1) / number_of_chunks

	for i := 0; i < wordsLen; i += chunkSize {
		end := i + chunkSize

		if end > wordsLen {
			end = wordsLen
		}

		divided = append(divided, words[i:end])
	}
	return divided
}

func parallelProcessingPermutations(words []string, wordsDictionary []string, notWordsDictionary []string) result {
	chunks := getChunks(words)
	channelResult := make(chan result)
	for _, chunk := range chunks {
		go worker(wordsDictionary, notWordsDictionary, chunk, channelResult)
	}
	return <-channelResult
}

func worker(wordsDictionary []string, notWordsDictionary []string, words []string, result chan<- result) {
	result <- processPermutations(words, wordsDictionary, notWordsDictionary)
}

func main() {
	start := time.Now()

	// wordsDictionary, err := readContentFile(path_words_dictionary)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// notWordsDictionary, err := readContentFile(path_not_words_dictionary)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// str := []string{"E", "J", "Q"} //, "O", "B", "I"} //, "U"}
	// subSets := generate_combinations([][]string{}, str, []string{}, 0)
	// words := []string{}
	// for _, subSet := range subSets {
	// 	lenStr := len(subSet)
	// 	if lenStr >= 3 {
	// 		words = permute(subSet, 0, lenStr-1, words)
	// 	}
	// }

	// fmt.Printf("Permutations to test: %d\n", len(words))
	// result := parallelProcessingPermutations(words, wordsDictionary, notWordsDictionary)
	// fmt.Println(result.foundWords)
	// errW := writeContentFile(result.newWords, path_words_dictionary)
	// if errW != nil {
	// 	log.Fatalf("Erro:", errW)
	// }
	// errR := writeContentFile(result.newNotWords, path_not_words_dictionary)
	// if err != nil {
	// 	log.Fatalf("Erro:", errR)
	// }
	in := []string{"O", "J", "T", "I", "B", "O", "E"}
	service.ProcessParallelVersion(in)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
