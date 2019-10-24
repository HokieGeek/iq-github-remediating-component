package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
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

func getGitHubPullRequest(req events.APIGatewayProxyRequest) (event githubPullRequest, resp *events.APIGatewayProxyResponse) {
	eventType, err := getGitHubEventType(req.Headers)
	if err != nil {
		return event, requestError(http.StatusBadRequest, fmt.Sprintf("could not parse request headers: %v", err))
	}

	// x-github-event: pull_request
	// x-gitHub-event: ping
	switch {
	case eventType == "ping":
		return event, &events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
	case eventType != "pull_request":
		log.Println("error: did not receive a supported github event")
		return event, requestError(http.StatusBadRequest, "did not receive a supported github event")
	}

	// if event.Action != "opened" {
	// 	return event, &events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
	// }

	decoded, err := url.QueryUnescape(req.Body)
	if err != nil {
		log.Println(err)
		return event, requestError(http.StatusBadRequest, fmt.Sprintf("error during url unescape of payload: %v", err))
	}
	re := regexp.MustCompile(`payload=({.*})(&.*)?$`)
	payload := re.FindAllStringSubmatch(decoded, -1)[0]
	// TODO: what if bad payload

	if err = json.Unmarshal([]byte(payload[1]), &event); err != nil {
		return event, requestError(http.StatusBadRequest, fmt.Sprintf("could not unmarshal payload as json: %v\nPAYLOAD>>%s\nDECODED>>%s", err, payload[1], decoded))
	}

	return event, nil
}

func addRemediationsToPR(token string, event githubPullRequest, remediations map[githubPullRequestFile]map[int64]component) *events.APIGatewayProxyResponse {

	for m, components := range remediations {
		for pos, comp := range components {
			err := addPullRequestComment(token, event, pos, m.Filename, comp.purl())
			if err != nil {
				log.Printf("WARN: could not add comment: %s", err)
			}
		}
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

	log.Printf("TRACE: Received Pull Request from: %s\n", event.Repository.HTMLURL)
	// log.Printf("DEBUG: %s\n", req.Body)

	token := req.QueryStringParameters["token"]

	files, err := getPullRequestFiles(token, event)
	if err != nil {
		log.Printf("ERROR: could not get files from pull request: %v\n", err)
		return *requestError(http.StatusInternalServerError, fmt.Sprintf("could not get files from pull request: %v\n", err)), nil
	}
	log.Printf("TRACE: Got %d files from full request\n", len(files))

	manifests, err := findComponentsFromManifest(files)
	if err != nil {
		log.Printf("ERROR: could not read files to find manifest: %v\n", err)
		return *requestError(http.StatusInternalServerError, fmt.Sprintf("could not read files to find manifest: %v\n", err)), nil
	}
	log.Printf("TRACE: Found manifests and added components: %q\n", manifests)

	iqURL := req.QueryStringParameters["iq_server"]
	iqAuth := strings.Split(req.QueryStringParameters["iq_auth"], ":")
	iq, err := nexusiq.New(iqURL, iqAuth[0], iqAuth[1])
	if err != nil {
		log.Printf("ERROR: could not create IQ client: %v\n", err)
		return *requestError(http.StatusInternalServerError, fmt.Sprintf("could not create IQ client: %v\n", err)), nil
	}

	iqApp := req.QueryStringParameters["iq_app"]
	remediations, err := evaluateComponents(iq, iqApp, manifests)
	if err != nil {
		log.Printf("ERROR: could not evaluate components: %v\n", err)
		return *requestError(http.StatusInternalServerError, fmt.Sprintf("could not evaluate components: %v\n", err)), nil
	}
	log.Printf("TRACE: retrieved %d remediations based on IQ app %s\n", len(remediations), iqApp)

	return *addRemediationsToPR(token, event, remediations), nil
}

func main() {
	lambda.Start(handleLambdaEvent)
}
