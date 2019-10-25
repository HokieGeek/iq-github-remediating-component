package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
)

type component struct {
	format, group, name, version string
}

func (c component) purl() string {
	/*
		purl := packageurl.NewPackageURL(c.format, c.group, c.name, c.version, nil, nil)
		if purl == nil {
			return errors.New("could not create PackageURL string")
		}
		return purl.ToString()
	*/
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

func isSupportedEventType(req events.APIGatewayProxyRequest) (bool, int) {
	eventType, err := getGitHubEventType(req.Headers)
	if err != nil {
		log.Printf("could not parse request headers: %v", err)
		return false, http.StatusBadRequest
	}

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
	comment := func(c component) string {
		var buf bytes.Buffer

		var href string
		switch c.format {
		case "npm":
			href = fmt.Sprintf("https://www.npmjs.com/package/%s/v/%s", c.name, c.version)
		case "maven":
			href = fmt.Sprintf("https://search.maven.org/artifact/%s/%s/%s/jar", c.group, c.name, c.version)
		}

		buf.WriteString("[Nexus Lifecycle](https://www.sonatype.com/product-nexus-lifecycle) has found that this version of `")
		buf.WriteString(c.name)
		buf.WriteString("` violates your company's policies.\n\n")
		buf.WriteString("Lifecycle recommends using version [")
		buf.WriteString(c.version)
		buf.WriteString("](")
		buf.WriteString(href)
		buf.WriteString(") instead as it does not violate any policies.\n\n")

		return buf.String()
	}

	for m, components := range remediations {
		for pos, comp := range components {
			err := addPullRequestComment(token, event, pos, m.Filename, comment(comp))
			if err != nil {
				log.Printf("WARN: could not add comment: %s", err)
			}
		}
	}

	return nil
}

func handlePullRequest(iqURL, iqUser, iqPassword, iqApp, token string, pull githubPullRequest) error {
	log.Printf("TRACE: Received Pull Request from: %s\n", pull.Repository.HTMLURL)

	iq, err := nexusiq.New(iqURL, iqUser, iqPassword)
	if err != nil {
		log.Printf("ERROR: could not create IQ client: %v", err)
		return fmt.Errorf("could not create IQ client: %v", err)
	}
	log.Printf("TRACE: created client to IQ server as: %s\n", iqApp)

	files, err := getPullRequestFiles(token, pull)
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

	remediations, err := getComponentRemediations(iq, iqApp, manifests)
	if err != nil {
		log.Printf("ERROR: could not find remediation version for components: %v\n", err)
		return fmt.Errorf("could not find remediation version for components: %v", err)
	}
	log.Printf("TRACE: retrieved %d remediations based on IQ app %s\n", len(remediations), iqApp)

	if err = addCommentsToPR(token, pull, remediations); err != nil {
		return fmt.Errorf("could not add PR comments: %v", err)
	}

	return nil
}

func handleLambdaEvent(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requestResponse := func(statusCode int, message string) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{
			StatusCode: statusCode,
			Body:       message,
		}
	}

	if supported, status := isSupportedEventType(req); !supported {
		return requestResponse(status, "Unsupported event"), nil
	}

	var event githubPullRequest
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
	if err := handlePullRequest(iqURL, iqAuth[0], iqAuth[1], iqApp, token, event); err != nil {
		return requestResponse(http.StatusInternalServerError, fmt.Sprintf("ERROR: error handling pull request: %v\n", err)), err
	}

	return requestResponse(http.StatusOK, "Evaluating new pull request"), nil
}

func main() {
	lambda.Start(handleLambdaEvent)
}
