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

func generateImage(message string) {
	apiKey := os.Getenv("ATLAS_API_KEY")
	if apiKey == "" {
		return
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

	if url := extractImageURL(data); url != "" {
		setImage(url)
		return
	}

	// Find poll URL and poll in background
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

	if pollURL != "" {
		go pollForImage(apiKey, pollURL)
	}
}

func pollForImage(apiKey, pollURL string) {
	client := &http.Client{Timeout: 10 * time.Second}

	for i := 0; i < 10; i++ {
		time.Sleep(3 * time.Second)

		req, _ := http.NewRequest("GET", pollURL, nil)
		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err := client.Do(req)
		if err != nil {
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
			setImage(url)
			return
		}

		if status, _ := data["status"].(string); status == "failed" || status == "canceled" {
			fmt.Printf("Image: generation %s\n", status)
			return
		}
	}

	fmt.Println("Image: gave up after 30s")
}

func setImage(url string) {
	fmt.Printf("Image: ready: %s\n", url)
	mtx.Lock()
	dailyImage = url
	mtx.Unlock()
}

func extractImageURL(data map[string]interface{}) string {
	if outputs, ok := data["outputs"]; ok && outputs != nil {
		if arr, ok := outputs.([]interface{}); ok {
			for _, item := range arr {
				if s, ok := item.(string); ok && isURL(s) {
					return s
				}
				if obj, ok := item.(map[string]interface{}); ok {
					for _, key := range []string{"url", "image", "src"} {
						if u, ok := obj[key].(string); ok && isURL(u) {
							return u
						}
					}
				}
			}
		}
	}

	if output, ok := data["output"].(string); ok && isURL(output) {
		return output
	}

	for _, key := range []string{"image", "image_url", "url"} {
		if u, ok := data[key].(string); ok && isURL(u) {
			return u
		}
	}

	return ""
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
