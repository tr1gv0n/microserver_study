package handler

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	"strconv"
	"time"

	example "sss/PostRet/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostRet(ctx context.Context, req *example.Request, rsp *example.Response) error {
	fmt.Println(" 注册服务  PostRet  /api/v1.0/users")
	//1初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//2连接redis
	bm,err := utils.RedisOpen(utils.G_server_name,utils.G_redis_addr,utils.G_redis_port,utils.G_redis_dbnum)
	if err !=nil {
		fmt.Println("redis 连接失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//3 获取到短信验证码
	value := bm.Get(req.Mobile)

	value_string ,_ := redis.String(value,nil)
	//4验证短信验证码是否正确
	if value_string != req.SmsCode {
		fmt.Println("短信验证码错误",value_string,req.SmsCode)
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//5 加密
	user := models.User{}
	user.Password_hash = utils.Getmd5string(req.Password)

	//6 插入数据到数据库中
	user.Mobile = req.Mobile
	user.Name = req.Mobile

	o:=orm.NewOrm()
	id ,err:= o.Insert(&user)
	if err != nil {
		fmt.Println("用户数据注册插入失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//7生成sessionid
	sessionid := utils.Getmd5string(req.Mobile+req.Password+strconv.Itoa(int(time.Now().UnixNano())))

	rsp.Sessionid = sessionid
	//8通过sessionid将数据存入redis
	bm.Put(sessionid+"user_id",id ,time.Second*600)
	bm.Put(sessionid+"mobile",user.Mobile ,time.Second*600)
	bm.Put(sessionid+"name",user.Name,time.Second*600)

	return nil
}
