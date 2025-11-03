package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/mfonda/simhash"
)

func readFile(filename string) ([]byte, error) {
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func removeHeaderFooter(content []byte) []byte {
	// Remove head section: from <head...> to </head>
	headRegex := regexp.MustCompile(`(?s)<head[^>]*>.*?</head>`)
	content = headRegex.ReplaceAll(content, []byte(""))

	// Remove all script tags: from <script...> to </script>
	scriptRegex := regexp.MustCompile(`(?s)<script[^>]*>.*?</script>`)
	content = scriptRegex.ReplaceAll(content, []byte(""))

	// Remove all style tags: from <style...> to </style>
	styleRegex := regexp.MustCompile(`(?s)<style[^>]*>.*?</style>`)
	content = styleRegex.ReplaceAll(content, []byte(""))

	// Remove div with id="ad-banner"
	adBannerRegex := regexp.MustCompile(`(?s)<div[^>]*id="ad-banner"[^>]*>.*?</div>`)
	content = adBannerRegex.ReplaceAll(content, []byte(""))

	// Remove header section: from <header...> to </header>
	headerRegex := regexp.MustCompile(`(?s)<header[^>]*>.*?</header>`)
	content = headerRegex.ReplaceAll(content, []byte(""))

	// Remove footer section: from <footer...> to </footer>
	footerRegex := regexp.MustCompile(`(?s)<footer[^>]*>.*?</footer>`)
	content = footerRegex.ReplaceAll(content, []byte(""))

	return content
}

func calculateSimilarityPercentage(distance uint8) float64 {
	const MAXIMUM = 1000.0
	percentage := 100.0 - ((float64(distance) / MAXIMUM) * 100.0)
	return percentage
}

func main() {
	files := []string{
		"https:__hexmos.com_freedevtools_png_icons_nodewebkit_nodewebkit-line_",
		"https:__hexmos.com_freedevtools_",
	}

	// Create output folder with current date-time
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	outputDir := filepath.Join("output", timestamp)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	var docs [][]byte

	for _, filename := range files {
		fmt.Printf("Reading %s...\n", filename)
		body, err := readFile(filename)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", filename, err)
			continue
		}

		// Remove header and footer
		cleaned := removeHeaderFooter(body)

		// Save cleaned file to output directory
		outputFilename := filepath.Join(outputDir, filename)
		err = os.WriteFile(outputFilename, cleaned, 0644)
		if err != nil {
			fmt.Printf("Error writing cleaned file: %v\n", err)
			continue
		}
		docs = append(docs, cleaned)
	}

	if len(docs) < 2 {
		fmt.Println("Error: Could not process both files")
		return
	}

	hashes := make([]uint64, len(docs))
	for i, d := range docs {
		hashes[i] = simhash.Simhash(simhash.NewWordFeatureSet(d))
		fmt.Printf("Simhash of file %d: %x\n", i+1, hashes[i])
	}

	distance := simhash.Compare(hashes[0], hashes[1])
	similarity := calculateSimilarityPercentage(distance)
	fmt.Printf("Distance between file 1 and file 2: %d\n", distance)
	fmt.Printf("Similarity percentage: %.2f%%\n", similarity)
}