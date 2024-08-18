package youtube

import (
	"fmt"
	"net/http"
	"os"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func GetVideosInfo(videos []string) (map[string]string, error) {
	formResponse := make(map[string]string, 0)

	key := os.Getenv("YOUTUBE_KEY")
	if key == "" {
		return formResponse, fmt.Errorf("Empty youtube api key")
	}
	if len(videos) == 0 {
		return formResponse, nil
	}

	client := &http.Client{
		Transport: &transport.APIKey{
			Key: key,
		},
	}
	youtubeService, err := youtube.New(client)
	if err != nil {
		return formResponse, fmt.Errorf("failed to create youtube instance - %w", err)
	}

	call := youtubeService.Videos.List([]string{"snippet"})
	response, err := call.Id(videos...).Do()
	if err != nil {

		return formResponse, fmt.Errorf("failed to get info about videos - %w", err)
	}

	for _, item := range response.Items {
		formResponse[item.Id] = item.Snippet.Thumbnails.Default.Url
	}

	return formResponse, nil

}
