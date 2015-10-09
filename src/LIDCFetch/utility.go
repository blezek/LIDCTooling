package main

import (
	"net/http"
	"net/url"
)

// Prepare the request by adding the api_key header
func prepareRequest(endpoint string, params url.Values) (*http.Request, error) {
	req, err := http.NewRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		logger.Error("Error fetching collections")
		return nil, err
	}
	req.Header.Add("api_key", APIKey)
	req.URL.RawQuery = params.Encode()
	return req, err
}
