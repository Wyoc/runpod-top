package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) StartPod(ctx context.Context, podID string) error {
	return c.mutate(ctx, "podResume", podID)
}

func (c *Client) StopPod(ctx context.Context, podID string) error {
	return c.mutate(ctx, "podStop", podID)
}

func (c *Client) mutate(ctx context.Context, mutation, podID string) error {
	query := fmt.Sprintf(`mutation { %s(input: {podId: "%s"}) { id desiredStatus } }`, mutation, podID)

	body, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		return fmt.Errorf("marshal mutation: %w", err)
	}

	url := fmt.Sprintf("%s?api_key=%s", c.baseURL, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("api request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	var result graphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if len(result.Errors) > 0 {
		return fmt.Errorf("graphql: %s", result.Errors[0].Message)
	}

	return nil
}
