package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func processRequestResponse(requestBody []byte) (requestID string, responseBody []byte) {
	log.Printf("got response: %s", requestBody)

	return "", nil
}

func main() {

	netClient := &http.Client{
		Timeout: time.Second * 5,
	}

	runtimeAPIBase := fmt.Sprintf("http://%s/2018-06-01/runtime", os.Getenv("AWS_LAMBDA_RUNTIME_API"))
	inocationNextURL := fmt.Sprintf("%s/invocation/next", runtimeAPIBase)

	for {
		resp, _ := netClient.Get(inocationNextURL)
		invocationReq, _ := ioutil.ReadAll(resp.Body)

		requestID, invocationResp := processRequestResponse(invocationReq)

		invocationResponseURL := fmt.Sprintf("%s/invocation/%s/response", runtimeAPIBase, requestID)
		req, err := http.NewRequest("POST", invocationResponseURL, bytes.NewReader(invocationResp))
		if err != nil {
			log.Fatalf("error: %s", err.Error())
		}
		netClient.Do(req)
	}
}
