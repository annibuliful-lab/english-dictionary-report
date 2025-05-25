package sequence

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func ReportFolderSizes(root string) {
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

func AnalyzeAndReport(words []string) {
	countLengthGT5 := CountWordsLongerThan(words, 5)
	countRepeatChars := CountWordsWithRepeatingChars(words, 2)
	countSameStartEnd := CountWordsSameStartEnd(words)

	fmt.Printf("\n7.1 Words longer than 5 characters: %d\n", countLengthGT5)
	fmt.Printf("7.2 Words with ≥2 repeating characters: %d\n", countRepeatChars)
	fmt.Printf("7.3 Words starting and ending with the same letter: %d\n", countSameStartEnd)
}

func ReportDirSize(dir string) int64 {
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
