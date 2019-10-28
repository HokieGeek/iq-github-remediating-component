package main

import (
	"bufio"
	"regexp"
	"strings"
)

func componentsSingleLineNameVersion(lines map[int64]string, re *regexp.Regexp, format string, fields []string) (map[int64]component, error) {
	converted := make(map[int64]component)

	for p, l := range lines {
		matches := re.FindAllStringSubmatch(l, -1)
		if len(matches) == 0 {
			continue
		}
		found := matches[0]
		c := component{format: format}
		for i, f := range fields {
			switch f {
			case "group":
				c.group = found[i+1]
			case "name":
				c.name = found[i+1]
			case "version":
				c.version = found[i+1]
			}
		}
		converted[p] = c
	}

	return converted, nil
}

func componentsFromNpm(lines map[int64]string) (map[int64]component, error) {
	re := regexp.MustCompile(`"([^"]*)": ".?([0-9](\.[0-9])+)",?`)
	return componentsSingleLineNameVersion(lines, re, "npm", []string{"name", "version"})
}

func componentsFromNuget(lines map[int64]string) (map[int64]component, error) {
	re := regexp.MustCompile(`<package id="([^"]*)" version="([^"]*)"`)
	return componentsSingleLineNameVersion(lines, re, "nuget", []string{"name", "version"})
}

func componentsFromPypi(lines map[int64]string) (map[int64]component, error) {
	re := regexp.MustCompile(`(.*)==([^\s#]*)`)
	return componentsSingleLineNameVersion(lines, re, "pypi", []string{"name", "version"})
}

func componentsFromGradle(lines map[int64]string) (map[int64]component, error) {
	reOld := regexp.MustCompile(`^.*group:\s*'([^']*)',\s+name:\s*'([^']*)',\s+version:\s*'([^']*)'\s*$`)
	reNew := regexp.MustCompile(`^[^\s(]*[\s(]["']([^:]*):([^:]*):([^:]*)["']\)?$`)
	fields := []string{"group", "name", "version"}
	components := make(map[int64]component)
	if comps, err := componentsSingleLineNameVersion(lines, reOld, "maven", fields); err == nil {
		for k, v := range comps {
			if _, ok := components[k]; !ok {
				components[k] = v
			}
			// TODO: what if it does find it?
		}
	}

	if comps, err := componentsSingleLineNameVersion(lines, reNew, "maven", fields); err == nil {
		for k, v := range comps {
			if _, ok := components[k]; !ok {
				components[k] = v
			}
			// TODO: what if it does find it?
		}
	}

	return components, nil
}

func parsePatchLineAdditions(patch string) map[int64]string {
	adds := make(map[int64]string)

	scanner := bufio.NewScanner(strings.NewReader(patch))
	var position int64
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}

		if scanner.Text()[0] == '+' {
			adds[position] = scanner.Text()[1:]
		}
		position++
	}

	return adds
}

func getMavenComponents(patch string) (map[int64]component, error) {
	components := make(map[int64]component)

	var (
		position, p int64
		comp        *component
		newVersion  bool
	)
	scanner := bufio.NewScanner(strings.NewReader(patch))
	tag := regexp.MustCompile(`<([^>]*)>([^<]*)</.*`)
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Printf("%s::: ", line)
		switch {
		case strings.Contains(line, "</dependency>") || line[:2] == "@@":
			if comp != nil && newVersion {
				// fmt.Printf("(created): %q\n", *comp)
				components[p] = *comp
			}
			comp = nil
		case strings.Contains(line, "<dependency>"):
			// fmt.Printf("(new)\n")
			comp = new(component)
			comp.format = "maven"
			newVersion = false
		default:
			if matches := tag.FindAllStringSubmatch(line, -1); len(matches) > 0 && comp != nil {
				t, v := matches[0][1], matches[0][2]
				// fmt.Printf("%q %s %s\n", matches, t, v)
				switch t {
				case "groupId":
					comp.group = v
				case "artifactId":
					comp.name = v
				case "version":
					newVersion = line[0] == '+'
					p = position
					comp.version = v
				}
			}
		}
		position++
	}

	return components, nil
}

func findComponentsFromManifest(files []githubPullRequestFile) (map[githubPullRequestFile]map[int64]component, error) {
	getComponents := func(patch string, linesToComponents func(lines map[int64]string) (map[int64]component, error)) (map[int64]component, error) {
		additions := parsePatchLineAdditions(patch)
		return linesToComponents(additions)
	}

	manifests := make(map[githubPullRequestFile]map[int64]component, 0)

	for _, f := range files {
		components := make(map[int64]component)
		// log.Printf("TRACE: %s patch:\n%s", f.Filename, f.Patch)
		var err error
		switch f.Filename {
		case "pom.xml":
			components, err = getMavenComponents(f.Patch)
		case "build.gradle":
			components, err = getComponents(f.Patch, componentsFromGradle)
		case "package.json":
			components, err = getComponents(f.Patch, componentsFromNpm)
		case "packages.config":
			components, err = getComponents(f.Patch, componentsFromNuget)
		case "requirements.txt":
			components, err = getComponents(f.Patch, componentsFromPypi)
			// case "go.sum":
			// case "go.mod":
		}

		if err != nil {
			// TODO
			continue
		}

		manifests[f] = components
	}

	return manifests, nil
}
