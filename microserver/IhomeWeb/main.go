package main

import (
        "github.com/julienschmidt/httprouter"
        "sss/IhomeWeb/handler"
        _ "sss/IhomeWeb/models"
        "github.com/micro/go-log"
        "github.com/micro/go-web"
        "net/http"
)

func main() {
	// create new web service
        service := web.NewService(
                web.Name("go.micro.web.IhomeWeb"),
                web.Version("latest"),
                web.Address(":9999"),
        )

	// initialise service
        if err := service.Init(); err != nil {
                log.Fatal(err)
        }

	// register html handler
	//service.Handle("/", http.FileServer(http.Dir("html")))

	// register call handler
	//service.HandleFunc("/example/call", handler.ExampleCall)

	rou := httprouter.New()
	//文件服务器 映射  静态页面
	rou.NotFound = http.FileServer(http.Dir("html"))

	rou.GET("/api/v1.0/areas",handler.GetArea)

	service.Handle("/",rou)
	// run service
        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}
