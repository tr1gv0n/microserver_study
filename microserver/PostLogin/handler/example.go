package handler

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/orm"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	"strconv"
	"time"

	example "sss/PostLogin/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostLogin(ctx context.Context, req *example.Request, rsp *example.Response) error {
	fmt.Println(" 登陆服务 PostLogin  /api/v1.0/sessions")

	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	/*2 接受数据 */
	var user models.User
	/*3 查询数据 */
	o := orm.NewOrm()
	qs := o.QueryTable("user")

	err := qs.Filter("mobile",req.Mobile).One(&user)
	if err != nil {
		fmt.Println("用户名查询失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/*4 判断密码*/
	if utils.Getmd5string(req.Password) != user.Password_hash{
		fmt.Println("密码错误",utils.Getmd5string(req.Password), user.Password_hash)
		rsp.Errno = utils.RECODE_PWDERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/*5 生成sessionid*/
	sessionid := utils.Getmd5string(req.Mobile+req.Password+strconv.Itoa(int(time.Now().UnixNano())))

	/*6 连接redis*/
	bm,err := utils.RedisOpen(utils.G_server_name,utils.G_redis_addr,
		utils.G_redis_port,utils.G_redis_dbnum)

	if err != nil {
		fmt.Println("连接redis失败",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil

	}

	/*7 拼接key 存入数据*/
	bm.Put(sessionid+"name",user.Name,time.Second*600)
	bm.Put(sessionid+"mobile",user.Mobile,time.Second*600)
	bm.Put(sessionid+"user_id",user.Id,time.Second*600)
	/*8  返回sessionid*/
	rsp.Sessionid = sessionid

	return nil
}

