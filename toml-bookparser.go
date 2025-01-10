package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

// Book struct to hold CSV data
type Book struct {
	Title       string
	ReadState   string
	AuthorName  string
	PubYear     string
	Recommender string
	Rating      string
	Date        string
}

// Function to sanitize file names
func sanitizeFileName(name string) string {
	replacer := strings.NewReplacer(" ", "_", "/", "-", ":", "", "'", "")
	return replacer.Replace(name)
}

// Function to create Markdown content with TOML frontmatter
func generateMarkdownContent(book Book) string {
	// Start TOML frontmatter
	content := "+++\n"
	content += fmt.Sprintf("title = \"%s\"\n", book.Title)
	content += fmt.Sprintf("date = %s\n", book.Date)
	content += "# if you don't use a taxonomy, delete it\n"
	content += "# empty fields not allowed\n"
	content += "[taxonomies]\n"
	content += fmt.Sprintf("  readstate = [\"%s\"]\n", book.ReadState)
	content += fmt.Sprintf("  authorname = [\"%s\"]\n", book.AuthorName)
	content += fmt.Sprintf("  pubyear = [\"%s\"]\n", book.PubYear)

	// Optional fields
	if book.Rating != "NR" && book.Rating != "" {
		content += fmt.Sprintf("  rating = [\"%s\"]\n", book.Rating)
	}
	if book.Recommender != "" {
		content += fmt.Sprintf("  recommender = [\"%s\"]\n", book.Recommender)
	}
	content += "+++\n\n"

	return content
}

func main() {
	// Open CSV file
	file, err := os.Open("book_import.csv")
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	// Read CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	// Check if CSV has headers
	if len(records) < 2 {
		fmt.Println("CSV file doesn't have enough data.")
		return
	}

	// Extract headers
	headers := records[0]
	headerIndices := make(map[string]int)
	for i, header := range headers {
		headerIndices[strings.ToLower(header)] = i
	}

	// Get current date in ISO 8601 format
	currentDate := time.Now().Format("2006-01-02")

	// Process each row after the header
	for _, record := range records[1:] {
		book := Book{
			Title:       record[headerIndices["title"]],
			ReadState:   record[headerIndices["readstate"]],
			AuthorName:  record[headerIndices["authorname"]],
			PubYear:     record[headerIndices["pubyear"]],
			Recommender: record[headerIndices["recommender"]],
			Rating:      record[headerIndices["rating"]],
			Date:        currentDate,
		}

		// Generate Markdown content
		content := generateMarkdownContent(book)

		// Sanitize file name
		fileName := sanitizeFileName(book.Title) + ".md"

		// Write content to Markdown file
		err = os.WriteFile(fileName, []byte(content), 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			continue
		}
		fmt.Printf("Generated %s\n", fileName)
	}
}
