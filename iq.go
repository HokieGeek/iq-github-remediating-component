package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/package-url/packageurl-go"

	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
)

func evaluateComponents(iq nexusiq.IQ, nexusApplication string, manifests map[githubPullRequestFile]map[int64]component) (map[githubPullRequestFile]map[int64]component, error) {
	asIQComponent := func(c component) (nexusiq.Component, error) {
		// TODO: how bout errors and validation?
		return nexusiq.Component{PackageURL: c.purl()}, nil
	}

	asComponent := func(c nexusiq.Component) (component, error) {
		log.Printf("TRACE: asComponent(): %q\n", c)

		switch {
		case c.ComponentID != nil:
			return component{
				c.ComponentID.Format,
				c.ComponentID.Coordinates.GroupID,
				c.ComponentID.Coordinates.ArtifactID,
				c.ComponentID.Coordinates.Version,
			}, nil

		case c.PackageURL != "":
			purl, err := packageurl.FromString(c.PackageURL)
			if err != nil {
				return component{}, fmt.Errorf("could not parse PackageURL: %v", err)
			}
			return component{
				purl.Type,
				purl.Namespace,
				purl.Name,
				purl.Version,
			}, nil
		}

		return component{}, errors.New("nexusiq.Component not formatted well enough to parse")
	}

	remediations := make(map[githubPullRequestFile]map[int64]component)

	for m, components := range manifests {
		log.Printf("TRACE: evaluating manifest: %s\n", m.Filename)
		remediated := make(map[int64]component)
		for p, c := range components {
			component, _ := asIQComponent(c)
			log.Printf("TRACE: evaluating component for manifest %s: %q\n", m.Filename, component)
			remediation, err := nexusiq.GetRemediationByApp(iq, component, nexusiq.StageBuild, nexusApplication)
			if err != nil {
				log.Printf("ERROR: could not evaluate component %s: %s\n", component, err)
				continue
			}

			rcomp, err := remediation.ComponentForRemediationType(nexusiq.RemediationTypeNoViolations)
			if err != nil {
				log.Printf("WARN: did not find remediating component for %s: %s\n", component, err)
				log.Printf("TRACE: remediation: %q\n", remediation)
				continue
			}
			comp, _ := asComponent(rcomp)
			remediated[p] = comp
		}
		if len(remediated) > 0 {
			remediations[m] = remediated
		}
	}

	return remediations, nil
}
