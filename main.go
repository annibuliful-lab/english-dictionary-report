package main

import (
	"english-dictionary-report/pkg/parallel"
	"english-dictionary-report/pkg/sequence"
	"english-dictionary-report/pkg/shared"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

func main() {

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

	if version == "v2" {
		v2Execute(wordFilePath, baseDir)
	} else if version == "v1" {
		V1Execute(wordFilePath, baseDir)
	} else {
		os.Exit(0)
	}

}

func V1Execute(wordFilePath string, baseDir string) {
	words, err := shared.LoadWordsFromFile(wordFilePath)
	if err != nil {
		log.Fatalf("Failed to read word list: %v", err)
	}

	words = shared.CapitalizeFirstLetter(words)

	if err := shared.ExportWordsToPDF(words, fmt.Sprintf("%s/word_list_output.pdf", baseDir)); err != nil {
		log.Fatalf("Failed to export PDF: %v", err)
	}
	fmt.Println("PDF exported to word_list_output.pdf")

	if err := parallel.WriteWordFiles(words, baseDir); err != nil {
		log.Printf("File writing errors: %v", err)
	}

	shared.ReportFolderSizes(baseDir)
	sequence.ZipTopLevelDirs(baseDir)

	shared.AnalyzeAndReport(words)
}

func v2Execute(wordFilePath string, baseDir string) {
	words, err := shared.LoadWordsFromFile(wordFilePath)
	if err != nil {
		log.Fatalf("Failed to read word list: %v", err)
	}

	if err := parallel.WriteWordFiles(words, baseDir); err != nil {
		log.Printf("File writing errors: %v", err)
	}

	if err := shared.ExportWordsToPDF(words, fmt.Sprintf("%s/word_list_output.pdf", baseDir)); err != nil {
		log.Fatalf("Failed to export PDF: %v", err)
	}
	fmt.Println("PDF exported to word_list_output.pdf")

	words = shared.CapitalizeFirstLetter(words)

	shared.ReportFolderSizes(baseDir)
	parallel.ZipTopLevelDirs(baseDir)

	shared.AnalyzeAndReport(words)
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
