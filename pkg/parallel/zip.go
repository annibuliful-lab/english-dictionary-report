package parallel

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func ZipTopLevelDirs(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Fatalf("Failed to read base directory: %v", err)
	}

	fmt.Println("\nZip Size Report:")
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		wg.Add(1)
		go func(entry fs.DirEntry) {
			defer wg.Done()
			l1Name := entry.Name()
			l1Path := filepath.Join(root, l1Name)
			zipPath := filepath.Join(root, l1Name+".zip")

			originalSize := computeDirSize(l1Path)

			if err := createZipArchive(l1Path, root, zipPath); err != nil {
				log.Printf("Failed to zip %s: %v", l1Path, err)
				return
			}

			zipInfo, err := os.Stat(zipPath)
			if err != nil {
				log.Printf("Cannot stat zip file %s: %v", zipPath, err)
				return
			}

			zipSize := zipInfo.Size()
			saved := float64(originalSize-zipSize) / float64(originalSize) * 100

			mu.Lock()
			fmt.Printf("[%s] %8d KB → %8d KB (Saved %.1f%%)\n", l1Name, originalSize/1024, zipSize/1024, saved)
			mu.Unlock()
		}(entry)
	}
	wg.Wait()
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
