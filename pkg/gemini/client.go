// Package gemini menyediakan klien untuk berinteraksi dengan Google Gemini API.
package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// TODO: Model dan URL bisa dibuat lebih dinamis melalui config jika diperlukan.
	geminiAPIURL = "https://generativelanguage.googleapis.com/v1/models/gemini-2.5-flash:generateContent"
)

// Client adalah klien untuk Gemini API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient membuat instance baru dari Gemini Client.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// GenerateContent mengirimkan prompt ke Gemini API dan mengembalikan respons teks.
func (c *Client) GenerateContent(ctx context.Context, prompt string) (string, error) {
	// 1. Membuat body request sesuai dengan struktur JSON yang dibutuhkan.
	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 2. Membuat HTTP request.
	url := fmt.Sprintf("%s?key=%s", geminiAPIURL, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 3. Mengirim request.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to gemini api: %w", err)
	}
	defer resp.Body.Close()

	// 4. Membaca dan memeriksa respons.
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini api returned non-200 status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	// 5. Mengekstrak teks dari respons.
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no content found in gemini response")
}

// --- Structs for JSON Marshalling/Unmarshalling ---

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}
