package main

import (
	"english-dictionary-report/pkg/parallel"
	"english-dictionary-report/pkg/sequence"
	"english-dictionary-report/pkg/shared"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

func main() {

	os.RemoveAll("./v1-output")
	os.RemoveAll("./v2-output")

	version, wordFilePath, baseDir := parseArgs()

	// Optional CPU profiling
	if f, err := os.Create(fmt.Sprintf("%s/cpu.prof", baseDir)); err == nil {
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		defer func() {
			mf, _ := os.Create(fmt.Sprintf("%s/memory.prof", baseDir))
			pprof.WriteHeapProfile(mf)
			mf.Close()
		}()
	}

	switch version {
	case "v1":
		execute("v1", shared.LoadWordsFromFile, sequence.WriteWordFiles, sequence.ZipTopLevelDirs, wordFilePath, baseDir)
	case "v2":
		execute("v2", parallel.LoadWordsFromFile, parallel.WriteWordFiles, parallel.ZipTopLevelDirs, wordFilePath, baseDir)
	default:
		log.Fatalf("Unknown version: %s", version)
	}
}

func execute(
	version string,
	loadFunc func(string) ([]string, error),
	writeFunc func([]string, string) error,
	zipFunc func(string),
	wordFilePath, baseDir string,
) {
	start := time.Now()

	words, err := loadFunc(wordFilePath)
	if err != nil {
		log.Fatalf("Failed to read word list: %v", err)
	}

	log.Println("Loaded words")

	if err := writeFunc(words, baseDir); err != nil {
		log.Printf("File writing errors: %v", err)
	}

	log.Println("TXT file to text")

	if err := shared.ExportWordsToPDF(words, fmt.Sprintf("%s/word_list_output.pdf", baseDir)); err != nil {
		log.Fatalf("Failed to export PDF: %v", err)
	}
	fmt.Println("PDF exported to word_list_output.pdf")

	words = shared.CapitalizeFirstLetter(words)

	if err := sequence.WriteWordsToTextFile(words, fmt.Sprintf("%s/capital-first-letter.txt", baseDir)); err != nil {
		log.Printf("Failed to write capital-first-letter.txt: %v", err)
	}

	shared.ReportFolderSizes(baseDir)
	zipFunc(baseDir)

	shared.AnalyzeAndReport(words)

	duration := time.Since(start)

	fmt.Println("Execution Duraction Time: ", duration)
	shared.WriteDurationCSV(version, duration, "./")
}

func parseArgs() (version, wordFilePath, baseDir string) {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <v1|v2> <word_list_file> [base_output_dir]")
	}
	version = os.Args[1]
	wordFilePath = os.Args[2]
	baseDir = "output"
	if len(os.Args) >= 4 {
		baseDir = os.Args[3]
	}

	baseDir = fmt.Sprintf("%s-%s", version, baseDir)

	return
}
