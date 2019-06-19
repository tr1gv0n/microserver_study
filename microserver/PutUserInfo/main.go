package main

import (
	"github.com/micro/go-grpc"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"sss/PutUserInfo/handler"
	"sss/PutUserInfo/subscriber"

	example "sss/PutUserInfo/proto/example"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.PutUserInfo"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example))

    // Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
