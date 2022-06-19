package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ResponseType struct {
	W http.ResponseWriter
	R *http.Request
}

func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}

func ErrorResponseHandler(output string, httpCode int, err error, respType ResponseType) {
	type Result struct {
		Status     int    `json:"status"`
		StatusText string `json:"status_text"`
	}
	var level string
	switch httpCode {

	case http.StatusBadRequest, http.StatusOK:
		level = "INFO"

	case http.StatusInternalServerError:
		level = "ERRR"

	case http.StatusUnauthorized:
		level = "WARN"
	}

	if err != nil {
		//	received error
		respType.W.WriteHeader(httpCode)
		respType.W.Header().Set("output", output)
		jErr := json.NewEncoder(respType.W).Encode(Result{Status: httpCode, StatusText: output})
		if jErr != nil { //Encode Error Handling
			respType.W.Header().Set("Level", level)
			respType.W.Header().Set("err", err.Error()+", and json encoding error: "+jErr.Error())
		}
		respType.W.Header().Set("Level", level)
		respType.W.Header().Set("err", err.Error())
	} else {
		//	no error
		respType.W.WriteHeader(httpCode)
		respType.W.Header().Set("output", output)
		jErr := json.NewEncoder(respType.W).Encode(Result{Status: httpCode, StatusText: output})
		if jErr != nil { //Encode Error Handling
			respType.W.Header().Set("Level", level)
			respType.W.Header().Set("err", "json encoding error: "+jErr.Error())
		}
		respType.W.Header().Set("Level", level)
		respType.W.Header().Set("err", "nil")
	}
}
