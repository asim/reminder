package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func generateImage(message string) string {
	apiKey := os.Getenv("ATLAS_API_KEY")
	if apiKey == "" {
		return ""
	}

	prompt := fmt.Sprintf(
		"A serene landscape scene inspired by the following spiritual reflection. "+
			"No people, no faces, no human figures, no animals. Only natural scenery. "+
			"No text, no words, no letters, no writing. "+
			"Mountains, rivers, deserts, forests, oceans, skies, or gardens. "+
			"Peaceful, contemplative, warm golden light, soft atmosphere. Photorealistic.\n\n%s",
		message,
	)

	body, _ := json.Marshal(map[string]interface{}{
		"model":  "black-forest-labs/flux-schnell",
		"prompt": prompt,
		"width":  1024,
		"height": 576,
	})

	req, err := http.NewRequest("POST", "https://api.atlascloud.ai/api/v1/model/generateImage", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Image: request error: %v\n", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Image: http error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	fmt.Printf("Image: initial response (%d): %s\n", resp.StatusCode, string(respBody))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ""
	}

	// Parse response as generic map to handle any shape
	var raw map[string]interface{}
	if err := json.Unmarshal(respBody, &raw); err != nil {
		fmt.Printf("Image: parse error: %v\n", err)
		return ""
	}

	// The response may have a "data" envelope or be flat
	data := raw
	if d, ok := raw["data"].(map[string]interface{}); ok {
		data = d
	}

	// Check for image URL in outputs
	if url := extractImageURL(data); url != "" {
		return url
	}

	// Find a poll URL
	pollURL := ""
	if urls, ok := data["urls"].(map[string]interface{}); ok {
		if u, ok := urls["get"].(string); ok {
			pollURL = u
		}
	}
	if pollURL == "" {
		if id, ok := data["id"].(string); ok && id != "" {
			pollURL = "https://api.atlascloud.ai/api/v1/model/prediction/" + id
		}
	}

	if pollURL == "" {
		fmt.Println("Image: no image URL and no poll URL found")
		return ""
	}

	fmt.Printf("Image: polling %s\n", pollURL)
	return pollForImage(apiKey, pollURL)
}

func pollForImage(apiKey, pollURL string) string {
	client := &http.Client{Timeout: 30 * time.Second}

	for i := 0; i < 30; i++ {
		time.Sleep(2 * time.Second)

		req, _ := http.NewRequest("GET", pollURL, nil)
		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Image: poll %d error: %v\n", i+1, err)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Printf("Image: poll %d (%d): %s\n", i+1, resp.StatusCode, string(body))

		var raw map[string]interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			continue
		}

		data := raw
		if d, ok := raw["data"].(map[string]interface{}); ok {
			data = d
		}

		if url := extractImageURL(data); url != "" {
			fmt.Printf("Image: ready after poll %d: %s\n", i+1, url)
			return url
		}

		// Check for terminal failure
		if status, _ := data["status"].(string); status == "failed" || status == "canceled" {
			if errMsg, _ := data["error"].(string); errMsg != "" {
				fmt.Printf("Image: generation failed: %s\n", errMsg)
			} else {
				fmt.Printf("Image: generation %s\n", status)
			}
			return ""
		}
	}

	fmt.Println("Image: timed out after 60s of polling")
	return ""
}

// extractImageURL looks for an image URL in any common response field
func extractImageURL(data map[string]interface{}) string {
	// Check "outputs" array (strings or objects with url field)
	if outputs, ok := data["outputs"]; ok && outputs != nil {
		if arr, ok := outputs.([]interface{}); ok {
			for _, item := range arr {
				if s, ok := item.(string); ok && isImageURL(s) {
					return s
				}
				if obj, ok := item.(map[string]interface{}); ok {
					for _, key := range []string{"url", "image", "src"} {
						if u, ok := obj[key].(string); ok && isImageURL(u) {
							return u
						}
					}
				}
			}
		}
	}

	// Check "output" singular
	if output, ok := data["output"].(string); ok && isImageURL(output) {
		return output
	}
	if output, ok := data["output"].(map[string]interface{}); ok {
		for _, key := range []string{"url", "image", "src"} {
			if u, ok := output[key].(string); ok && isImageURL(u) {
				return u
			}
		}
	}

	// Check "image" or "image_url" top-level
	for _, key := range []string{"image", "image_url", "url"} {
		if u, ok := data[key].(string); ok && isImageURL(u) {
			return u
		}
	}

	return ""
}

func isImageURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
