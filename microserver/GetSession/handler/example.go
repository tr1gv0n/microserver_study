package handler

import (
	"context"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sss/IhomeWeb/utils"

	example "sss/GetSession/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetSession(ctx context.Context, req *example.Request, rsp *example.Response) error {
	fmt.Println(" 登陆检查  GetSession  /api/v1.0/session")
	/*1 初始化返回值*/
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno )

	//2 连接redis
	bm,err := utils.RedisOpen(utils.G_server_name,utils.G_redis_addr,utils.G_redis_port,utils.G_redis_dbnum)
	if err != nil {
		fmt.Println("连接redis失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno )
		return nil
	}

	//3 接受web发送过来的sessionid
	sessionid := req.Sessionid

	//4 拼接对应的key
	namekey := sessionid + "name"

	//5 查询name
	value := bm.Get(namekey)
	name,_ := redis.String(value,nil)

	//6 将name返回
	rsp.Name = name

	return nil
}

