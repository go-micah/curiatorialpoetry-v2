package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler() {

	url := os.Getenv("AMPLIFY_WEBHOOK_URL")
	fmt.Printf("Sending to webhook: %s", url)

	body := []byte(`{}`)

	_, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	lambda.Start(handler)
}
