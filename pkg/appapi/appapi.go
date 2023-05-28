package appapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type V1ArticleShareRequest struct {
	// Description description
	Description *string `json:"description,omitempty"`

	// Thumbnail サムネイルのURL
	Thumbnail *string `json:"thumbnail,omitempty"`

	// Title タイトル
	Title *string `json:"title,omitempty"`

	// Url 記事のURL
	Url string `json:"url"`
}

func V1InternalArticleShare(
	ctx context.Context,
	endpoint string,
	key string,
	body V1ArticleShareRequest,
) (*http.Response, error) {
	client := &http.Client{Transport: NewTransport(key)}

	var bodyReader io.Reader

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed marshal body: %w", err)
	}

	bodyReader = bytes.NewReader(buf)

	url := fmt.Sprintf("%s/v1/internal/articles", endpoint)

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed new request: %w", err)
	}

	req = req.WithContext(ctx)

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed do request: %w", err)
	}

	return res, nil
}
