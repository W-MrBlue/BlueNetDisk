package ctl

import "BlueNetDisk/pkg/e"

type Response struct {
	Status int
	Msg    string
	Data   interface{}
	Err    string
}

type Datalist struct {
	Item  interface{} `json:"items"`
	Total int64       `json:"total"`
}

func RespSuccess(code ...int) *Response {
	status := e.Success
	if code != nil {
		status = code[0]
	}
	return &Response{
		Status: status,
		Msg:    e.GetMsg(status),
	}
}
func RespSuccessWithData(data interface{}, code ...int) *Response {
	status := e.Success
	if code != nil {
		status = code[0]
	}
	return &Response{
		Status: status,
		Msg:    e.GetMsg(status),
		Data:   data,
	}
}

func RespList(items interface{}, total int64) *Response {
	return &Response{
		Status: e.Success,
		Data: Datalist{
			Item:  items,
			Total: total,
		},
		Msg: e.GetMsg(e.Success),
	}
}

func RespError(err error, code ...int) *Response {
	status := e.InternalServerError
	if code != nil {
		status = code[0]
	}
	return &Response{
		Status: status,
		Msg:    e.GetMsg(status),
		Err:    err.Error(),
	}
}
