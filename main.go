package main

import (
	"archive/zip"
	"bufio"
	"runtime"
	"runtime/pprof"
	"time"

	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/phpdave11/gofpdf"
)

func main() {
	start := time.Now()

	// Enable CPU profiling
	fCPU, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatalf("Failed to create CPU profile: %v", err)
	}
	defer fCPU.Close()
	if err := pprof.StartCPUProfile(fCPU); err != nil {
		log.Fatalf("Failed to start CPU profiling: %v", err)
	}
	defer pprof.StopCPUProfile()

	// Main logic
	wordFilePath, baseDir := parseArgs()

	words, err := loadWords(wordFilePath)
	if err != nil {
		log.Fatalf("Failed to read word list: %v", err)
	}

	analyzeAndReport(words)
	capitalFirstLetterWords := capitalizeFirstLetter(words)

	if err := writeWordsToTextFile(capitalFirstLetterWords, "./output/word_list_capitalized.txt"); err != nil {
		log.Fatalf("Failed to write capitalized words to text file: %v", err)
	}
	fmt.Println("Capitalized word list written to word_list_capitalized.txt")

	if err := exportWordsToPDF(words, "./output/word_list_output.pdf"); err != nil {
		log.Fatalf("Failed to export PDF: %v", err)
	}
	fmt.Println("PDF exported to word_list_output.pdf")

	if err := writeWordFiles(words, baseDir); err != nil {
		log.Printf("File writing errors: %v", err)
	}

	reportFolderSizes(baseDir)
	zipTopLevelDirs(baseDir)

	// Measure runtime
	fmt.Printf("\nExecution time: %s\n", time.Since(start))

	// Measure Heap profiling
	fMem, err := os.Create("mem.prof")
	if err != nil {
		log.Fatalf("Failed to create memory profile: %v", err)
	}
	defer fMem.Close()
	runtime.GC() // run garbage collection before measuring
	if err := pprof.WriteHeapProfile(fMem); err != nil {
		log.Fatalf("Failed to write memory profile: %v", err)
	}

}

func parseArgs() (wordFilePath, baseDir string) {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <word_list_file> [base_output_dir]")
	}
	wordFilePath = os.Args[1]
	baseDir = "output"
	if len(os.Args) >= 3 {
		baseDir = os.Args[2]
	}
	return
}

func loadWords(path string) ([]string, error) {
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

func analyzeAndReport(words []string) {
	countLengthGT5 := countWordsLongerThan(words, 5)
	countRepeatChars := countWordsWithRepeatingChars(words, 2)
	countSameStartEnd := countWordsSameStartEnd(words)

	fmt.Printf("\n7.1 Words longer than 5 characters: %d\n", countLengthGT5)
	fmt.Printf("7.2 Words with ≥2 repeating characters: %d\n", countRepeatChars)
	fmt.Printf("7.3 Words starting and ending with the same letter: %d\n", countSameStartEnd)
}

func writeWordFiles(words []string, baseDir string) error {
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

func reportFolderSizes(root string) {
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

func zipTopLevelDirs(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Fatalf("Failed to read base directory: %v", err)
	}

	fmt.Println("\nZip Size Report:")
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		l1Name := entry.Name()
		l1Path := filepath.Join(root, l1Name)
		zipPath := filepath.Join(root, l1Name+".zip")

		originalSize := computeDirSize(l1Path)

		if err := createZipArchive(l1Path, root, zipPath); err != nil {
			log.Printf("Failed to zip %s: %v", l1Path, err)
			continue
		}

		zipInfo, err := os.Stat(zipPath)
		if err != nil {
			log.Printf("Cannot stat zip file %s: %v", zipPath, err)
			continue
		}

		zipSize := zipInfo.Size()
		saved := float64(originalSize-zipSize) / float64(originalSize) * 100
		fmt.Printf("[%s] %8d KB → %8d KB (Saved %.1f%%)\n", l1Name, originalSize/1024, zipSize/1024, saved)
	}
}

func computeDirSize(dir string) int64 {
	var totalSize int64
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
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

func createZipArchive(srcDir, baseRoot, destZip string) error {
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

// ---------------- Word Analysis ---------------- //

func countWordsLongerThan(words []string, length int) int {
	count := 0
	for _, word := range words {
		if len(word) > length {
			count++
		}
	}
	return count
}

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

func countWordsSameStartEnd(words []string) int {
	count := 0
	for _, word := range words {
		if len(word) >= 1 && word[0] == word[len(word)-1] {
			count++
		}
	}
	return count
}

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

func exportWordsToPDF(words []string, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 14)

	lineHeight := 10.0
	marginTop := 20.0
	marginBottom := 270.0
	y := marginTop

	for _, word := range words {
		if y > marginBottom {
			pdf.AddPage()
			y = marginTop
		}
		pdf.Text(20, y, word)
		y += lineHeight
	}

	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return fmt.Errorf("failed to write PDF: %w", err)
	}
	return nil
}

func writeWordsToTextFile(words []string, outputPath string) error {
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
