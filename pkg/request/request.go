package request

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

/**
* Generic HTTP request handler
*/
func Request(url string) (string, error) {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return url, errors.New("Failed to fetch " + url)
	}

	// Format body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", errors.New("Couldn't read body")
	}

	response := string(body)

	return response, nil
}

/**
* Moralis API request helper
*/
func APIRequest(url string) (string, error) {

	requestUrl := os.Getenv("MORALIS_API_URL") + url

	// Create GET request
	req, _ := http.NewRequest("GET", requestUrl, nil)


	req.Header = http.Header{
		"x-api-key": []string{os.Getenv("MORALIS_API_KEY")},
	}

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return "", errors.New("Request failed")
	}

	// Format body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
        return "", errors.New("Couldn't read body")
    }

	response := string(body)

	// Send any response errors
	if (resp.StatusCode != 200) && (resp.StatusCode != 202) {
		errorMessage := "ERROR: Moralis API request failed for " + requestUrl

		log.Println(errorMessage)
		return response, errors.New(errorMessage)
	}

	return response, nil
}