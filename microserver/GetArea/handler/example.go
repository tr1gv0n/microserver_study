package handler

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/orm"

	"sss/IhomeWeb/utils"
	"sss/IhomeWeb/models"
	example "sss/GetArea/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetArea (ctx context.Context, req *example.Request, rsp *example.Response) error {
	fmt.Println("获取地域信息服务  GetArea  /api/v1.0/areas")
	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//查询数据库
	o:= orm.NewOrm()
	var areas []models.Area
	qs := o.QueryTable("Area")

	num ,err := qs.All(&areas)
	if err != nil {
		fmt.Println("查询数据库错误",err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	if num == 0 {
		fmt.Println("五数据",err)
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	for key,value := range areas{
		fmt.Println(key,value)

		area := example.ResponseAddress{Aid:int32(value.Id),Aname:string(value.Name)}
		rsp.Data = append(rsp.Data,&area)
	}
	return nil
}


