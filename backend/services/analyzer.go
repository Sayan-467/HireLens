package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	pdf "github.com/ledongthuc/pdf"
)

// ExtractTextFromPdfFile extracts clean text from a local PDF file
func ExtractTextFromPdfFile(filePath string) (string, error) {
	// Open the PDF file
	file, reader, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %v", err)
	}
	defer file.Close()

	var textBuilder strings.Builder
	totalPages := reader.NumPage()

	fmt.Printf("📄 PDF has %d pages\n", totalPages)

	// Extract text from each page
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page := reader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			fmt.Printf("⚠️  Warning: Could not extract text from page %d: %v\n", pageNum, err)
			continue
		}

		textBuilder.WriteString(text)
		textBuilder.WriteString("\n")
	}

	extractedText := textBuilder.String()

	// Validate the extracted text
	if len(extractedText) == 0 {
		return "", fmt.Errorf("no text could be extracted from PDF")
	}

	fmt.Printf("Extracted text length: %d characters\n", len(extractedText))

	// Limit text size to 50KB
	maxChars := 50000
	if len(extractedText) > maxChars {
		fmt.Printf("Text too long (%d chars), truncating to %d chars\n", len(extractedText), maxChars)
		extractedText = extractedText[:maxChars]
	}
	return extractedText, nil
}

func AnalyzeResumeText(text string, jobDescription string) (string, error) {
	// Call the local FastAPI analyzer service
	analyzerURL := os.Getenv("ANALYZER_URL")
	if analyzerURL == "" {
		analyzerURL = "http://localhost:8000/analyze" // default
	}

	fmt.Println("Calling analyzer at:", analyzerURL)
	fmt.Println("Text length:", len(text))
	if jobDescription != "" {
		fmt.Println("Job description provided, length:", len(jobDescription))
	}

	requestBody := map[string]string{
		"text":            text,
		"job_description": jobDescription,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", analyzerURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Analyzer connection error:", err)
		return "", fmt.Errorf("failed to call analyzer service: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Analyzer response status:", resp.StatusCode)

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Analyzer response body:", string(respData))
		return "", fmt.Errorf("analyzer service returned status: %d", resp.StatusCode)
	}

	fmt.Println("Analyzer response:", string(respData))
	return string(respData), nil
}
