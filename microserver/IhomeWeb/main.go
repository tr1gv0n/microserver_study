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
        //获取图片验证码服务/api/v1.0/imagecode/61bb9476-2a28-4180-9ce6-56a99e16003a
    rou.GET("/api/v1.0/imagecode/:uuid", handler.GetImageCd)

        // 短信验证码服务
    rou.GET("/api/v1.0/smscode/:mobile", handler.GetSmscd)
        // 注册
    rou.POST("/api/v1.0/users", handler.PostRet)

        //欺骗浏览器  session index
    rou.GET("/api/v1.0/session", handler.GetSession)
        //session
    rou.GET("/api/v1.0/house/index", handler.GetIndex)

    rou.POST("/api/v1.0/sessions", handler.PostLogin)

	rou.DELETE("/api/v1.0/session",handler.DeleteSession)

    service.Handle("/",rou)
	// run service
        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}
