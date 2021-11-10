package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetLogs(ctx context.Context) (*Order, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	log := Order{}
	if err := json.Unmarshal(body, &log); err != nil {
		return nil, err
	}

	return &log, nil
}

func (c *Client) CreateLog(ctx context.Context, logReq *Order) (*Order, error) {
	rb, err := json.Marshal(logReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	log := Order{}
	if err := json.Unmarshal(body, &log); err != nil {
		return nil, err
	}

	return &log, nil
}

func (c *Client) UpdateLog(ctx context.Context, orderID string, logReq *Order) (*Order, error) {
	rb, err := json.Marshal(logReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/%s", c.HostURL, orderID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	log := Order{}
	if err := json.Unmarshal(body, &log); err != nil {
		return nil, err
	}

	return &log, nil
}

func (c *Client) DeleteLog(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/", c.HostURL), nil)
	if err != nil {
		return err
	}

	if _, err := c.doRequest(req); err != nil {
		return err
	}

	return nil
}
