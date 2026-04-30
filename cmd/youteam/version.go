package main

import (
	"regexp"
	"strings"
)

const (
	defaultVersion    = "0.0.0"
	defaultCommit     = "dev"
	semverPatternText = `^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-((0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*)(\.(0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*))*))?(\+([0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*))?$`
)

var semverPattern = regexp.MustCompile(semverPatternText)

func displayVersion(versionValue, commitValue string) string {
	normalizedCommit := strings.TrimSpace(commitValue)
	if normalizedCommit == "" {
		normalizedCommit = defaultCommit
	}

	return normalizeVersion(versionValue) + "-" + normalizedCommit
}

func normalizeVersion(versionValue string) string {
	normalizedVersion := strings.TrimPrefix(strings.TrimSpace(versionValue), "v")
	if semverPattern.MatchString(normalizedVersion) {
		return normalizedVersion
	}

	return defaultVersion
}
