package main

import (
	"bufio"
	"regexp"
	"strings"
)

var testPatch = `@@ -75,6 +75,7 @@
],
"dependencies": {
"body-parser": "^1.19.0",
+ "chalk": "^1.0.0",
"check-dependencies": "^1.1.0",
"clarinet": "^0.12.3",
"colors": "^1.3.3",
@@ -85,7 +86,6 @@
"cors": "^2.8.5",
"dottie": "^2.0.1",
"download": "^7.1.0",
- "errorhandler": "^1.5.1",
"express": "^4.17.1",
"express-jwt": "0.1.3",
"express-rate-limit": "^4.0.1",
@@ -112,7 +112,7 @@
"libxmljs2": "^0.21.3",
"marsdb": "^0.6.11",
"morgan": "^1.9.1",
- "moment": "^2.3.0",
+ "moment": "^2.1.0",
"multer": "^1.4.1",
"node-pre-gyp": "^0.13.0",
"notevil": "^1.3.1 ",
`

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
