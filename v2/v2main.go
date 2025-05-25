package v2

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func Main() {
	wordFilePath, baseDir := parseArgs()
	words, err := readWordList(wordFilePath)
	if err != nil {
		log.Fatalf("Failed to read word list: %v", err)
	}

	parallelWriteWordFiles(baseDir, words)

	fmt.Println("Folder Size Report (Level 1):")
	reportTopLevelSizes(baseDir)

	fmt.Println("\n  Zip Size Report:")
	parallelZipLevelOneDirs(baseDir)
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

func parallelWriteWordFiles(baseDir string, words []string) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 16) // limit to 16 concurrent goroutines

	for _, word := range words {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(word string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			if err := writeWordFile(baseDir, word); err != nil {
				log.Printf("Error writing file for word %s: %v", word, err)
			}
		}(word)
	}

	wg.Wait()
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

func parallelZipLevelOneDirs(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Fatalf("Failed to read base directory: %v", err)
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 4) // limit to 4 concurrent zips

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}
		go func(entry fs.DirEntry) {
			defer wg.Done()
			defer func() { <-semaphore }()

			l1Name := entry.Name()
			l1Path := filepath.Join(root, l1Name)
			zipPath := filepath.Join(root, l1Name+".zip")

			originalSize := computeDirSize(l1Path)

			if err := createZip(l1Path, root, zipPath); err != nil {
				log.Printf("Failed to zip %s: %v", l1Path, err)
				return
			}

			zipInfo, err := os.Stat(zipPath)
			if err != nil {
				log.Printf("Cannot stat zip file %s: %v", zipPath, err)
				return
			}

			zipSize := zipInfo.Size()
			diffPercent := float64(originalSize-zipSize) / float64(originalSize) * 100
			fmt.Printf("[%s] %8d KB → %8d KB (Saved %.1f%%)\n", l1Name, originalSize/1024, zipSize/1024, diffPercent)
		}(entry)
	}

	wg.Wait()
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
