package handler

import (
	"context"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"reflect"
	"sss/IhomeWeb/utils"
	"strconv"
	"time"

	example "sss/GetSmscd/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetSmscd(ctx context.Context, req *example.Request, rsp *example.Response) error {
	fmt.Println(" 获取短信验证码 GetSmscd  /api/v1.0/smscode/:mobile")

	//1 初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

/*	//2 手机号 验证 带入 数据库  如果手机号存在说明 是老用户

	//创建数据库句柄
	o := orm.NewOrm()

	user :=models.User{Mobile:req.Mobile}

	err := o.Read(&user)
	if err ==nil{
		fmt.Println("该用户已经注册",req.Mobile  )
		rsp.Errno = utils.RECODE_USERONERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
*/
	//3 连接redis
	bm ,err :=utils.RedisOpen(utils.G_server_name,utils.G_redis_addr,
		utils.G_redis_port,utils.G_redis_dbnum)
	if err!=nil{
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//4 获取图片验证码
	value :=bm.Get(req.Uuid)

	fmt.Println(reflect.TypeOf(value),value)

	value_string ,_:=redis.String(value ,nil)

	//5 数据比对 验证图片验证码是否正确
	if req.Text !=value_string{
		fmt.Println("图片验证码错误 ",req.Text , value_string )
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//6 生成随机数
	t := rand.New(rand.NewSource(time.Now().UnixNano()))

	//四位 1000-9999
	size := t.Intn(8999)+1000
	fmt.Println("随机数：",size)

	/*//7调用短信借口发送短信

	//发送短信的配置信息
	messageconfig := make(map[string]string)
	//预先创建好的appid
	messageconfig["appid"] = "29672"
	//预先获得的app的key
	messageconfig["appkey"] = "89d90165cbea8cae80137d7584179bdb"
	//加密方式默认
	messageconfig["signtype"] = "md5"



	//messagexsend
	//创建短信发送的句柄
	messagexsend := submail.CreateMessageXSend()
	//短信发送的手机号
	submail.MessageXSendAddTo(messagexsend, req.Mobile)
	//短信发送的模板
	submail.MessageXSendSetProject(messagexsend, "NQ1J94")
	//验证码
	submail.MessageXSendAddVar(messagexsend, "code", strconv.Itoa(size))
	//发送短信的请求
	send :=submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig)

	fmt.Println("MessageXSend ", send)

	//8验证短信验证码是否成功
	bo :=strings.Contains(send,"success")
	if bo !=true{
		fmt.Println("短信验证码发送失败 ！" )
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}*/

	//9 将随机数 与手机号存入 redis
	err =bm.Put(req.Mobile,strconv.Itoa(size),time.Second*300)
	if err!=nil{
		fmt.Println("随机数存储失败 ！" )
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	return nil
}

