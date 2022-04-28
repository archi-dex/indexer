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

func parse(logger util.Logger, matcher *regexp.Regexp, names []string, path string) *Entity {
	base := filepath.Base(path)
	matches := matcher.FindStringSubmatch(base)
	result := make(map[string]string, len(names)-1)

	if len(names) != len(matches) {
		logger.Warnw("did not match all named fields", "path", path, "names", names, "matches", matches)
		return nil
	}

	for i, name := range names {
		if name == path || name == "" {
			continue
		}

		if value := strings.Trim(matches[i], " "); value != "" {
			result[name] = value
		}
	}

	return &Entity{Filepath: path, Data: result}
}

func Parser(logger util.Logger, files <-chan string, entities chan<- Entity, pattern string) error {
	matcher := regexp.MustCompile(pattern)
	names := matcher.SubexpNames()

	total := 0
	skipped := 0
	defer func() { logger.Infow("parsing complete", "total", total, "skipped", skipped) }()

	for path := range files {
		total += 1

		if entity := parse(logger, matcher, names, path); entity != nil {
			logger.Debugw("parsed entity", "entity", entity)
			entities <- *entity
			continue
		}

		skipped += 1
	}

	return nil
}
