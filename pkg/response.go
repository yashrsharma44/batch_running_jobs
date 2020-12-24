package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WorkerStatus struct {
	ID           string `json:"worker-id"`
	Status       string `json:"task-status"`
	ResponseCode int    `json:"reponse-status-code"`
	Message      string `json:"message"`
}

func NewResponse(ID string, Status string, ResponseCode int, Message string, w *http.ResponseWriter) {
	(*w).WriteHeader(ResponseCode)

	respObj := &WorkerStatus{
		ID:           ID,
		Status:       Status,
		ResponseCode: ResponseCode,
		Message:      Message,
	}

	if err := json.NewEncoder((*w)).Encode(respObj); err != nil {
		fmt.Fprintf((*w), "Server Error. err=%v", err)
	}
}
