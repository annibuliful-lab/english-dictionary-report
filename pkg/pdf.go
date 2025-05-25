package pkg

import (
	"fmt"

	"github.com/phpdave11/gofpdf"
)

func ExportWordsToPDF(words []string, outputPath string) error {
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
