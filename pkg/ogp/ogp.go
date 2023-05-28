package ogp

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/dyatlov/go-opengraph/opengraph"
)

type OGP struct {
	URL         string
	Title       string
	Description string
	Thumbnail   string
}

func Create(
	ctx context.Context,
	url string,
) (OGP, error) {
	client := http.DefaultClient

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return OGP{}, err
	}

	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return OGP{}, err
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	og := opengraph.NewOpenGraph()

	if err = og.ProcessHTML(strings.NewReader(string(body))); err != nil {
		return OGP{}, err
	}

	thumbnail := ""
	if len(og.Images) > 0 {
		thumbnail = og.Images[0].URL
	}

	return OGP{
		URL:         url,
		Title:       og.Title,
		Description: og.Description,
		Thumbnail:   thumbnail,
	}, nil
}
