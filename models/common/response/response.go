package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	CtxUserIDKey   = "userID"
	CtxUserNameKey = "username"
)

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"` // omitempty 字段为空就忽略
}

func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	//responseData := make(map[string]interface{})
	//val := reflect.ValueOf(data)
	//if val.Kind() == reflect.Struct {
	//	for i := 0; i < val.NumField(); i++ {
	//		field := val.Field(i)
	//		switch field.Kind() {
	//		case reflect.Int64:
	//			responseData[reflect.TypeOf(data).Field(i).Name] = strconv.FormatInt(field.Int(), 10)
	//		default:
	//			responseData[reflect.TypeOf(data).Field(i).Name] = field.Interface()
	//		}
	//	}
	//} else if val.Kind() == reflect.Int64 {
	//	uidString := strconv.FormatInt(val.Int(), 10)
	//	//jsonData, _ = json.Marshal(map[string]string{"uid": uidString})
	//}
	switch d := data.(type) {
	case int64:
		uidString := strconv.FormatInt(d, 10)
		c.JSON(http.StatusOK, &ResponseData{
			Code: CodeSuccess,
			Msg:  CodeSuccess.Msg(),
			Data: uidString,
		})

	default:
		c.JSON(http.StatusOK, &ResponseData{
			Code: CodeSuccess,
			Msg:  CodeSuccess.Msg(),
			Data: data,
		})
	}

}
