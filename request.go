package acmerelay

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	apiKey string
}

type Record struct {
	Subdomain string `json:"subdomain"`
	Target    string `json:"target"`
	TTL       int    `json:"ttl"`
}

type apiResponse struct {
	Success bool                   `json:"success"`
	Detail  string                 `json:"detail,omitempty"`
	Result  map[string]interface{} `json:"result,omitempty"`
}

func (client *Client) doRequest(req *http.Request) (*apiResponse, error) {
	// Set the correct headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.apiKey))
	req.Header.Set("Content-Type", "application/json")
	// Send the request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	// Read the data sent by the server
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the data
	var apiResp apiResponse
	if json.Unmarshal(body, &apiResp) != nil {
		return nil, fmt.Errorf("unable to parse response: %s", string(body))
	}

	// Check if the api returned an error
	if len(apiResp.Detail) > 0 {
		return nil, fmt.Errorf("acmerelay api returned error in response. Err: %s", apiResp.Detail)
	}

	// Otherwise return the api's response
	return &apiResp, nil
}
