package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	return `
	{
		"user_id": 1
	}
`, nil

}

func main() {
	lambda.Start(HandleRequest)
}
