package handler

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/orm"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	"strconv"
	"github.com/garyburd/redigo/redis"

	example "sss/GetUserInfo/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetUserInfo(ctx context.Context, req *example.Request, rsp *example.Response) error {
	fmt.Println("获取用户信息 GetUserInfo /api/v1.0/user")

	/*1 初始化返回值 */
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno )

	/*2 获取sessionid */
	sessionid := req.Sessionid

	/*3 连接redis*/
	bm ,err :=utils.RedisOpen(utils.G_server_name ,utils.G_redis_addr ,utils.G_redis_port ,utils.G_redis_dbnum)
	if err!=nil{
		fmt.Println("redis 连接失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno )
		return nil
	}

	/*4 拼接可以  查询 user_id*/
	value :=bm.Get(sessionid+"user_id")
	//	[]uint8
	value_string ,_:=redis.String(value,nil)
	//	string
	value_int ,_:= strconv.Atoi(value_string)
	//	int


	/*5 连接数据库*/
	o := orm.NewOrm()
	user := models.User{Id:value_int}

	/*6查询数据*/
	err =o.Read(&user)
	if err !=nil{
		fmt.Println("用户数据查询失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno )
		return nil
	}

	/*7返回数据*/

	//"user_id": 1,
	//"name": "Panda",
	//"mobile": "110",
	//"real_name": "熊猫",
	//"id_card": "210112244556677",
	//"avatar_url":

	rsp.UserId = strconv.Itoa(user.Id)
	rsp.Name = user.Name
	rsp.Mobile = user.Mobile
	rsp.RealName = user.Real_name
	rsp.IdCard = user.Id_card
	rsp.AvatarUrl = user.Avatar_url



	return nil
}

