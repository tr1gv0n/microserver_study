package handler

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/orm"
	"sss/IhomeWeb/utils"
	"time"
	"github.com/garyburd/redigo/redis"
	"sss/IhomeWeb/models"

	example "sss/PostUserAuth/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostUserAuth(ctx context.Context, req *example.Request, rsp *example.Response) error {
	fmt.Println(" 实名认证服务  PostUserAuth   /api/v1.0/user/auth  ")
	/*1 初始化返回值*/
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno )


	/*TODO   关联到相关部门 平台   身份证号 和 姓名 进行实名认证  返回1个对或错*/



	/*2连接redis*/
	bm ,err :=utils.RedisOpen(utils.G_server_name,utils.G_redis_addr ,utils.G_redis_port,utils.G_redis_dbnum)
	if err!=nil{
		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(rsp.Errno )
		return nil
	}

	/*3 拼接key*/
	key :=req.Sessionid +"user_id"

	/*4 查询userid*/
	value :=bm.Get(key)

	value_int ,_:=redis.Int(value ,nil)

	fmt.Println("user_id",value_int)






	/*5 连接数据库*/
	o := orm.NewOrm()

	user:= models.User{Id:value_int ,Real_name: req.Realname ,Id_card: req.Idcard}

	/*6 更新 数据*/
	_,err =o.Update(&user ,"real_name","id_card")
	if err!=nil{

		fmt.Println("实名认证数据 更新失败 ",err)
		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(rsp.Errno )
		return nil
	}



	/*7更新session信息时间*/
	bm.Put(req.Sessionid +"user_id",user.Id ,time.Second*600)
	bm.Put(req.Sessionid +"mobile",user.Mobile ,time.Second*600)
	bm.Put(req.Sessionid +"name",user.Name ,time.Second*600)

	return nil
}
