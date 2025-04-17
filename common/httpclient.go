package common

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetContent(URL string) ([]byte, error) {

	client := http.Client{Timeout: time.Second * 5}
	resp, err := client.Get(URL)

	if err != nil {
		return nil, fmt.Errorf("HTTP_CLIENT: Can't connect to API")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"HTTP_CLIENT: Unexpected API response: %d (%s)",
			resp.StatusCode,
			http.StatusText(resp.StatusCode),
		)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("HTTP_CLIENT: Can't read data from API")
	}

	return body, nil
}
