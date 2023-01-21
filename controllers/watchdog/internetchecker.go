package watchdog

import (
	"fmt"
	"net/http"
	"time"
)

type InternetCheckerImpl struct {
	url string
}

func NewInternetChecker(url string) *InternetCheckerImpl {
	return &InternetCheckerImpl{
		url: url,
	}
}

func (ic *InternetCheckerImpl) IsInternetAvailable() (bool, error) {
	fmt.Printf("About to send request to check if Internet is available: %s \n", ic.url)

	req, err := http.NewRequest(http.MethodGet, ic.url, nil)
	if err != nil {
		fmt.Printf("Error while creating request: %s \n", err.Error())
		return false, err
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error while sending request: %s \n", err.Error())
		return false, nil
	}

	if resp.StatusCode >= 400 {
		err = fmt.Errorf("response status is >= 400: %d", resp.StatusCode)
		fmt.Printf("Error: %s \n", err.Error())
		return false, nil
	}

	return true, nil
}
