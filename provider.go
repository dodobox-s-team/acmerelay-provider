package acmerelay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/libdns/libdns"
)

type Provider struct {
	APIKey      string `json:"api_key,omitempty"`
	APIEndpoint string `json:"api_endpoint,omitempty"`
	mu          sync.Mutex
}

func (p *Provider) getClient() (*Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return &Client{p.APIKey}, nil
}

func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	client, err := p.getClient()
	if err != nil {
		return nil, err
	}

	for _, r := range records {
		record := Record{
			Subdomain: r.Name,
			Target:    r.Value,
			TTL:       int(r.TTL),
		}

		data, err := json.Marshal(record)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, "POST", p.APIEndpoint, bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}

		_, err = client.doRequest(req)
		if err != nil {
			return nil, err
		}
	}
	return records, err
}

func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	client, err := p.getClient()
	if err != nil {
		return nil, err
	}

	for _, r := range records {
		path := fmt.Sprintf("%s?subdomain=%s", p.APIEndpoint, r.Name)
		req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
		if err != nil {
			return nil, err
		}

		_, err = client.doRequest(req)
		if err != nil {
			return nil, err
		}
	}
	var result []libdns.Record
	return result, err
}
