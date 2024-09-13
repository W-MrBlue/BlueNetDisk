package api

import (
	"BlueNetDisk/pkg/ctl"
	"encoding/json"
)

func ErrorResponse(err error) *ctl.Response {
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return ctl.RespError(err)
	}
	return ctl.RespError(err)
}
