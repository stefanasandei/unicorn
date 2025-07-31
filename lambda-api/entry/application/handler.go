package application

import (
	"encoding/json"
	"fmt"
	"net/http"

	common "github.com/common/model"
)

func (a *App) ExecuteRequest(w http.ResponseWriter, r *http.Request) {
	// prepare for event streaming
	err := SetupStreamingResponse(w)
	if FailIfError(err, w, "Failed to setup the response") {
		return
	}

	// read the body JSON
	var execReq common.ExecutionRequest
	err = json.NewDecoder(r.Body).Decode(&execReq)
	if FailIfError(err, w, "Failed to decode the request body") {
		return
	}

	execRes, err := a.ExecuteCode(execReq)
	if FailIfError(err, w, "Failed to execute code") {
		return
	}

	responseTask, err := json.Marshal(execRes)
	if FailIfError(err, w, "Failed to write the response") {
		return
	}

	// flush the stream to send immediate data
	_, err = fmt.Fprintf(w, string(responseTask))
	if FailIfError(err, w, "Failed to send the request") {
		return
	}

	w.(http.Flusher).Flush()
}

func (a *App) TestRequest(w http.ResponseWriter, r *http.Request) {
	// prepare for event streaming
	err := SetupStreamingResponse(w)
	if FailIfError(err, w, "Failed to setup the response") {
		return
	}

	//// read the body JSON
	//var testReq common.TestRequest
	//err = json.NewDecoder(r.Body).Decode(&testReq)
	//if FailIfError(err, w, "Failed to decode the request body") {
	//	return
	//}
	//
	//testRes, err := a.TestCode(testReq)
	//if FailIfError(err, w, "Failed to execute code") {
	//	return
	//}
	//
	//responseTask, err := json.Marshal(testRes)
	//if FailIfError(err, w, "Failed to write the response") {
	//	return
	//}

	// flush the stream to send immediate data
	_, err = fmt.Fprintf(w, string("test response"))
	if FailIfError(err, w, "Failed to send the request") {
		return
	}

	w.(http.Flusher).Flush()
}

/*
Known issues:
	- don't copy large objects at every request
	- find a better metric than cpu load to balance workload among workers
*/
