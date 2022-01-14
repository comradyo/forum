package response

import (
	"encoding/json"
	routing "github.com/qiangxue/fasthttp-routing"
)

type Response struct {
	Status int
	Body   interface{}
}

func SendResponse(ctx *routing.Context, status int, body interface{}) {
	ctx.Response.SetStatusCode(status)
	ctx.Response.Header.Set("Content-Type", "application/json")
	if body != nil {
		err := json.NewEncoder(ctx.Response.BodyWriter()).Encode(body)
		if err != nil {
			return
		}
	}
}
