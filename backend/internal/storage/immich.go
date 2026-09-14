package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Immich struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewImmich(apiKey, baseURL string) *Immich {
	return &Immich{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (i *Immich) Upload(bucket, objectName string, data []byte, contentType string) error {
	url := fmt.Sprintf("%s/api/assets", i.baseURL)

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-api-key", i.apiKey)
	req.Header.Set("Content-Type", contentType)

	resp, err := i.client.Do(req)
	if err != nil {
		return fmt.Errorf("upload to immich: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("immich upload failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (i *Immich) Download(bucket, objectName string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/assets/%s/original", i.baseURL, objectName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-api-key", i.apiKey)

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download from immich: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("asset not found in immich")
	}

	return io.ReadAll(resp.Body)
}

func (i *Immich) Delete(bucket, objectName string) error {
	url := fmt.Sprintf("%s/api/assets/%s", i.baseURL, objectName)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-api-key", i.apiKey)

	resp, err := i.client.Do(req)
	if err != nil {
		return fmt.Errorf("delete from immich: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("immich delete failed (status %d)", resp.StatusCode)
	}

	return nil
}

func (i *Immich) Exists(bucket, objectName string) bool {
	url := fmt.Sprintf("%s/api/assets/%s", i.baseURL, objectName)

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false
	}

	req.Header.Set("x-api-key", i.apiKey)

	resp, err := i.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func (i *Immich) Serve(bucket, objectName string, w io.Writer) error {
	data, err := i.Download(bucket, objectName)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func (i *Immich) GetURL(bucket, objectName string) string {
	return fmt.Sprintf("%s/api/assets/%s/thumbnail", i.baseURL, objectName)
}

func (i *Immich) Ping() error {
	url := fmt.Sprintf("%s/api/server-info", i.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-api-key", i.apiKey)

	resp, err := i.client.Do(req)
	if err != nil {
		return fmt.Errorf("ping immich: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("immich ping failed (status %d)", resp.StatusCode)
	}

	var info struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return fmt.Errorf("decode immich response: %w", err)
	}

	fmt.Printf("✅ Immich connected (v%s)\n", info.Version)
	return nil
}
