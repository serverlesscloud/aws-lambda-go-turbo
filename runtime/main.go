package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

func processRequestResponse(requestBody []byte) (responseBody []byte) {
	log.Printf("got response: %s", requestBody)

	return nil
}

func main() {

	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	runtimeAPIBase := fmt.Sprintf("http://%s/2018-06-01/runtime", os.Getenv("AWS_LAMBDA_RUNTIME_API"))
	inocationNextURL := fmt.Sprintf("%s/invocation/next", runtimeAPIBase)

	for {
		resp, err := netClient.Get(inocationNextURL)
		if err != nil {
			log.Fatalf("error getting next invocation: %s", err.Error())
		}

		invocationReq, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatalf("error reading invocation body: %s", err.Error())
		}

		lc := &lambdacontext.LambdaContext{
			AwsRequestID:       resp.Header.Get("Lambda-Runtime-Aws-Request-Id"),
			InvokedFunctionArn: resp.Header.Get("Lambda-Runtime-Invoked-Function-Arn"),
			Identity: lambdacontext.CognitoIdentity{
				CognitoIdentityID: resp.Header.Get("Lambda-Runtime-Cognito-Identity"),
			},
		}

		invocationResp := processRequestResponse(invocationReq)

		invocationResponseURL := fmt.Sprintf("%s/invocation/%s/response", runtimeAPIBase, lc.AwsRequestID)
		req, err := http.NewRequest("POST", invocationResponseURL, bytes.NewReader(invocationResp))
		if err != nil {
			log.Fatalf("error: %s", err.Error())
		}
		netClient.Do(req)
	}
}
