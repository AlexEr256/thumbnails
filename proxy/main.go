package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/AlexEr256/thumbnail/internal/api"
	proxyserver "github.com/AlexEr256/thumbnail/internal/grpc"
	sqlite "github.com/AlexEr256/thumbnail/internal/storage"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to parse dotenv file ", err)
	}

	path := os.Getenv("SQLITE_PATH")
	if path == "" {
		log.Fatal("Empty sqlite3 path")
	}

	s, err := sqlite.New(path)
	if err != nil {
		log.Fatal("Failed to connect to sqlite3 ", err)
	}
	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("Can't init storage", err)
	}

	server := grpc.NewServer()

	srv := proxyserver.GRPCServer{Storage: s}

	api.RegisterProxyServer(server, srv)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Can't access port ", err)
	}

	if err := server.Serve(l); err != nil {
		log.Fatal("Can't establish server connection ", err)
	}
}
