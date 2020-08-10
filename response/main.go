package response

import (
	"encoding/json"
	"net/http"
	"time"
)

type TextResponse struct {
	Message string `json:"message"`
}

type ApiResponse struct {
	Meta   interface{} `json:"meta"`
	Data   interface{} `json:"data"`
	Errors []string    `json:"errors"`
}

type MetaResponse struct {
	Version     string    `json:"version"`
	RequestedAt time.Time `json:"requestedAt"`
	Application string `json:"application"`
}

func SendResponse(w http.ResponseWriter, data interface{}, status int, errs []error) {
	var sendResponse []byte
	var errorList []string
	var noNilErrors []error

	for _, err := range errs {
		if err != nil {
			noNilErrors = append(noNilErrors, err)
		}
	}

	if len(noNilErrors) > 0 {
		var blankData interface{}
		data = blankData
		errorList = errorsToStrings(noNilErrors)

		// Force a 404 for any record not found errors
		if len(errorList) == 1 && errorList[0] == "record not found" {
			status = 404
		}
	}
	if status == 0 {
		status = 200
	}

	apiResponse := ApiResponse{getMetaData(), data, errorList}
	sendResponse, _ = json.Marshal(apiResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(sendResponse)
	return
}

func SendUnauthorized(w http.ResponseWriter) {
	var sendResponse []byte
	apiResponse := ApiResponse{getMetaData(), EmptyInterface(), []string{"You don't have permssion to this."}}
	sendResponse, _ = json.Marshal(apiResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(403)
	w.Write(sendResponse)
	return
}

func EmptyInterface() interface{} {
	var empty interface{}
	return empty
}

func getMetaData() interface{} {
	newMeta := MetaResponse{
		"1",
		time.Now(),
		"Test App",
	}
	return newMeta
}

func errorsToStrings(errors []error) []string {
	var strings []string
	for _, err := range errors {
		strings = append(strings, err.Error())
	}
	return strings
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	var errors []string
	apiResponse := ApiResponse{getMetaData(), payload, errors}
	response, _ := json.Marshal(apiResponse)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithJSONError(w http.ResponseWriter, code int, errors []string) {
	var data interface{}
	apiResponse := ApiResponse{getMetaData(), data, errors}
	response, _ := json.Marshal(apiResponse)

	// Force a 404 for any record not found errors
	if len(errors) == 1 && errors[0] == "record not found" {
		code = 404
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
