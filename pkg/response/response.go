package response

import (
	"encoding/json"

	"net/http"
)

type Response struct {
	Status int
	Body   interface{}
}

func SendResponse(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		err := json.NewEncoder(w).Encode(body)
		if err != nil {
			return
		}
	}
}
