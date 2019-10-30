package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/package-url/packageurl-go"

	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
)

func httpreq(method, url string, payload io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth("admin", "admin1234")

	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return client.Do(request)
}

func getComponentRemediations(iq nexusiq.IQ, nexusApplication string, manifests componentRemediations) (componentRemediations, error) {
	asIQComponent := func(c component) (nexusiq.Component, error) {
		// TODO: how bout errors and validation?
		return nexusiq.Component{PackageURL: c.purl()}, nil
	}

	asComponent := func(c nexusiq.Component) (component, error) {
		log.Printf("TRACE: asComponent(): %#v\n", c)

		switch {
		case c.PackageURL != "":
			purl, err := packageurl.FromString(c.PackageURL)
			log.Printf("TRACE: PURL: %#v\n", purl)
			if err != nil {
				return component{}, fmt.Errorf("could not parse PackageURL: %v", err)
			}
			return component{
				purl.Type,
				purl.Namespace,
				purl.Name,
				purl.Version,
			}, nil

		case c.ComponentID != nil:
			log.Printf("TRACE: CID: %#v\n", c.ComponentID)
			log.Printf("TRACE: C.PURL: %s\n", c.PackageURL)
			return component{
				c.ComponentID.Format,
				c.ComponentID.Coordinates.GroupID,
				c.ComponentID.Coordinates.ArtifactID,
				c.ComponentID.Coordinates.Version,
			}, nil
		}

		return component{}, errors.New("nexusiq.Component not formatted well enough to parse")
	}

	remediations := make(componentRemediations)

	for m, components := range manifests {
		log.Printf("TRACE: evaluating manifest: %s\n", m.Filename)
		remediated := make(map[changeLocation]component)
		log.Printf("TRACE: manifest components: %v\n", components)
		for loc, c := range components {
			iqcomponent, _ := asIQComponent(c)
			log.Printf("TRACE: evaluating %s component for manifest %s: %v\n", nexusApplication, m.Filename, iqcomponent)

			log.Println("TRACE: retrieving remediating component")
			remediation, err := nexusiq.GetRemediationByApp(iq, iqcomponent, nexusiq.StageBuild, nexusApplication)
			if err != nil {
				log.Printf("ERROR: could not evaluate component %v: %v\n", iqcomponent, err)
				continue
			}

			rcomp, err := remediation.ComponentForRemediationType(nexusiq.RemediationTypeNoViolations)
			if err != nil {
				log.Printf("WARN: did not find remediating component for %v: %v\n", iqcomponent, err)
				log.Printf("TRACE: remediation: %v\n", remediation)
				continue
			}

			comp, err := asComponent(rcomp)
			if err != nil {
				log.Printf("ERROR: could not parse remediating component object %v: %v\n", rcomp, err)
				log.Printf("TRACE: remediation: %v\n", remediation)
				continue
			}

			// Only add if it's a different version
			if comp.version == c.version {
				continue
			}

			// TODO: evaluate the component to determine what is wrong with it

			log.Printf("TRACE: adding suggestion: %v[%#v] = %v\n", iqcomponent, loc, comp)
			remediated[loc] = comp
		}

		if len(remediated) > 0 {
			remediations[m] = remediated
		}
	}

	return remediations, nil
}
