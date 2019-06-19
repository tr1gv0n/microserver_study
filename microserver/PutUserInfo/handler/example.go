package handler

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	"time"

	example "sss/PutUserInfo/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PutUserInfo(ctx context.Context, req *example.Request, rsp *example.Response) error {
	fmt.Println("更新用户名   PutUserInfo   /api/v1.0/user/name")

	//1 初始化返回之
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//2 连接redis
	bm, err := utils.RedisOpen(utils.G_server_name,utils.G_redis_addr,utils.G_redis_port,utils.G_redis_dbnum)

	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//3 拼接key
	key:= req.Sessionid+"user_id"

	//4 查询user_id
	value := bm.Get(key)
	value_int,_ := redis.Int(value,nil)
	fmt.Println("userid",value_int)

	//5 连接mysql
	o := orm.NewOrm()
	user := models.User{Id:value_int,Name:req.Name}

	//6 更新用户名
	size ,err := o.Update(&user,"name")
	if err != nil {
		fmt.Println("更新用户名失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	fmt.Println("返回之",size)

	//7 更新session信息
	bm.Put(req.Sessionid+"name",user.Name,time.Second*600)
	bm.Put(req.Sessionid+"user_id",user.Id,time.Second*600)
	bm.Put(req.Sessionid+"mobile",user.Mobile,time.Second*600)

	rsp.Name = req.Name
	return nil
}
