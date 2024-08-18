package dto

import "github.com/AlexEr256/thumbnail/internal/api"

type Command struct {
	Videos  []*api.Video
	IsAsync bool
}

type VideoInfo struct {
	Link string
	Url  string
}
