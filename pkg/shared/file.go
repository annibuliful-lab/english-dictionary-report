package shared

import (
	"bufio"
	"os"
	"strings"
)

func LoadWordsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if len(word) >= 2 {
			words = append(words, word)
		}
	}
	return words, scanner.Err()
}
