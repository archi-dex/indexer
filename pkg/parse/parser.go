package parse

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hans-m-song/archidex/indexer/pkg/util"
)

type Entity struct {
	Filepath string            `json:"filepath"`
	Data     map[string]string `json:"data"`
}

func parse(logger util.Logger, matcher *regexp.Regexp, names []string, filename string) *Entity {
	base := filepath.Base(filename)
	matches := matcher.FindStringSubmatch(base)
	result := make(map[string]string, len(names)-1)

	if len(names) != len(matches) {
		logger.Warnw("did not match all named fields", "filename", filename, "names", names, "matches", matches)
		return nil
	}

	for i, name := range names {
		if name == filename || name == "" {
			continue
		}

		value := strings.Trim(matches[i], " ")
		if value == "" {
			continue
		}

		result[name] = value
	}

	return &Entity{Filepath: filename, Data: result}
}

func Parser(logger util.Logger, files <-chan string, entities chan<- Entity, pattern string) error {
	matcher := regexp.MustCompile(pattern)
	names := matcher.SubexpNames()

	for filepath := range files {
		if entity := parse(logger, matcher, names, filepath); entity != nil {
			logger.Debugw("parsed entity", "entity", entity)
			entities <- *entity
		}
	}

	return nil
}
