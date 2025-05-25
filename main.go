package main

import (
	"english-dictionary-report/pkg"
	"runtime"
	"runtime/pprof"
	"time"

	"fmt"
	"log"
	"os"
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

	words, err := pkg.LoadWordsFromFile(wordFilePath)
	if err != nil {
		log.Fatalf("Failed to read word list: %v", err)
	}

	pkg.AnalyzeAndReport(words)
	capitalFirstLetterWords := pkg.CapitalizeFirstLetter(words)

	if err := pkg.WriteWordsToTextFile(capitalFirstLetterWords, "./output/word_list_capitalized.txt"); err != nil {
		log.Fatalf("Failed to write capitalized words to text file: %v", err)
	}
	fmt.Println("Capitalized word list written to word_list_capitalized.txt")

	if err := pkg.ExportWordsToPDF(words, "./output/word_list_output.pdf"); err != nil {
		log.Fatalf("Failed to export PDF: %v", err)
	}
	fmt.Println("PDF exported to word_list_output.pdf")

	if err := pkg.WriteWordFiles(words, baseDir); err != nil {
		log.Printf("File writing errors: %v", err)
	}

	pkg.ReportFolderSizes(baseDir)
	pkg.ZipTopLevelDirs(baseDir)

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
