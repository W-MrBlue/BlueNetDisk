package api

import (
	"BlueNetDisk/pkg/ctl"
	"encoding/json"
)

// ErrorResponse is tended to classify error into json.UnmarshalTypeError and others
func ErrorResponse(err error) *ctl.Response {
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return ctl.RespError(err)
	}
	return ctl.RespError(err)
}
