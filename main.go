package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleLambdaEvent(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requestResponse := func(statusCode int, message string) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{
			StatusCode: statusCode,
			Body:       message,
		}
	}

	if supported, status := IsValidGithubWebhookPullRequestEvent(req.Headers); !supported {
		return requestResponse(status, "Unsupported event"), nil
	}

	var event GithubPullRequest
	if err := json.Unmarshal([]byte(req.Body), &event); err != nil {
		return requestResponse(http.StatusBadRequest, fmt.Sprintf("could not unmarshal payload as json: %v", err)), err
	}

	if event.Action != "opened" {
		return requestResponse(http.StatusNoContent, "Only processing new pull requests"), nil
	}

	token := req.QueryStringParameters["token"]
	iqApp := req.QueryStringParameters["iq_app"]
	iqURL := req.QueryStringParameters["iq_url"]
	iqAuth := strings.Split(req.QueryStringParameters["iq_auth"], ":")
	if err := ProcessPullRequestForRemediations(iqURL, iqAuth[0], iqAuth[1], iqApp, token, event); err != nil {
		return requestResponse(http.StatusInternalServerError, fmt.Sprintf("ERROR: error handling pull request: %v\n", err)), err
	}

	return requestResponse(http.StatusOK, "Evaluating new pull request"), nil
}

func main() {
	lambda.Start(handleLambdaEvent)
}
