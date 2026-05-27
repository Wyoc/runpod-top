package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const podsQuery = `{
  myself {
    pods {
      id
      name
      desiredStatus
      costPerHr
      gpuCount
      memoryInGb
      vcpuCount
      uptimeSeconds
      machine {
        gpuDisplayName
        location
      }
      runtime {
        uptimeInSeconds
        gpus {
          id
          gpuUtilPercent
          memoryUtilPercent
        }
        container {
          cpuPercent
          memoryPercent
        }
        ports {
          ip
          isIpPublic
          privatePort
          publicPort
          type
        }
      }
    }
  }
}`

type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.runpod.io/graphql",
	}
}

func (c *Client) FetchPods(ctx context.Context) ([]Pod, error) {
	body, err := json.Marshal(map[string]string{"query": podsQuery})
	if err != nil {
		return nil, fmt.Errorf("marshal query: %w", err)
	}

	url := fmt.Sprintf("%s?api_key=%s", c.baseURL, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("api request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	var result graphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(result.Errors) > 0 {
		msgs := make([]string, len(result.Errors))
		for i, e := range result.Errors {
			msgs[i] = e.Message
		}
		return nil, fmt.Errorf("graphql: %s", strings.Join(msgs, "; "))
	}

	return result.Data.Myself.Pods, nil
}
