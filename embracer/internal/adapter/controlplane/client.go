package controlplane

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

type registerRequest struct {
	Cluster string `json:"cluster"`
	Addr    string `json:"addr"`
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	return &Client{httpClient: httpClient, baseURL: baseURL}
}

func (c *Client) Register(cluster, addr string) error {
	reqBody := registerRequest{Cluster: cluster, Addr: addr}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}
	
	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/register", c.baseURL),
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}
