package response

import (
	"encoding/json"
	log "forum/forum/pkg/logger"
	"net/http"
)

type Response struct {
	Status int
	Body   interface{}
}

func SendResponse(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	if body != nil {
		err := json.NewEncoder(w).Encode(body)
		if err != nil {
			log.Error("send_response err =", err)
			return
		}
	}
}
