package main

import (
	"net/http"
	"net/url"
)

// Prepare the request by adding the api_key header
func prepareRequest(endpoint string, inParams url.Values) (*http.Request, error) {

	params := url.Values{}
	for key, _ := range inParams {
		params.Set(key, inParams.Get(key))
	}

	// add the APIKey to the query parameters
	params.Set("api_key", APIKey)

	req, err := http.NewRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		logger.Error("Error fetching collections")
		return nil, err
	}
	req.Header.Add("api_key", APIKey)
	req.URL.RawQuery = params.Encode()
	return req, err
}
