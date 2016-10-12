package middleware

import (
	"fmt"
	"net/http"
	"encoding/json"
)

var Result HandlerResult

type HandlerResult struct {
	StatusCode int
	Result     interface{}
	Error      *error
}

type TextResponse struct {
	Result	string
}

type ResultHandler func(http.ResponseWriter, *http.Request) HandlerResult

func (inner ResultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Result = inner(w, r)
	switch Result.Result.(type) {
	case string:
		if Result.Error != nil {
			http.Error(w, Result.Result.(string), Result.StatusCode)
			return
		}
	default:
		json, err := json.Marshal(Result.Result)
		if err != nil {
			http.Error(w, "Error encoding object", http.StatusInternalServerError)
			return
		}
		Result.Result = string(json[:])
	}
	fmt.Fprint(w, Result.Result.(string))
}
