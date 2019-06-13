package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-grpc"
	"net/http"
	GETAREA "sss/GetArea/proto/example"
	"sss/IhomeWeb/models"

	//example "github.com/micro/examples/template/srv/proto/example"
)

func ExampleCall(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	//// decode the incoming request as json
	//var request map[string]interface{}
	//if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	//	http.Error(w, err.Error(), 500)
	//	return
	//}
	//
	//cli := grpc.NewService()
	//cli.Init()
	//// call the backend service
	//exampleClient := example.NewExampleService("go.micro.srv.template", cli.Client())
	//rsp, err := exampleClient.Call(context.TODO(), &example.Request{
	//	Name: request["name"].(string),
	//})
	//if err != nil {
	//	http.Error(w, err.Error(), 500)
	//	return
	//}
	//
	//// we want to augment the response
	//response := map[string]interface{}{
	//	"msg": rsp.Msg,
	//	"ref": time.Now().UnixNano(),
	//}
	//
	//// encode and write the response as json
	//if err := json.NewEncoder(w).Encode(response); err != nil {
	//	http.Error(w, err.Error(), 500)
	//	return
	//}
}

func GetArea(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	fmt.Println("获取地域信息服务 GetArea  /api/v1.0/areas")

	cli := grpc.NewService()
	cli.Init()
	// call the backend service
	exampleClient := GETAREA.NewExampleService("go.micro.srv.GetArea", cli.Client())
	rsp, err := exampleClient.GetArea(context.TODO(), &GETAREA.Request{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	area_list := []models.Area{}

	for _,value := range rsp.Data{
		tmp:= models.Area{Id:int(value.Aid),Name:value.Aname}
		area_list = append(area_list,tmp)
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":rsp.Errno,
		"errmsg":rsp.Errmsg,
		"data":area_list,
	}

	w.Header().Set("Content-Type","application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
