package main

import (
	"archive/zip"
	"bufio"
	"english-dictionary-report/pkg"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	wordFilePath, baseDir := parseArgs()
	words, err := readWordList(wordFilePath)
	if err != nil {
		log.Fatalf("Failed to read word list: %v", err)
	}

	// Q7.1 – Q7.4
	countLengthGT5 := countWordsLongerThan(words, 5)
	countRepeatChars := countWordsWithRepeatingChars(words, 2)
	countSameStartEnd := countWordsSameStartEnd(words)
	words = capitalizeFirstLetter(words)

	if err := pkg.ExportWordsToPDF(words, "word_list_output.pdf"); err != nil {
		log.Fatalf("Failed to export PDF: %v", err)
	}
	fmt.Println("✔ PDF exported to word_list_output.pdf")

	fmt.Printf("\n7.1 Words longer than 5 characters: %d\n", countLengthGT5)
	fmt.Printf("7.2 Words with ≥2 repeating characters: %d\n", countRepeatChars)
	fmt.Printf("7.3 Words starting and ending with the same letter: %d\n", countSameStartEnd)

	for _, word := range words {
		if err := writeWordFile(baseDir, word); err != nil {
			log.Printf("Failed to write word file: %v", err)
		}
	}

	fmt.Println("\nFolder Size Report (Level 1):")
	reportTopLevelSizes(baseDir)

	fmt.Println("\nZip Size Report:")
	zipLevelOneDirs(baseDir)
}

func parseArgs() (string, string) {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <word_list_file> [base_output_dir]")
	}
	wordFilePath := os.Args[1]
	baseDir := "output"
	if len(os.Args) >= 3 {
		baseDir = os.Args[2]
	}
	return wordFilePath, baseDir
}

func readWordList(path string) ([]string, error) {
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

func writeWordFile(baseDir, word string) error {
	l1, l2 := string(word[0]), string(word[1])
	dir := filepath.Join(baseDir, l1, l2)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("creating dir %s: %w", dir, err)
	}

	content := strings.Repeat(word+"\n", 100)
	filePath := filepath.Join(dir, word+".txt")
	return os.WriteFile(filePath, []byte(content), 0644)
}

func reportTopLevelSizes(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Fatalf("Failed to read base directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		l1Path := filepath.Join(root, entry.Name())
		var totalSize int64

		fmt.Printf("\n[%s/]\n", l1Path)
		filepath.WalkDir(l1Path, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			info, err := os.Stat(path)
			if err == nil {
				fmt.Printf("  %8d KB  %s\n", info.Size()/1024, path)
				totalSize += info.Size()
			}
			return nil
		})

		fmt.Printf("==> Total size: %d KB\n", totalSize/1024)
	}
}

func zipLevelOneDirs(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Fatalf("Failed to read base directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		l1Name := entry.Name()
		l1Path := filepath.Join(root, l1Name)
		zipPath := filepath.Join(root, l1Name+".zip")

		originalSize := computeDirSize(l1Path)

		if err := createZip(l1Path, root, zipPath); err != nil {
			log.Printf("Failed to zip %s: %v", l1Path, err)
			continue
		}

		zipInfo, err := os.Stat(zipPath)
		if err != nil {
			log.Printf("Cannot stat zip file %s: %v", zipPath, err)
			continue
		}

		zipSize := zipInfo.Size()
		diffPercent := float64(originalSize-zipSize) / float64(originalSize) * 100
		fmt.Printf("[%s] %8d KB → %8d KB (Saved %.1f%%)\n", l1Name, originalSize/1024, zipSize/1024, diffPercent)
	}
}

func computeDirSize(dir string) int64 {
	var totalSize int64
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := os.Stat(path)
		if err == nil {
			totalSize += info.Size()
		}
		return nil
	})
	return totalSize
}

func createZip(srcDir, baseRoot, destZip string) error {
	zipFile, err := os.Create(destZip)
	if err != nil {
		return fmt.Errorf("creating zip: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(srcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(baseRoot, path)
		if err != nil {
			return err
		}

		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		reader, err := os.Open(path)
		if err != nil {
			return err
		}
		defer reader.Close()

		_, err = io.Copy(writer, reader)
		return err
	})
}

// 7.1
func countWordsLongerThan(words []string, length int) int {
	count := 0
	for _, word := range words {
		if len(word) > length {
			count++
		}
	}
	return count
}

// 7.2
func countWordsWithRepeatingChars(words []string, minRepeat int) int {
	count := 0
	for _, word := range words {
		charCount := make(map[rune]int)
		for _, ch := range word {
			charCount[ch]++
		}
		for _, v := range charCount {
			if v >= minRepeat {
				count++
				break
			}
		}
	}
	return count
}

// 7.3
func countWordsSameStartEnd(words []string) int {
	count := 0
	for _, word := range words {
		if len(word) >= 1 && word[0] == word[len(word)-1] {
			count++
		}
	}
	return count
}

// 7.4
func capitalizeFirstLetter(words []string) []string {
	newWords := make([]string, len(words))
	for i, word := range words {
		if len(word) == 0 {
			newWords[i] = word
			continue
		}
		newWords[i] = strings.ToUpper(string(word[0])) + word[1:]
	}
	return newWords
}
