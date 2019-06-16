package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"

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

	//连接redis

	bm,err := utils.RedisOpen(utils.G_server_name,utils.G_redis_addr,utils.G_redis_port,utils.G_redis_dbnum)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	key := "area_info"

	area_info_value := bm.Get(key)
	//接受数据
	var areas []models.Area
	if area_info_value  != nil{
		fmt.Println("获取数据送给web")
		err = json.Unmarshal(area_info_value.([]byte),&areas)

		for key ,value := range areas{
			fmt.Println(key,value)

			area := example.ResponseAddress{Aid:int32(value.Id),Aname:string(value.Name)}

			rsp.Data = append(rsp.Data,&area)
		}
		return nil
	}
	//查询数据库
	o:= orm.NewOrm()

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
	//2.5  存入数据
	//json编码
	area_info_json,_ := json.Marshal(areas)
	//存入数据
	err= bm.Put(key,area_info_json,time.Second*7200)
	if err != nil {
		fmt.Println("redis 存入数据失败",err)
		rsp.Errno= utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//3  将查询到的数据转化类型
	for key,value := range areas{
		fmt.Println(key,value)
		//	结构体 ---》proto
		area := example.ResponseAddress{Aid:int32(value.Id),Aname:string(value.Name)}
		//4 数据返回
		rsp.Data = append(rsp.Data,&area)
	}
	return nil
}


