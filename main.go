package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
)

func handleLambdaEvent(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requestResponse := func(statusCode int, message string) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{
			StatusCode: statusCode,
			Body:       message,
		}
	}

	token := req.QueryStringParameters["token"]
	iqApp := req.QueryStringParameters["iq_app"]
	iqURL := req.QueryStringParameters["iq_url"]
	iqAuth := strings.Split(req.QueryStringParameters["iq_auth"], ":")

	iq, err := nexusiq.New(iqURL, iqAuth[0], iqAuth[1])
	if err != nil {
		err2 := fmt.Errorf("could not create IQ client: %v", err)
		log.Printf("ERROR: %v", err2)
		return requestResponse(http.StatusInternalServerError, err2.Error()), err2
	}
	log.Printf("TRACE: created client to IQ server as: %s\n", iqApp)

	// Github webhook comes in two parts.
	// One is a ping to verify the connection
	// The other is the actual event
	supported, status := IsValidGithubWebhookPullRequestEvent(req.Headers)
	switch {
	case !supported && status == http.StatusOK:
		return requestResponse(status, "Acknowledging ping"), nil
	case !supported && status != http.StatusOK:
		log.Println("WARN: Did not receive a valid Github webhook")
		// We don't return here in case what we got was a Gitlab webhook
	case supported:
		status, err := HandleGithubWebhookPullRequestEvent(iq, iqApp, token, []byte(req.Body))
		if err != nil {
			log.Printf("ERROR: %v", err)
			return requestResponse(status, err.Error()), err
		}
		return requestResponse(status, "Evaluating new Github pull request"), nil
	}

	// Gitlab's webhook event is simpler. Just need a quick binary check
	supported, status = IsValidGitlabWebhookMergeRequestEvent(req.Headers)
	if !supported {
		log.Println("WARN: Did not receive a valid Gitlab webhook")
		return requestResponse(status, "Did not receive a valid Gitlab webhook"), nil
	}

	status, err = HandleGitlabWebhookMergeRequestEvent(iq, iqApp, token, []byte(req.Body))
	if err != nil {
		log.Printf("ERROR: %v", err)
		return requestResponse(status, err.Error()), err
	}
	return requestResponse(status, "Evaluating new Gitlab merge request"), nil
}

func main() {
	lambda.Start(handleLambdaEvent)
}
