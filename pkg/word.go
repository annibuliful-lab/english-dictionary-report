package pkg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadWordList(path string) ([]string, error) {
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

func WriteWordFile(baseDir, word string) error {
	l1, l2 := string(word[0]), string(word[1])
	dir := filepath.Join(baseDir, l1, l2)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("creating dir %s: %w", dir, err)
	}

	content := strings.Repeat(word+"\n", 100)
	filePath := filepath.Join(dir, word+".txt")
	return os.WriteFile(filePath, []byte(content), 0644)
}
