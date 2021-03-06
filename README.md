# AWS Lambda Go Turbo Runtime
A Lambda custom runtime for Go that eliminates the need for unecessary reflection nor having to run a RPC server.

Came about as a result of seeing the slightly dissapointing results in [Lambda Performence Benchmarks 2018](https://read.acloud.guru/comparing-aws-lambda-performance-of-node-js-python-java-c-and-go-29c1163c2581)

## Features

* [X] Custom Runtime boostrapper
* [ ] AWS SAM Support - pending https://github.com/aws/aws-cli/issues/3789
* [ ] Serverless Framework Support
* [ ] Lambda Layers
* [ ] Use Go plugins for dynamic linking

## History

I and a few keen contributors created [serverless-golang](https://github.com/yunspace/serverless-golang/) back in Jan 2017, 11 months before official [aws-lambda-go](https://github.com/aws/aws-lambda-go) support was announced during ReInvent 2017. It used a `go plugin` approach based on the awesome work from [eawsy](https://github.com/eawsy/) team. The framework was used in production and quite flexible in that all `Handers` can be written into a [single file](https://github.com/yunspace/serverless-golang/blob/master/examples/aws-golang-event/handler.go).

Since moving on to the official `aws-lambda-go`, I always felt a bit uneasy with the use of RPC and amount of reflection for what I perceive to be a quite straight forward function call. 

To address reflection, I submitted [PR #69](https://github.com/aws/aws-lambda-go/pull/69) to allow for a custom non-reflection based Handler to be passed. However there was still no way to side step RPC in the existing `aws-lambda-go` runtime.

Then came ReInvent 2018 and [lambda custom runtimes](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-custom.html). Thus this project was born.
