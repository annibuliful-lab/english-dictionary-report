package parallel

import (
	"io/fs"
	"path/filepath"
)

func computeDirSize(dir string) int64 {
	var totalSize int64
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if info, err := d.Info(); err == nil {
			totalSize += info.Size()
		}
		return nil
	})
	return totalSize
}
