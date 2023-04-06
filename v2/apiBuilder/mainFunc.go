package apiBuilder

type IMainFunc interface {
	SendMsg() ISendMsg
	QueryUin() IQueryUin
	Qrcode() IQrcode
	Upload() IUpload
}
