package handler

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"path"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"

	example "sss/PostAvatar/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostAvatar(ctx context.Context, req *example.Request, rsp *example.Response) error {
	fmt.Println(" 上传头像  PostAvatar  /api/v1.0/user/avatar ")
	/*1 初始化返回值*/
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg =utils.RecodeText(rsp.Errno)

	//2 数据对比
	if req.Filesize != int64(len(req.Buffer)) {
		fmt.Println("接受数据异常")
		rsp.Errno = utils.RECODE_IOERR
		rsp.Errmsg =utils.RecodeText(rsp.Errno)
		return nil
	}

	//上传图片到fastdfs
	ext := path.Ext(req.Filename)

	fileid,err := utils.Uploadbybuf(req.Buffer,ext[1:])
	if err != nil {
		fmt.Println("文件上传fastdfs失败")
		rsp.Errno = utils.RECODE_IOERR
		rsp.Errmsg =utils.RecodeText(rsp.Errno)
		return nil
	}

	//4 连接redis
	bm, err := utils.RedisOpen(utils.G_server_name,utils.G_redis_addr,utils.G_redis_port,
		utils.G_redis_dbnum)

	if err!=nil{
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg =utils.RecodeText(rsp.Errno)
		return nil
	}

	//5 拼接key
	key  := req.Sessionid +"user_id"

	//6 查询user_id
	value := bm.Get(key)
	value_int ,_ := redis.Int(value,nil)

	//7 注册mysql
	o := orm.NewOrm()
	user := models.User{Id:value_int,Avatar_url:fileid}

	//8 将数据更新到表中
	_,err = o.Update(&user,"avatar_url")
	if err !=nil{
		fmt.Println("更新数据库 文件地址失败 ",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg =utils.RecodeText(rsp.Errno)
		return nil
	}
	rsp.Fileid = fileid
	return nil
}

