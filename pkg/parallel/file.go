package parallel

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func WriteWordFiles(words []string, baseDir string) error {
	const workers = 16
	jobs := make(chan [2]string, len(words))
	errs := make(chan error, len(words))

	for _, word := range words {
		l1, l2 := string(word[0]), string(word[1])
		dir := filepath.Join(baseDir, l1, l2)
		filePath := filepath.Join(dir, word+".txt")
		jobs <- [2]string{word, filePath}
	}
	close(jobs)

	var sbPool = sync.Pool{
		New: func() any {
			return new(strings.Builder)
		},
	}

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				word, filePath := job[0], job[1]
				dir := filepath.Dir(filePath)
				if err := os.MkdirAll(dir, os.ModePerm); err != nil {
					errs <- err
					continue
				}

				sb := sbPool.Get().(*strings.Builder)
				sb.Reset()
				for j := 0; j < 100; j++ {
					sb.WriteString(word)
					sb.WriteByte('\n')
				}
				err := os.WriteFile(filePath, []byte(sb.String()), 0644)
				sbPool.Put(sb)

				if err != nil {
					errs <- err
				}
			}
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		log.Println(err)
	}
	return nil
}
