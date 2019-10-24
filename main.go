package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
)

func requestError(statusCode int) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       http.StatusText(statusCode),
	}
}

func getGitHubPullRequest(req events.APIGatewayProxyRequest) (event githubPullRequest, resp *events.APIGatewayProxyResponse) {
	eventType, err := getGitHubEventType(req.Headers)
	if err != nil {
		return event, requestError(http.StatusBadRequest)
	}

	// x-github-event: pull_request
	// x-gitHub-event: ping
	switch {
	case eventType == "ping":
		return event, &events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
	case eventType != "pull_request":
		log.Println("error: did not receive a supported github event")
		return event, requestError(http.StatusBadRequest)
	}

	decoded, err := url.QueryUnescape(req.Body)
	if err != nil {
		log.Println(err)
		return event, requestError(http.StatusBadRequest)
	}
	payload := strings.TrimPrefix(decoded, "payload=")

	if err = json.Unmarshal([]byte(payload), &event); err != nil {
		return event, requestError(http.StatusInternalServerError)
	}

	return event, nil
}

func findComponentsFromManifest(files []githubPullRequestFile) ([]string, error) {
	manifests := make([]githubPullRequestFile, 0)

	for _, f := range files {
		switch f.Filename {
		case "package.json":
			manifests = append(manifests, f)
		}
	}

	log.Println("DEBUG: Changed manifests")

	for _, m := range manifests {
		log.Printf("DEBUG: %s: %s\n", m.Filename, m.Patch)
	}

	return nil, nil
}

func addRemediationsToPR(token string, event githubPullRequest, remediations map[string]string) *events.APIGatewayProxyResponse {

	err := addPullRequestComment(event, token, "THINGY")
	if err != nil {
		return nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		// Body:       string(buf),
	}
}

func handleLambdaEvent(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	event, resp := getGitHubPullRequest(req)
	if resp != nil {
		return *resp, nil
	}

	log.Printf("Received Pull Request from: %s\n", event.Repository.HTMLURL)
	log.Printf("DEBUG: %s\n", req.Body)

	token := req.QueryStringParameters["token"]

	files, err := getPullRequestFiles(event, token)
	if err != nil {
		log.Printf("ERROR: could not get files from pull request: %v\n", err)
		return *requestError(http.StatusInternalServerError), nil
	}

	components, err := findComponentsFromManifest(files)
	if err != nil {
		log.Printf("ERROR: could not read files to find manifest: %v\n", err)
		return *requestError(http.StatusInternalServerError), nil
	}

	iqURL := req.QueryStringParameters["iq_server"]
	iqAuth := strings.Split(req.QueryStringParameters["iq_auth"], ":")
	iq, err := nexusiq.New(iqURL, iqAuth[0], iqAuth[1])
	if err != nil {
		log.Printf("ERROR: could not create IQ client: %v\n", err)
		return *requestError(http.StatusInternalServerError), nil
	}

	iqApp := req.QueryStringParameters["iq_app"]
	remediations, err := evaluateComponents(iq, iqApp, components)
	if err != nil {
		log.Printf("ERROR: could not evaluate components: %v\n", err)
		return *requestError(http.StatusInternalServerError), nil
	}

	return *addRemediationsToPR(token, event, remediations), nil
}

func main() {
	lambda.Start(handleLambdaEvent)
}
