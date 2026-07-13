package myhttp

import (
	"fmt"
	"net/http"
)

func Send(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		return nil, err
	}
	return resp, nil
}
