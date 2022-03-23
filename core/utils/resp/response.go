package resp

import (
	"github.com/gofiber/fiber/v2"
)

// Response api response
type Response struct {
	Code  int         `json:"code" comment:"111"`        // msg
	Msg   string      `json:"msg"`                       // code
	Data  interface{} `json:"data,omitempty" form:"111"` // data
	Count int         `json:"count,omitempty"`           // data count
}

// JSON gin resp to json
func JSON(c *fiber.Ctx, code int, data ...interface{}) error {
	resp := Response{
		Code: code,
		Msg:  GetMsg(code),
		Data: data[0],
	}
	if len(data) == 2 {
		resp.Count = data[1].(int)
	}
	c.Set("statuscode", "10000")
	return c.JSON(resp)
}