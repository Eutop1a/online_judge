package response

type Response struct {
	Code int `json:"code"`
}
type ResponseWithData struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
