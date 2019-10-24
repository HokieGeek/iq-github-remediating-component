package main

import (
	"log"

	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
)

func evaluateComponents(iq nexusiq.IQ, nexusApplication string, manifests map[githubPullRequestFile]map[int64]component) (map[githubPullRequestFile]map[int64]component, error) {
	asIQComponent := func(c component) (nexusiq.Component, error) {
		// TODO: how bout errors and validation?
		return nexusiq.Component{PackageURL: c.purl()}, nil
	}

	asComponent := func(c nexusiq.Component) (component, error) {
		return component{
			c.ComponentID.Format,
			c.ComponentID.Coordinates.GroupID,
			c.ComponentID.Coordinates.ArtifactID,
			c.ComponentID.Coordinates.Version,
		}, nil
	}

	remediations := make(map[githubPullRequestFile]map[int64]component)

	for m, components := range manifests {
		remediated := make(map[int64]component)
		for p, c := range components {
			component, _ := asIQComponent(c)
			remediation, err := nexusiq.GetRemediationByApp(iq, component, nexusiq.StageBuild, nexusApplication)
			if err != nil {
				// TODO: handle
				continue
			}

			log.Printf("TRACE: evaluating component for manifest %s: %q\n", m.Filename, component)

			rcomp, _ := remediation.ComponentForRemediationType(nexusiq.RemediationTypeNoViolations)
			comp, _ := asComponent(rcomp)
			remediated[p] = comp
		}
		if len(remediated) > 0 {
			remediations[m] = remediated
		}
	}

	return remediations, nil
}
