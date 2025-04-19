package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Status struct {
	UpTime                float64 `json:"upTime"`
	FreeHeap              float64 `json:"freeHeap"`
	MinHeap               float64 `json:"minHeap"`
	OpeningsCount         float64 `json:"openingsCount"`
	OpenDuration          float64 `json:"openDuration"`
	WiFiRSSI              string  `json:"wifiRSSI"`
	GarageDoorState       string  `json:"garageDoorState"`
	GarageDoorLockedState string  `json:"garageLockState"`
	DeviceName            string  `json:"deviceName"`
	GarageLight           bool    `json:"garageLightOn"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Authenticate() error {
	req, err := http.NewRequest("GET", c.baseURL+"/auth", nil)
	if err != nil {
		return fmt.Errorf("failed to build auth request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Successfully called ratgdo Auth API")

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) GetStatus() (*Status, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/status.json")
	if err != nil {
		return nil, fmt.Errorf("error fetching status.json: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status.json returned %d", resp.StatusCode)
	}

	log.Printf("Successfully called ratgdo Status API")

	var status Status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("error decoding status.json: %w", err)
	}

	return &status, nil
}
