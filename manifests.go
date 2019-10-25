package main

import (
	"bufio"
	"regexp"
	"strings"
)

func componentsFromNpm(lines map[int64]string) (map[int64]component, error) {
	converted := make(map[int64]component)
	re := regexp.MustCompile(`"([^"]*)": ".?([0-9](\.[0-9])+)",?`)

	for p, l := range lines {
		found := re.FindAllStringSubmatch(l, -1)[0]
		name := found[1]
		version := found[2]
		converted[p] = component{format: "npm", name: name, version: version}
	}

	return converted, nil
}

func parsePatchAdditions(patch string) map[int64]string {
	adds := make(map[int64]string)

	scanner := bufio.NewScanner(strings.NewReader(patch))
	var position int64
	for scanner.Scan() {
		if scanner.Text()[:2] == "+ " {
			adds[position] = scanner.Text()[2:]
		}
		position++
	}

	return adds
}

func findComponentsFromManifest(files []githubPullRequestFile) (map[githubPullRequestFile]map[int64]component, error) {
	getComponents := func(patch string, linesToComponents func(lines map[int64]string) (map[int64]component, error)) (map[int64]component, error) {
		additions := parsePatchAdditions(patch)
		return linesToComponents(additions)
	}

	manifests := make(map[githubPullRequestFile]map[int64]component, 0)

	for _, f := range files {
		components := make(map[int64]component)
		var err error
		switch f.Filename {
		// case "go.mod":
		// case "build.gradle":
		// case "pom.xml":
		case "package.json":
			components, err = getComponents(f.Patch, componentsFromNpm)
		}

		if err != nil {
			// TODO
			continue
		}

		manifests[f] = components
	}

	return manifests, nil
}
