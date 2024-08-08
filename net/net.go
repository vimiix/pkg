// Copyright (c) 2024 Qian Yao
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package net

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DownloadParams struct {
	URL      string
	Output   io.Writer
	Timeout  time.Duration
	Headers  map[string]string
	AuthUser string
	AuthPass string
}

// Download downloads a resource from the specified URL
// and saves it to the provided output writer.
func Download(ctx context.Context, params *DownloadParams) error {
	if params == nil {
		return fmt.Errorf("params cannot be nil")
	}

	if params.Timeout < 0 {
		params.Timeout = 0
	}

	cli := http.Client{Timeout: params.Timeout}

	req, err := http.NewRequestWithContext(ctx, "GET", params.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range params.Headers {
		req.Header.Set(k, v)
	}

	if params.AuthUser != "" && params.AuthPass != "" {
		req.SetBasicAuth(params.AuthUser, params.AuthPass)
	}

	resp, err := cli.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(params.Output, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy response body: %w", err)
	}

	return nil
}
