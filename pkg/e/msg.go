package e

import "net/http"

var Msg = map[int]string{
	http.StatusOK:                  "OK",
	http.StatusBadRequest:          "Bad Request",
	http.StatusInternalServerError: "Operation failed",
}

func GetMsg(code int) string {
	msg, ok := Msg[code]
	if !ok {
		return Msg[http.StatusInternalServerError]
	}
	return msg
}
