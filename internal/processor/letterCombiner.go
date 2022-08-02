package processor

func GenerateGroupOfLetterCombinations(subset [][]string, strIn []string, strOutput []string, index int) [][]string {
	if index == len(strIn) {
		subset = append(subset, strOutput)
		return subset
	}
	internalSubset := GenerateGroupOfLetterCombinations(subset, strIn, strOutput, index+1)
	strOutput = append(strOutput, strIn[index])
	return GenerateGroupOfLetterCombinations(internalSubset, strIn, strOutput, index+1)
}
