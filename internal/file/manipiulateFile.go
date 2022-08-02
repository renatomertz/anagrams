package file

import (
	"bufio"
	"fmt"
	"os"
)

func openFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		return nil, err
	}
	return file, nil

}

func getFileContent(file *os.File) []string {
	contentFile := []string{}

	var scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		contentFile = append(contentFile, scanner.Text())
	}

	return contentFile
}

func ReadContentFile(fileName string) ([]string, error) {
	file, err := openFile(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content := getFileContent(file)
	return content, nil
}

func WriteContentFile(content []string, fileName string) error {
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
