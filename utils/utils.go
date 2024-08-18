package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func DownloadFile(URL string, link string) error {
	if URL == "" {
		return fmt.Errorf("Empty url provided. Check link of the video")
	}
	response, err := http.Get(URL)
	fmt.Println(URL, response.StatusCode)
	if err != nil {
		return fmt.Errorf("failed to get thumbnail - %w", err)
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code for getting thumbnail - %d", response.StatusCode)
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body - %w", err)
	}
	fileName := link + ".jpg"

	return os.WriteFile(fileName, data, 0644)
}

func DownloadAsyncFiles(info map[string]string) {
	var wg sync.WaitGroup
	resultCh := make(chan string, len(info))
	errorCh := make(chan error, len(info))

	for link, url := range info {
		wg.Add(1)

		go func(url string, link string) {
			defer wg.Done()
			err := DownloadFile(url, link)

			if err != nil {
				errorCh <- err
			} else {
				resultCh <- fmt.Sprintf("Downloaded: %s", link)
			}
		}(url, link)
	}

	go func() {
		wg.Wait()
		close(resultCh)
		close(errorCh)
	}()

	for result := range resultCh {
		fmt.Println(result)
	}

	for err := range errorCh {
		fmt.Println("Error:", err)
	}
}
