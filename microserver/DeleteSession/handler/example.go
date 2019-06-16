package handler

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"reflect"
	"sss/IhomeWeb/utils"

	example "sss/DeleteSession/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) DeleteSession(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info(" DELETE session    /api/v1.0/session !!!")

	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	beego.Info(req.Sessionid,reflect.TypeOf(req.Sessionid))

	bm,err :=utils.RedisOpen(utils.G_server_name,utils.G_redis_addr,
		utils.G_redis_port,utils.G_redis_dbnum)
	if err != nil {
		fmt.Println("连接redis失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	sessionidname :=  req.Sessionid + "name"
	sessioniduserid :=  req.Sessionid + "user_id"
	sessionidmobile :=  req.Sessionid + "mobile"

	//从缓存中获取session 那么使用唯一识别码
	bm.Delete(sessionidname)
	bm.Delete(sessioniduserid)
	bm.Delete(sessionidmobile)
	return nil
}

