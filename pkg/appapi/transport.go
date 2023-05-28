package appapi

import "net/http"

type Transport struct {
	APIKey    string
	Transport http.RoundTripper
}

func NewTransport(
	key string,
) *Transport {
	return &Transport{
		APIKey:    key,
		Transport: http.DefaultTransport,
	}
}

func (tsp *Transport) transport() http.RoundTripper {
	return tsp.Transport
}

func (tsp *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Api-Key", tsp.APIKey)

	res, err := tsp.transport().RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return res, err
}
