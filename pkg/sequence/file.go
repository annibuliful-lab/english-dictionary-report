package sequence

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func WriteWordFiles(words []string, baseDir string) error {
	for _, word := range words {
		if err := writeWordFile(baseDir, word); err != nil {
			log.Printf("Failed to write %s: %v", word, err)
		}
	}
	return nil
}

func writeWordFile(baseDir, word string) error {
	l1, l2 := string(word[0]), string(word[1])
	dir := filepath.Join(baseDir, l1, l2)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("creating dir %s: %w", dir, err)
	}

	filePath := filepath.Join(dir, word+".txt")
	content := strings.Repeat(word+"\n", 100)
	return os.WriteFile(filePath, []byte(content), 0644)
}

func WriteWordsToTextFile(words []string, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create text file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, word := range words {
		if _, err := writer.WriteString(word + "\n"); err != nil {
			return fmt.Errorf("failed to write to text file: %w", err)
		}
	}
	return writer.Flush()
}
