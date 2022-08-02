package processor

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//"http://api.dicionario-aberto.net/word/%s/1"
const uriBase string = "https://www.dicio.com.br/%s"

func isValidWord(word string) bool {
	body := getdWordInictionary(word)
	bodyString := getBodyStringFormat(body)
	return foundWord(bodyString)
}

func foundWord(body string) bool {
	return !strings.Contains(body, "Ops, a palavra que procura não foi encontrada no Dicio.")
}

func getBodyStringFormat(resp *http.Response) string {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	bodyString := string(body)
	return bodyString
}

func getdWordInictionary(word string) *http.Response {
	wordLowerCase := strings.ToLower(word)

	url := fmt.Sprintf(uriBase, wordLowerCase)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}
