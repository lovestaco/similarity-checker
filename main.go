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
	similarity := float64(64-distance) / 64.0 * 100
	return similarity
}

func interpretRelationship(distance uint8) string {
	switch {
	case distance == 0:
		return "Identical"
	case distance >= 1 && distance <= 3:
		return "Near duplicates"
	case distance >= 4 && distance <= 10:
		return "Minor variants"
	case distance >= 11 && distance <= 25:
		return "Somewhat related"
	case distance >= 26 && distance <= 38:
		return "Unrelated"
	case distance >= 39 && distance <= 64:
		return "Maximally different"
	default:
		return "Unknown"
	}
}

func main() {
	files := []string{
		"https:__hexmos.com_freedevtools_",
		"https:__hexmos.com_freedevtools_png_icons_nodewebkit_nodewebkit-line_",
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
		fmt.Printf("Simhash of file %d: %d\n", i+1, hashes[i])
	}

	distance := simhash.Compare(hashes[0], hashes[1])
	similarity := calculateSimilarityPercentage(distance)
	relationship := interpretRelationship(distance)
	fmt.Printf("Distance between file 1 and file 2: %d\n", distance)
	fmt.Printf("Similarity percentage: %.2f%%\n", similarity)
	fmt.Printf("Relationship: %s\n", relationship)
}