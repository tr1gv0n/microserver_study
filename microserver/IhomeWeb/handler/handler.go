package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-grpc"
	"image"
	"image/png"
	"net/http"
	DELETESESSION "sss/DeleteSession/proto/example"
	GETAREA "sss/GetArea/proto/example"
	GETIMAGECD "sss/GetImageCd/proto/example"
	GETSESSION "sss/GetSession/proto/example"
	GETSMSCD "sss/GetSmscd/proto/example"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	POSTLOGIN "sss/PostLogin/proto/example"
	POSTRET "sss/PostRet/proto/example"
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

func GetIndex(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {



	// we want to augment the response
	//准备返回给前端的map
	response := map[string]interface{}{
		"errno": "0",
		"errmsg": "ok",
	}
	// encode and write the response as json
	//设置返回数据的格式
	w.Header().Set("Content-Type","application/json")
	//将map转化为json 返回给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetSession(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	fmt.Println(" 登陆检查  GetSession  /api/v1.0/session")
	//获取sessionid
	cookie ,err :=r.Cookie("ihomelogin")
	if err != nil {
		// we want to augment the response
		//准备返回给前端的map
		response := map[string]interface{}{
			"errno": "4101",
			"errmsg": "用户未登录",
		}
		// encode and write the response as json
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}


	//创建 grpc 客户端
	cli :=grpc.NewService()
	//客户端初始化
	cli.Init()

	// call the backend service
	//通过protobuf 生成文件 创建 连接服务端 的客户端句柄
	exampleClient := GETSESSION.NewExampleService("go.micro.srv.GetSession", cli.Client())
	//通过句柄调用服务端函数
	rsp, err := exampleClient.GetSession(context.TODO(), &GETSESSION.Request{
		Sessionid:cookie.Value,
	})
	//判断是否成功
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := make(map[string]string)
	data["name"] = rsp.Name

	// we want to augment the response
	//准备返回给前端的map
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":data,
	}
	// encode and write the response as json
	//设置返回数据的格式
	w.Header().Set("Content-Type","application/json")
	//将map转化为json 返回给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetImageCd(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	fmt.Println("获取验证码图片   GetImageCd  /api/v1.0/imagecode/:uuid")
	//获取url代参的uuid
	uuid := ps.ByName("uuid")

	cli := grpc.NewService()
	//客户端初始化
	cli.Init()
	// call the backend service
	//通过protobuf 生成文件 创建 连接服务端 的客户端句柄
	exampleClient := GETIMAGECD.NewExampleService("go.micro.srv.GetImageCd", cli.Client())
	//通过句柄调用服务端函数
	rsp, err := exampleClient.GetImageCd(context.TODO(), &GETIMAGECD.Request{
		Uuid:uuid,
	})
	//判断是否成功
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//判断返回值 如果不等于 0直接返回错误
	if rsp.Errno != "0"{
		// we want to augment the response
		//准备返回给前端的map
		response := map[string]interface{}{
			"errno": rsp.Errno,
			"errmsg": rsp.Errmsg,
		}
		// encode and write the response as json
		//设置返回数据的格式
		w.Header().Set("Content-Type","application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	//	接受图片数据 拼接图片结构体 将图片发送给前端

	//	接受图片变量
	var img  image.RGBA
	//	pix
	for _, value := range rsp.Pix {

		img.Pix =append(img.Pix , uint8(value))
	}
	//	Stride
	img.Stride = int(rsp.Stride)
	//  point
	img.Rect.Min.X = int(rsp.Min.X)
	img.Rect.Min.Y = int(rsp.Min.Y)
	img.Rect.Max.X = int(rsp.Max.X)
	img.Rect.Max.Y = int(rsp.Max.Y)


	var image captcha.Image

	image.RGBA  = &img
	//将图片发送给前端
	png.Encode(w ,image)


}

func GetSmscd(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	fmt.Println(" 获取短信验证码 GetSmscd  /api/v1.0/smscode/:mobile")

	//获取手机号
	mobile := ps.ByName("mobile")
	////验证手机号
	//myreg := regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)
	//
	//bo := myreg.MatchString(mobile)
	//
	//if bo == false {
	//	//准备返回给前段的map
	//	response := map[string]interface{}{
	//		"errno": utils.RECODE_MOBILEERR,
	//		"errmsg": utils.RecodeText(utils.RECODE_MOBILEERR),
	//	}
	//
	//	//设置返回数据格式
	//	w.Header().Set("Content-Type","application/json")
	//
	//	//将map转化为json返回给前段
	//	if err := json.NewEncoder(w).Encode(response);err != nil  {
	//		http.Error(w,err.Error(),500)
	//		return
	//	}
	//	return
	//}

	//获取图片验证码与uuid
	text := r.URL.Query()["text"][0]
	uuid := r.URL.Query()["id"][0]


	cli := grpc.NewService()
	cli.Init()
	// call the backend service
	exampleClient := GETSMSCD.NewExampleService("go.micro.srv.GetSmscd", cli.Client())
	fmt.Println("sssss")

	rsp, err := exampleClient.GetSmscd(context.TODO(), &GETSMSCD.Request{
		Mobile:mobile,
		Uuid:uuid,
		Text:text,
	})
	fmt.Println("aaaaaaaa")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostRet(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {

	fmt.Println(" 注册服务  PostRet  /api/v1.0/users")
	var request map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&request);err !=nil {
		http.Error(w,err.Error(),500)
		return
	}

	if request["mobile"].(string) == "" ||request["password"].(string) == "" ||request["sms_code"].(string) == "" {
		response := map[string]interface{}{
			"errno": utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}

		w.Header().Set("Content-Type","application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	cli := grpc.NewService()
	cli.Init()
	// call the backend service
	exampleClient := POSTRET.NewExampleService("go.micro.srv.PostRet", cli.Client())
	rsp, err := exampleClient.PostRet(context.TODO(), &POSTRET.Request{
		Mobile:request["mobile"].(string),
		Password:request["password"].(string),
		SmsCode:request["sms_code"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	//设置cookie
	cookie,err := r.Cookie("ihmoelogin")
	if err != nil || cookie.Value == "" {
		cookie := http.Cookie{Name:"ihomelogin",Value:rsp.Sessionid,MaxAge:600,Path:"/"}
		http.SetCookie(w,&cookie)
	}
	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	fmt.Println(" 登陆服务 PostLogin  /api/v1.0/sessions")

	var request map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&request);err != nil {
		http.Error(w,err.Error(),500)
		return
	}

	if request["mobile"].(string) == "" || request["password"].(string) == "" {
		response := map[string]interface{}{
			"errno": utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}

		w.Header().Set("Content-Type","application/json")

		if err := json.NewEncoder(w).Encode(response);err != nil {
			http.Error(w,err.Error(),500)
			return
		}

		return
	}

	cli := grpc.NewService()
	cli.Init()
	// call the backend service
	exampleClient := POSTLOGIN.NewExampleService("go.micro.srv.PostLogin", cli.Client())
	rsp, err := exampleClient.PostLogin(context.TODO(), &POSTLOGIN.Request{
		Mobile:request["mobile"].(string),
		Password:request["password"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//将获取到的sessionid 设置到 cookie

	cookie ,err := r.Cookie("ihomelogin")

	if err != nil || cookie.Value == "" {
		cookie := http.Cookie{Name:"ihomelogin",Value:rsp.Sessionid,MaxAge:600,Path:"/"}

		http.SetCookie(w,&cookie)
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func DeleteSession(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	beego.Info("---------------- DELETE  /api/v1.0/session Deletesession() ------------------")

	cli := grpc.NewService()
	cli.Init()
	// call the backend service
	exampleClient := DELETESESSION.NewExampleService("go.micro.srv.DeleteSession", cli.Client())
	//获取session
	cookie,err := r.Cookie("ihomelogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno": utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type","application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w,err.Error(),503)
			beego.Info(err)
			return
		}
		return
	}
	rsp, err := exampleClient.DeleteSession(context.TODO(), &DELETESESSION.Request{
		Sessionid : cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		beego.Info(err)
		return
	}

	cookie ,err := r.Cookie("ihomelogin")

	if err != nil || cookie.Value == "" {
		return
	}else {
		cookie := http.Cookie{Name:"ihomelogin",Path:"/",MaxAge:-1}
		http.SetCookie(w,&cookie)
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	w.Header().Set("Content-Type","application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 503)
		beego.Info(err)
		return
	}
	return
}
