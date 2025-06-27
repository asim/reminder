package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort" // Import the sort package
	"strconv"
	"strings"
)

// Verse represents a single verse entry in the JSON output.
type Verse struct {
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
	Text    string `json:"text"`
}

// ArabicData represents the overall structure of the JSON output.
// The keys are chapter numbers (as strings), and values are slices of Verse structs.
type ArabicData map[string][]Verse

func main() {
	inputCSVPath := "quran/data/uthmani.txt"
	outputJSONPath := "quran/data/arabic.json"

	err := parseCSVToJSON(inputCSVPath, outputJSONPath)
	if err != nil {
		fmt.Printf("Error processing CSV: %v\n", err)
		return
	}

	fmt.Printf("CSV parsing complete. Output saved to %s\n", outputJSONPath)
}

// parseCSVToJSON reads the CSV, processes lines, and writes to JSON with ordered map keys.
func parseCSVToJSON(csvFilePath, jsonFilePath string) error {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	outputData := make(ArabicData)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Stop processing on a blank line
		if line == "" {
			fmt.Println("Encountered blank line, stopping processing.")
			break
		}

		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			// Skip lines that don't have exactly 3 parts
			continue
		}

		chapterStr := parts[0]
		verseStr := parts[1]
		text := parts[2]

		chapter, err := strconv.Atoi(chapterStr)
		if err != nil {
			// Skip lines where chapter is not a valid number
			continue
		}

		verse, err := strconv.Atoi(verseStr)
		if err != nil {
			// Skip lines where verse is not a valid number
			continue
		}

		// Only process lines where the chapter is between 1 and 114
		if chapter >= 1 && chapter <= 114 {
			// Convert chapter back to string for the map key
			chapterKey := strconv.Itoa(chapter)
			outputData[chapterKey] = append(outputData[chapterKey], Verse{
				Chapter: chapter,
				Verse:   verse,
				Text:    text,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading CSV file: %w", err)
	}

	// --- START: Logic to ensure numeric order of map keys in JSON ---

	// 1. Extract all unique chapter numbers (as integers) from the map keys
	var chapterNumbers []int
	for k := range outputData {
		num, err := strconv.Atoi(k)
		if err == nil { // Should not happen if conversion to int was successful previously
			chapterNumbers = append(chapterNumbers, num)
		}
	}

	// 2. Sort the chapter numbers numerically
	sort.Ints(chapterNumbers)

	// 3. Manually construct the JSON output string to preserve order
	var orderedJSONParts []string
	for _, chapterNum := range chapterNumbers {
		chapterKey := strconv.Itoa(chapterNum)
		verses := outputData[chapterKey]

		// Marshal the slice of verses for the current chapter with indentation
		// The inner content (verses array) will also be nicely formatted
		verseJSON, err := json.MarshalIndent(verses, "  ", "  ") // Indent the value (array) by 2 spaces relative to its key
		if err != nil {
			return fmt.Errorf("failed to marshal verses for chapter %s: %w", chapterKey, err)
		}

		// Format the "key": value pair for this chapter
		// Ensure the key itself is quoted correctly in JSON
		orderedJSONParts = append(orderedJSONParts, fmt.Sprintf("  \"%s\": %s", chapterKey, string(verseJSON)))
	}

	// Combine all parts into the final JSON string, adding the outer curly braces and newlines
	finalJSONString := "{\n" + strings.Join(orderedJSONParts, ",\n") + "\n}"
	jsonData := []byte(finalJSONString)

	// --- END: Logic to ensure numeric order of map keys in JSON ---

	// Write the JSON data to the output file
	err = os.WriteFile(jsonFilePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

