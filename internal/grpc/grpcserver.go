package proxyserver

import (
	"context"

	"github.com/AlexEr256/thumbnail/internal/api"
	dto "github.com/AlexEr256/thumbnail/internal/domain"
	sqlite "github.com/AlexEr256/thumbnail/internal/storage"
	"github.com/AlexEr256/thumbnail/internal/youtube"
)

type GRPCServer struct {
	api.UnimplementedProxyServer
	Storage *sqlite.Storage
}

func (s GRPCServer) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	videos := req.GetVideos()

	//init response
	response := make(map[string]string, 0)
	for _, video := range videos {
		response[video.Link] = ""
	}

	//check if videos already exist in sqlite3 and update response with existed urls
	videosList := make([]string, 0)
	for _, video := range videos {
		videosList = append(videosList, video.Link)
	}
	cachedVideos, err := s.Storage.ListVideos(videosList)
	if err != nil {
		return &api.GetResponse{Info: response, Error: err.Error()}, err
	}
	for _, video := range cachedVideos {
		link, url := video.Link, video.Url
		_, ok := response[link]
		if ok {
			response[link] = url
		}
	}

	// form reqs to get info about videos from youtube
	videosReqs := make([]string, 0)
	for key, value := range response {
		if value == "" {
			videosReqs = append(videosReqs, key)
		}
	}

	youtubeVideosInfo, err := youtube.GetVideosInfo(videosReqs)
	if err != nil {
		return &api.GetResponse{Info: response, Error: err.Error()}, err
	}

	//update response with new information from youtube
	for key, link := range youtubeVideosInfo {
		_, ok := response[key]
		if ok {
			response[key] = link
		}
	}
	//insert new videos into cache
	videosToInsertInCache := make([]dto.VideoInfo, 0)
	for link, url := range youtubeVideosInfo {
		videosToInsertInCache = append(videosToInsertInCache, dto.VideoInfo{Link: link, Url: url})
	}

	if len(videosToInsertInCache) != 0 {
		err = s.Storage.SaveVideos(videosToInsertInCache)
		if err != nil {
			return &api.GetResponse{Info: response, Error: err.Error()}, err
		}
	}

	return &api.GetResponse{Info: response, Error: ""}, nil

}
