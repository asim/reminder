package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type imageRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type imageResponse struct {
	ID      string   `json:"id"`
	Outputs []string `json:"outputs"`
	Status  string   `json:"status"`
}

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

	body, _ := json.Marshal(imageRequest{
		Model:  "black-forest-labs/flux-schnell",
		Prompt: prompt,
		Width:  1024,
		Height: 576,
	})

	req, err := http.NewRequest("POST", "https://api.atlascloud.ai/api/v1/model/generateImage", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Image generation request error: %v\n", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Image generation error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("Image generation returned status %d: %s\n", resp.StatusCode, string(respBody))
		return ""
	}

	fmt.Printf("Image generation response (%d): %s\n", resp.StatusCode, string(respBody))

	var result imageResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Printf("Image generation parse error: %v\n", err)
		return ""
	}

	if len(result.Outputs) > 0 && result.Outputs[0] != "" {
		fmt.Printf("Image generated: %s\n", result.Outputs[0])
		return result.Outputs[0]
	}

	// If async, poll for result
	if result.ID != "" {
		return pollImageResult(apiKey, result.ID)
	}

	fmt.Println("Image generation: no outputs and no async ID in response")
	return ""
}

func pollImageResult(apiKey, id string) string {
	client := &http.Client{Timeout: 30 * time.Second}
	url := fmt.Sprintf("https://api.atlascloud.ai/api/v1/model/generateImage/%s", id)

	for i := 0; i < 20; i++ {
		time.Sleep(3 * time.Second)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Printf("Image poll %d (%d): %s\n", i+1, resp.StatusCode, string(body))

		var result imageResponse
		if err := json.Unmarshal(body, &result); err != nil {
			continue
		}

		if len(result.Outputs) > 0 && result.Outputs[0] != "" {
			fmt.Printf("Image ready (poll %d): %s\n", i+1, result.Outputs[0])
			return result.Outputs[0]
		}
	}

	fmt.Println("Image generation timed out after polling")
	return ""
}
