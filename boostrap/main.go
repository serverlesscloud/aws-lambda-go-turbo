package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

)

const RuntimeAPIBase := fmt.Sprintf("http://%s/2018-06-01/runtime", os.Getenv("AWS_LAMBDA_RUNTIME_API"))

func processRequestResponse(requestBody []byte) requestID string, responseBody []byte {
	log.Printf("got response: %s", body)

	return "", nil
}

func main() {

	netClient := &http.Client{
		Timeout: time.Second * 5,
	}

	inocationNextURL := fmt.Sprintf("%s/invocation/next", RuntimeAPIBase)

	for {
		resp, _ := netClient.Get(inocationNextURL)
		invocationReq, _ := ioutil.ReadAll(resp.Body)

		requestID, invocationResp := processRequestResponse(invocationReq)

		invocationResponseURL := fmt.Sprintf("%s/invocation/%s/response", RuntimeAPIBase, requestID)
		req, err := http.NewRequest("POST", invocationResponseURL, req)
		netClient.Do(req)
	}
}
