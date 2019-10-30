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

	if supported, _ := IsValidGithubWebhookPullRequestEvent(req.Headers); supported {
		status, err := HandleGithubWebhookPullRequestEvent(iq, iqApp, token, []byte(req.Body))
		if err != nil {
			log.Printf("ERROR: %v", err)
			return requestResponse(status, err.Error()), err
		}
		return requestResponse(status, "Evaluating new Github pull request"), nil
	}

	if supported, _ := IsValidGitlabWebhookMergeRequestEvent(req.Headers); supported {
		status, err := HandleGitlabWebhookMergeRequestEvent(iq, iqApp, token, []byte(req.Body))
		if err != nil {
			log.Printf("ERROR: %v", err)
			return requestResponse(status, err.Error()), err
		}
		return requestResponse(status, "Evaluating new Gitlab merge request"), nil
	}

	return requestResponse(http.StatusNotFound, "Did not recognize webhook event"), nil
}

func main() {
	lambda.Start(handleLambdaEvent)
}
