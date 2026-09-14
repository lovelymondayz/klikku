package storage

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
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
			Timeout: 60 * time.Second,
		},
	}
}

func (i *Immich) Upload(bucket, objectName string, data []byte, contentType string) error {
	url := fmt.Sprintf("%s/api/assets", i.baseURL)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("assetData", filepath.Base(objectName))
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}
	if _, err := part.Write(data); err != nil {
		return fmt.Errorf("write data: %w", err)
	}

	writer.WriteField("deviceAssetId", objectName)
	writer.WriteField("deviceId", "klikku")
	writer.WriteField("fileCreatedAt", time.Now().UTC().Format(time.RFC3339))
	writer.WriteField("fileModifiedAt", time.Now().UTC().Format(time.RFC3339))
	writer.WriteField("isFavorite", "false")

	writer.Close()

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-api-key", i.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := i.client.Do(req)
	if err != nil {
		return fmt.Errorf("upload to immich: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
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
		return nil, fmt.Errorf("asset not found in immich (status %d)", resp.StatusCode)
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
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("immich delete failed (status %d): %s", resp.StatusCode, string(body))
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
	url := fmt.Sprintf("%s/api/albums", i.baseURL)

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

	return nil
}
