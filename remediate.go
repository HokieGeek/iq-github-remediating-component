package main

import (
	"bytes"
	"fmt"
	"log"

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
	case "nuget":
		return fmt.Sprintf("pkg:nuget/%s@%s", c.name, c.version)
	case "pypi":
		return fmt.Sprintf("pkg:pypi/%s@%s?extension=%s", c.name, c.version, "tar.gz")
	case "maven":
		return fmt.Sprintf("pkg:maven/%s/%s@%s?type=%s", c.group, c.name, c.version, "jar")
	case "golang":
		return fmt.Sprintf("pkg:golang/%s@%s", c.name, c.version)
	case "ruby":
		return fmt.Sprintf("pkg:gem/%s@%s?platform=ruby", c.name, c.version)
	default:
		return ""
	}
}

func addRemediationsToPullRequest(token string, pull GithubPullRequest, remediations map[githubPullRequestFile]map[int64]component) error {
	comment := func(c component) string {
		var buf bytes.Buffer

		var href string
		switch c.format {
		case "npm":
			href = fmt.Sprintf("https://www.npmjs.com/package/%s/v/%s", c.name, c.version)
		case "maven":
			href = fmt.Sprintf("https://search.maven.org/artifact/%s/%s/%s/jar", c.group, c.name, c.version)
		case "nuget":
			href = fmt.Sprintf("https://www.nuget.org/packages/%s/%s", c.name, c.version)
		case "pypi":
			href = fmt.Sprintf("https://pypi.org/project/%s/%s", c.name, c.version)
		case "golang":
			href = fmt.Sprintf("https://%s/%s/releases/tag/%s", c.group, c.name, c.version)
		case "ruby":
			fallthrough
		case "gem":
			href = fmt.Sprintf("https://rubygems.org/gems/%s/versions/%s", c.name, c.version)
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
			err := addPullRequestComment(token, pull, pos, m.Filename, comment(comp))
			if err != nil {
				log.Printf("WARN: could not add comment: %s", err)
			}
		}
	}

	return nil
}

// ProcessPullRequestForRemediations will take a Github pull request and add any remediations if a manifest is found
func ProcessPullRequestForRemediations(iqURL, iqUser, iqPassword, iqApp, token string, pull GithubPullRequest) error {
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

	if err = addRemediationsToPullRequest(token, pull, remediations); err != nil {
		return fmt.Errorf("could not add PR comments: %v", err)
	}

	return nil
}
