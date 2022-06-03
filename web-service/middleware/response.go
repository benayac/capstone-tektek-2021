package middleware

import (
	"encoding/json"
	"net/http"
)

type BaseResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (r *BaseResponse) FailResponse(message string) {
	r.Status = false
	r.Message = message
}

func (r *BaseResponse) InsertData(data interface{}, message string) {
	r.Status = true
	r.Message = message
	if data != nil {
		r.Data = data
	}
}

func (r *BaseResponse) WriteResponse(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(r)
}
