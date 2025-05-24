package pkg

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func ZipLevelOneDirs(root string) {
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
