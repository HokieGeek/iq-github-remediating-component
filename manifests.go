package main

import (
	"bufio"
	"log"
	"regexp"
	"strings"
)

func componentsSingleLineNameVersion(lines map[int64]string, re *regexp.Regexp, format string) (map[int64]component, error) {
	converted := make(map[int64]component)

	for p, l := range lines {
		matches := re.FindAllStringSubmatch(l, -1)
		if len(matches) == 0 {
			continue
		}
		found := matches[0]
		name := found[1]
		version := found[2]
		converted[p] = component{format: format, name: name, version: version}
	}

	return converted, nil
}

func componentsFromNpm(lines map[int64]string) (map[int64]component, error) {
	re := regexp.MustCompile(`"([^"]*)": ".?([0-9](\.[0-9])+)",?`)
	return componentsSingleLineNameVersion(lines, re, "npm")
}

func componentsFromNuget(lines map[int64]string) (map[int64]component, error) {
	re := regexp.MustCompile(`<package id="([^"]*)" version="([^"]*)"`)
	return componentsSingleLineNameVersion(lines, re, "nuget")
}

func componentsFromPypi(lines map[int64]string) (map[int64]component, error) {
	re := regexp.MustCompile(`(.*)==([^\s#]*)`)
	log.Printf("TRACE: pypi adds: %v\n", lines)
	return componentsSingleLineNameVersion(lines, re, "pypi")
}

/*
func componentsFromGradle(lines map[int64]string) (map[int64]component, error) {
	reOld := regexp.MustCompile(`^.*group:\s*'([^']*)',\s+name:\s*'([^']*)',\s+version:\s*'([^']*)'\s*$`)
	reNew := regexp.MustCompile(`^[^\s(]*[\s(]["']([^:]*):([^:]*):([^:]*)["']\)?$`)
}
*/

func parsePatchAdditions(patch string) map[int64]string {
	adds := make(map[int64]string)

	scanner := bufio.NewScanner(strings.NewReader(patch))
	var position int64
	for scanner.Scan() {
		if scanner.Text()[0] == '+' {
			adds[position] = scanner.Text()[1:]
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
		case "packages.config":
			components, err = getComponents(f.Patch, componentsFromNuget)
		case "requirements.txt":
			log.Printf("TRACE: pypi patch: %s", f.Patch)
			components, err = getComponents(f.Patch, componentsFromPypi)
		}

		if err != nil {
			// TODO
			continue
		}

		manifests[f] = components
	}

	return manifests, nil
}
