package processor

const number_of_chunks int = 5

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
