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

var pendingImagePollURL string
var pendingImageAPIKey string

func generateImage(message string) string {
	apiKey := os.Getenv("ATLAS_API_KEY")
	if apiKey == "" {
		return ""
	}

	// Check if there's a pending image from a previous request
	if pendingImagePollURL != "" {
		img := checkImageResult(pendingImageAPIKey, pendingImagePollURL)
		pendingImagePollURL = ""
		if img != "" {
			// Also kick off the next generation
			submitImageRequest(apiKey, message)
			return img
		}
	}

	// Submit a new request (non-blocking)
	submitImageRequest(apiKey, message)
	return ""
}

func submitImageRequest(apiKey, message string) {
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
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Image: http error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	fmt.Printf("Image: submitted (%d): %s\n", resp.StatusCode, string(respBody))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(respBody, &raw); err != nil {
		return
	}

	data := raw
	if d, ok := raw["data"].(map[string]interface{}); ok {
		data = d
	}

	// If image is ready immediately, save it
	if url := extractImageURL(data); url != "" {
		fmt.Printf("Image: ready immediately: %s\n", url)
		pendingImagePollURL = ""
		mtx.Lock()
		dailyImage = url
		mtx.Unlock()
		return
	}

	// Store the poll URL for next cycle
	if urls, ok := data["urls"].(map[string]interface{}); ok {
		if u, ok := urls["get"].(string); ok {
			pendingImagePollURL = u
			pendingImageAPIKey = apiKey
			fmt.Printf("Image: will check result next cycle: %s\n", u)
			return
		}
	}
	if id, ok := data["id"].(string); ok && id != "" {
		pendingImagePollURL = "https://api.atlascloud.ai/api/v1/model/prediction/" + id
		pendingImageAPIKey = apiKey
		fmt.Printf("Image: will check result next cycle: %s\n", pendingImagePollURL)
	}
}

func checkImageResult(apiKey, pollURL string) string {
	req, _ := http.NewRequest("GET", pollURL, nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Image: check error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	fmt.Printf("Image: check result (%d): %s\n", resp.StatusCode, string(body))

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return ""
	}

	data := raw
	if d, ok := raw["data"].(map[string]interface{}); ok {
		data = d
	}

	return extractImageURL(data)
}

func extractImageURL(data map[string]interface{}) string {
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
