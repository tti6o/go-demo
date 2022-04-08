package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Response struct {
	Error
	RequestId string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

type Error struct {
	Prefix       string `json:"-"`        // 错误码前缀，有时候需要使用App.10000的错误码形式 App.就为错误码前缀
	ErrCode      string `json:"err_code"` // 错误码
	ErrMsg       string `json:"err_msg"`  // 详细错误信息，暴露给前端
	PrimitiveErr error  `json:"-"`        // error包装的原始错误信息或者debug信息，不暴露给前端，只用%+v打印在日志中
}

type ToMyOAListResp struct {
	Code int        `json:"code"` //myoa错误码
	Data []ListItem `json:"data"` //data
}

type ListItem struct {
	ID string `json:"id"` //单据ID
	WorkItem
}

type WorkItem struct {
	Category      string `json:"category"`        // 业务类别
	ProcessName   string `json:"process_name"`    // 流程名称
	ProcessInstId string `json:"process_inst_id"` // 流程实例名称
	Title         string `json:"title"`           // 标题
	Applicant     string `json:"applicant"`       // 申请人的英文名
	Handler       string `json:"handler"`         // 当前处理人的英文名
	Activity      string `json:"activity"`        // 此审批单据对应的业务的审批活动或步骤
}

func TestJsonTag(t *testing.T) {
	//var response Response
	var jsonStr string
	jsonStr = `{"code":200,"data":[]}`
	//jsonStr = `{"code":500,"data":[]}`
	//jsonStr = `{ "err_code": "00000", "err_msg": "Success", "request_id": "c705ji98d3b44d48s0t0", "data": { "msg": "hello test" } }`
	//err := json.Unmarshal([]byte(jsonStr), &response)
	//fmt.Printf("err:%v\n", err)
	//fmt.Printf("response:%#v\n", response)
	resp := ToMyOAListResp{}
	//var resp []ListItem
	response2 := Response{Data: &resp}
	err2 := json.Unmarshal([]byte(jsonStr), &response2)
	fmt.Printf("err2:%v\n", err2)
	fmt.Printf("response2:%#v\n", response2)
}
