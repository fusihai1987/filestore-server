package common

import (
	"encoding/json"
	"fmt"
)

type Resp struct {
	Code int
	Msg  string
	Data interface{}
}

func NewResp(code int, msg string, data interface{}) Resp{
	return Resp{
		Code: code,
		Msg: msg,
		Data: data,
	}
}

func (resp *Resp) JsonBytes() []byte{
	jsonRes, err := json.Marshal(resp)

	if err != nil {
		fmt.Println(err.Error())
	}

	return jsonRes

}

func (resp *Resp) JsonString() string{
	jsonRes, err := json.Marshal(resp)

	if err != nil {
		fmt.Println(err.Error())
	}

	return string(jsonRes)

}