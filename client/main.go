package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AlexEr256/thumbnail/internal/api"
	dto "github.com/AlexEr256/thumbnail/internal/domain"
	"github.com/AlexEr256/thumbnail/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const operation = "get"

func parseCommand(command string) (*dto.Command, error) {
	commandDto := &dto.Command{}
	parts := strings.Split(command, " ")

	if len(parts) <= 1 {
		return nil, fmt.Errorf("invalid command. Provide arguments")
	}
	if parts[0] != operation {
		return nil, fmt.Errorf("invalid operations - %q. Should be 'get'", parts[0])
	}

	videos := strings.Split(parts[1], ",")
	videoLinks := make([]*api.Video, 0)

	for _, video := range videos {
		videoLinks = append(videoLinks, &api.Video{Link: video})
	}

	switch len(parts) {
	case 2:
		commandDto.IsAsync = false
	case 3:
		if parts[2] == "--async" {
			commandDto.IsAsync = true
		} else if parts[2] != "--async" {
			return nil, fmt.Errorf("invalid flag - %q. Remove invalid flags or provide --async", parts[2])
		}
	default:
		return nil, fmt.Errorf("failed to parse command. Check input data")
	}

	commandDto.Videos = videoLinks

	return commandDto, nil

}

func SetupEnterCommandHandler(client api.ProxyClient) {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		command := sc.Text()
		dto, err := parseCommand(command)
		if err != nil {
			log.Println(err)
			continue
		}

		resp, err := client.Get(context.Background(), &api.GetRequest{Videos: dto.Videos})
		if err != nil {
			log.Println(err)
			continue
		}

		if !dto.IsAsync {
			for link, url := range resp.Info {
				err := utils.DownloadFile(url, link)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		} else {
			utils.DownloadAsyncFiles(resp.Info)
		}

	}
}
func main() {

	c, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Can't establish connection with server ", err)
	}
	client := api.NewProxyClient(c)

	SetupEnterCommandHandler(client)

}
