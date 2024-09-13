package ctl

import "BlueNetDisk/pkg/e"

type Response struct {
	Status int
	Msg    string
	Data   interface{}
	Err    string
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
