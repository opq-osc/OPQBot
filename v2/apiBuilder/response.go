package apiBuilder

import (
	"encoding/json"
)

type ResponseStruct struct {
	CgiBaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"CgiBaseResponse"`
	ResponseData json.RawMessage `json:"ResponseData,omitempty"`
	Data         interface{}     `json:"Data"`
}
type Response struct {
	originMsg []byte
	response  *ResponseStruct
}

type GroupMessageResponse struct {
	MsgSeq  int64 `json:"MsgSeq"`
	MsgTime int64 `json:"MsgTime"`
}

func NewResponse(msg []byte) *Response {
	return &Response{originMsg: msg}
}
func (r *Response) unmarshal() error {
	if r.response == nil {
		r.response = &ResponseStruct{}
		err := json.Unmarshal(r.originMsg, r.response)
		return err
	}
	return nil
}
func (r *Response) GetData(data interface{}) error {
	if err := r.unmarshal(); err != nil {
		return err
	}
	return json.Unmarshal(r.response.ResponseData, data)
}
func (r *Response) Ok() bool {
	if err := r.unmarshal(); err != nil {
		return false
	}
	return r.response.CgiBaseResponse.Ret == 0
}
func (r *Response) ErrorMsg() string {
	if err := r.unmarshal(); err != nil {
		return ""
	}
	return r.response.CgiBaseResponse.ErrMsg
}
func (r *Response) Result() (Ret int, ErrMsg string) {
	if err := r.unmarshal(); err != nil {
		return -1, err.Error()
	}
	return r.response.CgiBaseResponse.Ret, r.response.CgiBaseResponse.ErrMsg
}
func (r *Response) GetOrigin() []byte {
	return r.originMsg
}

func (r *Response) GetGroupMessageResponse() (*GroupMessageResponse, error) {
	if err := r.unmarshal(); err != nil {
		return nil, err
	}
	data := &GroupMessageResponse{}
	err := json.Unmarshal(r.response.ResponseData, data)
	return data, err
}
