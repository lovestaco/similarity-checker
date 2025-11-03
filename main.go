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

func main() {
	files := []string{
		"https___hexmos.com_freedevtools_png_icons_nodewebkit_nodewebkit-plain_",
		"https___hexmos.com_freedevtools_svg_icons_nodewebkit_nodewebkit-plain_",
	}

	// Create output folder with current date-time
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	outputDir := filepath.Join("output", timestamp)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}
	fmt.Printf("Created output directory: %s\n", outputDir)

	var docs [][]byte

	for _, filename := range files {
		fmt.Printf("Reading %s...\n", filename)
		body, err := readFile(filename)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", filename, err)
			continue
		}
		fmt.Printf("Read %d bytes\n", len(body))

		// Remove header and footer
		cleaned := removeHeaderFooter(body)
		fmt.Printf("After cleaning: %d bytes\n", len(cleaned))

		// Save cleaned file to output directory
		outputFilename := filepath.Join(outputDir, filename)
		err = os.WriteFile(outputFilename, cleaned, 0644)
		if err != nil {
			fmt.Printf("Error writing cleaned file: %v\n", err)
			continue
		}
		fmt.Printf("Saved cleaned file to: %s\n", outputFilename)
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

	comparison := simhash.Compare(hashes[0], hashes[1])
	fmt.Printf("Comparison of file 1 and file 2: %d (lower is more similar)\n", comparison)
}