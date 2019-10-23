package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/sonatype-nexus-community/gonexus"
)

type MyEvent struct {
	Ping string `json:"ping"`
}

type MyResponse struct {
	Pong string `json:"pong"`
}

func HandleLambdaEvent(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var event MyEvent
	err := json.Unmarshal([]byte(req.Body), &event)

	buf, err := json.Marshal(MyResponse{Pong: fmt.Sprintf("Pong!: %s", event.Ping)})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(buf),
	}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
