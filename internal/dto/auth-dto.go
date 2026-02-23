package dto

type VerifyReq struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

type VerifyRes struct {
	Token   string `json:"token"`
	Address string `json:"address"`
}
