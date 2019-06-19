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
	GETUSERINFO "sss/GetUserInfo/proto/example"
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
	fmt.Println("退出登陆  DeleteSession   /api/v1.0/session")

	//从cookie当中获取sessionid
	cookie ,err :=r.Cookie("ihomelogin")
	if err!=nil || cookie.Value ==""{

		response := map[string]interface{}{
			"errno": utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
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
	exampleClient := DELETESESSION.NewExampleService("go.micro.srv.DeleteSession", cli.Client())
	//通过句柄调用服务端函数
	rsp, err := exampleClient.DeleteSession(context.TODO(), &DELETESESSION.Request{
		Sessionid: cookie.Value,
	})
	//判断是否成功
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	/*删除cookie当中的sessionid*/

	if rsp.Errno =="0"{
		cookie ,err :=r.Cookie("ihomelogin")
		if err!=nil ||  cookie.Value==""{
			return
		}else {
			cookie:=  http.Cookie{Name:"ihomelogin",Value:"",Path:"/",MaxAge:-1}
			http.SetCookie(w,&cookie)
		}

		// we want to augment the response
		//准备返回给前端的map
		response := map[string]interface{}{
			"errno": rsp.Errno,
			"errmsg":rsp.Errmsg,
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




	// we want to augment the response
	//准备返回给前端的map
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg":rsp.Errmsg,
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

func GetUserInfo(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	fmt.Println("GetUserInfo  获取用户信息   /api/v1.0/user")

	cli := grpc.NewService()
	cli.Init()
	// call the backend service
	exampleClient := GETUSERINFO.NewExampleService("go.micro.srv.GetUserInfo", cli.Client())
	cookie,err := r.Cookie("ihomelogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno": utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}
	rsp, err := exampleClient.GetUserInfo(context.TODO(), &GETUSERINFO.Request{
		Sessionid:cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 准备1个数据的map
	data := make(map[string]interface{})
	//将信息发送给前端
	data["user_id"] = rsp.UserId
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IdCard
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)
	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data" : data,
	}

	w.Header().Set("Content-Type", "application/json")


	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostAvatar(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	fmt.Println(" 上传头像  PostAvatar  /api/v1.0/user/avatar ")

	//接受前段发送过来的二进制图片数据
	File ,FileHeader ,err := r.FormFile("avatar")
	if err != nil {
		response := map[string]interface{}{
			"errno": utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type","application/json")

		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	//获取文件的基本信息
	fmt.Println("文件名",FileHeader.Filename)
	fmt.Println("文件大小",FileHeader.Size)

	//创建1个二进制空间
	filebuffer := make([]byte,FileHeader.Size)

	//把接受过来的文件读入二进制空间
	_,err = File.Read(filebuffer)
	if err != nil {
		response := map[string]interface{}{
			"errno": utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type","application/json")

		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	//获取sessionid
	cookie,err := r.Cookie("ihomelogin")
	if err != nil || cookie.Value == "" {
		response := map[string]interface{}{
			"errno": utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type","application/json")

		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

	}

	cli := grpc.NewService()
	cli.Init()
	// call the backend service
	exampleClient := POSTAVATAR.NewExampleService("go.micro.srv.PostAvatar", cli.Client())
	rsp, err := exampleClient.PostAvatar(context.TODO(), &POSTAVATAR.Request{
		Sessionid:cookie.Value,
		Buffer:filebuffer,
		Filename:FileHeader.Filename,
		Filesize:FileHeader.Size,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}


	data := make(map[string]string)
	data["avatar_url"] = utils.AddDomain2Url(rsp.Fileid)
	fmt.Println(data)
	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":data,
	}
	w.Header().Set("Content-Type","application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

/*func PutUserInfo(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
	fmt.Println("更新用户名   PutUserInfo   /api/v1.0/user/name")

	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//参数校验
	if request["name"].(string) == "" {
		response := map[string]interface{}{
			"errno": utils.RECODE_NODATA,
			"errmsg": utils.RecodeText(utils.RECODE_NODATA),
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

	//获取sessionid
	cookie,err := r.Cookie("ihomelogin")
	if err != nil {
		response := map[string]interface{}{
			"errno": utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
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

	cli := grpc.NewService()
	cli.Init()
	// call the backend service
	exampleClient := PUTUSERINFO.NewExampleService("go.micro.srv.template", cli.Client())
	rsp, err := exampleClient.PutUserInfo(context.TODO(), &PUTUSERINFO.Request{
		Name: request["name"].(string),
		Sessionid:cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//重新设置 cookie的时间
	//登陆的时候 600妙
	//时间消耗  300
	//进行了操作 刷新世界 将时间再次恢复到600
	cookienew := http.Cookie{Name:"ihomelogin",Value:cookie.Value,Path:"/",MaxAge:600}

	http.SetCookie(w,&cookienew)

	data :=make(map[string]string)
	data["name"] = rsp.Name
	// we want to augment the response
	response := map[string]interface{}{
		"errno": rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":data,
	}

	// encode and write the response as json
	w.Header().Set("Content-Type","application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
*/
/*func GetUserAuth(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {

	fmt.Println("获取用户信息 GetUserAuth /api/v1.0/user/auth")
	//获取cookie当中的sessionid
	cookie,err :=r.Cookie("ihomelogin")
	if err!=nil||cookie.Value==""{
		response := map[string]interface{}{
			"errno": utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
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
	exampleClient := GETUSERINFO.NewExampleService("go.micro.srv.GetUserInfo", cli.Client())
	//通过句柄调用服务端函数
	rsp, err := exampleClient.GetUserInfo(context.TODO(), &GETUSERINFO.Request{
		Sessionid:cookie.Value,
	})
	//判断是否成功
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//接受数据

	//"user_id": 1,
	//"name": "Panda",
	//"mobile": "110",
	//"real_name": "熊猫",
	//"id_card": "210112244556677",
	//"avatar_url":
	data := make(map[string]interface{})
	data["user_id"] = rsp.UserId
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IdCard
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)

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
*/
//func PostUserAuth(w http.ResponseWriter, r *http.Request,ps httprouter.Params) {
//	fmt.Println(" 实名认证服务  PostUserAuth   /api/v1.0/user/auth  ")
//
//
//	// decode the incoming request as json
//	//接受 前端发送过来数据的
//	var request map[string]interface{}
//	// 将前端 json 数据解析到 map当中
//	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
//		http.Error(w, err.Error(), 500)
//		return
//	}
//	/*  real_name: "熊猫",
//    id_card:
//	*/
//
//	//数据校验
//	if request["real_name"].(string) ==""|| request["id_card"].(string) ==""{
//		response := map[string]interface{}{
//			"errno": utils.RECODE_NODATA,
//			"errmsg": utils.RecodeText(utils.RECODE_NODATA),
//		}
//		// encode and write the response as json
//		//设置返回数据的格式
//		w.Header().Set("Content-Type","application/json")
//		//将map转化为json 返回给前端
//		if err := json.NewEncoder(w).Encode(response); err != nil {
//			http.Error(w, err.Error(), 500)
//			return
//		}
//		return
//	}
//	/*TODO  需要对身份证号 进行 正则 校验 */
//
//
//	//获取sessionid
//	cookie ,err :=r.Cookie("ihomelogin")
//	if err!=nil||cookie.Value==""{
//		response := map[string]interface{}{
//			"errno": utils.RECODE_SESSIONERR,
//			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
//		}
//		// encode and write the response as json
//		//设置返回数据的格式
//		w.Header().Set("Content-Type","application/json")
//		//将map转化为json 返回给前端
//		if err := json.NewEncoder(w).Encode(response); err != nil {
//			http.Error(w, err.Error(), 500)
//			return
//		}
//		return
//	}
//
//	//创建 grpc 客户端
//	cli :=grpc.NewService()
//	//客户端初始化
//	cli.Init()
//
//	// call the backend service
//	//通过protobuf 生成文件 创建 连接服务端 的客户端句柄
//	exampleClient := POSTUSERAUTH.NewExampleService("go.micro.srv.PostUserAuth", cli.Client())
//	//通过句柄调用服务端函数
//	rsp, err := exampleClient.PostUserAuth(context.TODO(), &POSTUSERAUTH.Request{
//		Idcard:request["id_card"].(string),
//		Realname:request["real_name"].(string)  ,
//		Sessionid:cookie.Value,
//	})
//	//判断是否成功
//	if err != nil {
//		http.Error(w, err.Error(), 500)
//		return
//	}
//
//
//	//刷新cookie的时间
//	cookienew :=http.Cookie{Name:"ihomelogin", Value:cookie.Value ,Path:"/",MaxAge:600}
//
//	http.SetCookie(w,&cookienew)
//
//
//	// we want to augment the response
//	//准备返回给前端的map
//	response := map[string]interface{}{
//		"errno": rsp.Errno,
//		"errmsg": rsp.Errmsg,
//	}
//	// encode and write the response as json
//	//设置返回数据的格式
//	w.Header().Set("Content-Type","application/json")
//	//将map转化为json 返回给前端
//	if err := json.NewEncoder(w).Encode(response); err != nil {
//		http.Error(w, err.Error(), 500)
//		return
//	}
//}
