package apiBuilder

type IQueryUin interface {
	DoApi
	SetUin(uid string)
}
type QueryUinStruct struct {
	Uin       int    `json:"Uin"`
	Uid       string `json:"Uid"`
	Nick      string `json:"Nick"`
	Head      string `json:"Head"`
	Signature string `json:"Signature"`
	Sex       int    `json:"Sex"`
	Level     int    `json:"Level"`
}

func (b *Builder) QueryUin() IQueryUin {
	cmd := "QueryUinByUid"
	b.CgiCmd = &cmd
	return b
}
func (b *Builder) SetUin(uid string) {
	if b.CgiRequest == nil {
		b.CgiRequest = &CgiRequest{}
	}
	b.CgiRequest.Uid = &uid
}
