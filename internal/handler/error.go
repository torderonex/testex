package handler

import "encoding/json"

type Error struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func newError(message string, statusCode int) Error {
	return Error{Message: message, StatusCode: statusCode}
}

func (e Error) ToJson() string {
	data, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(data)
}
