package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
)

type component struct {
	format, group, name, version string
}

func (c component) purl() string {
	switch c.format {
	case "npm":
		return fmt.Sprintf("pkg:npm/%s@%s", c.name, c.version)
	// case "maven":
	// 	return fmt.Sprintf("pkg:maven/%s/%s@%s?type=%s", "group", c.name, c.version, "type")
	// case "nuget":
	// 	return fmt.Sprintf("pkg:nuget/%s@%s", c.name, c.version)
	// case "golang":
	// 	return fmt.Sprintf("pkg:golang/%s@%s", c.name, c.version)
	// case "pypi":
	// 	return fmt.Sprintf("pkg:pypi/%s@%s?extension=%s", c.name, c.version, "ext")
	// case "ruby":
	// 	return fmt.Sprintf("pkg:gem/%s@%s?platform=ruby", c.name, c.version)
	default:
		return ""
	}
}

func requestError(statusCode int, message string) *events.APIGatewayProxyResponse {
	body := message
	if body == "" {
		body = http.StatusText(statusCode)
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
	}
}

func isSupportedEventType(req events.APIGatewayProxyRequest) (bool, int) {
	eventType, err := getGitHubEventType(req.Headers)
	if err != nil {
		log.Printf("could not parse request headers: %v", err)
		return false, http.StatusBadRequest
	}

	// x-github-event: pull_request
	// x-gitHub-event: ping
	switch {
	case eventType == "ping":
		return false, http.StatusOK
	case eventType != "pull_request":
		log.Println("ERROR: did not receive a supported github event")
		return false, http.StatusBadRequest
	}

	return true, http.StatusOK
}

func addCommentsToPR(token string, event githubPullRequest, remediations map[githubPullRequestFile]map[int64]component) error {
	for m, components := range remediations {
		for pos, comp := range components {
			err := addPullRequestComment(token, event, pos, m.Filename, comp.purl())
			if err != nil {
				log.Printf("WARN: could not add comment: %s", err)
			}
		}
	}

	return nil
}

func evaluatePullRequest(iq nexusiq.IQ, token, iqApp string, event githubPullRequest) error {
	log.Printf("TRACE: Received Pull Request from: %s\n", event.Repository.HTMLURL)
	// log.Printf("DEBUG: %s\n", req.Body)

	files, err := getPullRequestFiles(token, event)
	if err != nil {
		log.Printf("ERROR: could not get files from pull request: %v\n", err)
		return fmt.Errorf("could not get files from pull request: %v", err)
	}
	log.Printf("TRACE: Got %d files from pull request\n", len(files))

	manifests, err := findComponentsFromManifest(files)
	if err != nil {
		log.Printf("ERROR: could not read files to find manifest: %v\n", err)
		return fmt.Errorf("could not read files to find manifest: %v", err)
	}
	log.Printf("TRACE: Found manifests and added components: %q\n", manifests)

	remediations, err := evaluateComponents(iq, iqApp, manifests)
	if err != nil {
		log.Printf("ERROR: could not evaluate components: %v\n", err)
		return fmt.Errorf("could not evaluate components: %v", err)
	}
	log.Printf("TRACE: retrieved %d remediations based on IQ app %s\n", len(remediations), iqApp)

	if err = addCommentsToPR(token, event, remediations); err != nil {
		return fmt.Errorf("could not add PR comments: %v", err)
	}

	return nil
}

func handlePullRequest(iqURL, iqUser, iqPassword, iqApp, githubToken string, pull githubPullRequest) error {
	iq, err := nexusiq.New(iqURL, iqUser, iqPassword)
	if err != nil {
		log.Printf("ERROR: could not create IQ client: %v", err)
		return fmt.Errorf("could not create IQ client: %v", err)
	}
	log.Printf("TRACE: created client to IQ server as: %s\n", iqApp)

	err = evaluatePullRequest(iq, githubToken, iqApp, pull)
	if err != nil {
		log.Printf("ERROR: could not parse pull request: %v", err)
		return fmt.Errorf("could not parse pull request: %v", err)
	}
	log.Println("TRACE: done")

	return nil
}

func handleLambdaEvent(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if supported, status := isSupportedEventType(req); !supported {
		return *requestError(status, "Unsupported event"), nil
	}

	var event githubPullRequest
	if err := json.Unmarshal([]byte(req.Body), &event); err != nil {
		return *requestError(http.StatusBadRequest, fmt.Sprintf("could not unmarshal payload as json: %v", err)), err
	}

	if event.Action != "opened" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNoContent,
			Body:       "Only processing new pull requests",
		}, nil
	}

	log.Printf("TRACE: Received Pull Request from: %s\n", event.Repository.HTMLURL)
	// log.Printf("DEBUG: %s\n", req.Body)

	token := req.QueryStringParameters["token"]
	iqApp := req.QueryStringParameters["iq_app"]
	iqURL := req.QueryStringParameters["iq_url"]
	iqAuth := strings.Split(req.QueryStringParameters["iq_auth"], ":")
	log.Println("TRACE: handling request")
	if err := handlePullRequest(iqURL, iqAuth[0], iqAuth[1], iqApp, token, event); err != nil {
		panic(fmt.Sprintf("ERROR: error handling pull request: %v\n", err))
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Evaluating new pull request",
	}, nil
}

func testPR() {
	tokenf, err := os.Open("token")
	if err != nil {
		panic(err)
	}
	defer tokenf.Close()

	buf, _ := ioutil.ReadAll(tokenf)
	token := string(buf)

	pr, err := os.Open("pullrequest_sample.json")
	if err != nil {
		panic(err)
	}
	defer pr.Close()
	buf, _ = ioutil.ReadAll(pr)
	var pull githubPullRequest
	if err := json.Unmarshal(buf, &pull); err != nil {
		fmt.Println(string(buf))
		panic(err)
	}

	url := os.Args[1]
	user := os.Args[2]
	pass := os.Args[3]
	if err := handlePullRequest(url, user, pass, "jshop", token, pull); err != nil {
		panic(fmt.Sprintf("ERROR: error handling pull request: %v\n", err))
	}
}

func main() {
	lambda.Start(handleLambdaEvent)
	// testPR()
}
