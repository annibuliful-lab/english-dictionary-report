package parallel

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

func LoadWordsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	numWorkers := runtime.NumCPU()
	lines := make(chan string, numWorkers*2)
	results := make(chan string, numWorkers*2)

	var wg sync.WaitGroup

	// Workers: process and print
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lines {
				word := strings.ToLower(strings.TrimSpace(line))
				if len(word) >= 2 {
					// fmt.Println(word) // Print each word immediately
					results <- word
				}
			}
		}()
	}

	// Feed lines
	go func() {
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	// Close results after all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	words := make([]string, 0, 100000)
	for word := range results {
		words = append(words, word)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
