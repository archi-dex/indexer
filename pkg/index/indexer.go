package index

import (
	"io/fs"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hans-m-song/archidex/indexer/pkg/util"
)

func IndexFromWalk(logger util.Logger, files chan<- string, cwd, dir, ignorePattern, matchPattern string) error {
	ignore := regexp.MustCompile(ignorePattern)
	match := regexp.MustCompile(matchPattern)

	logger.Infow("indexing directory", "cwd", cwd, "dir", dir)
	directory := path.Join(cwd, dir)

	total := 0
	skipped := 0
	errors := 0
	defer func() { logger.Infow("indexing complete", "total", total, "skipped", skipped, "errors", errors) }()

	return filepath.WalkDir(directory, func(filename string, d fs.DirEntry, err error) error {
		total += 1
		relative := strings.Replace(filename, cwd, "", 1)

		if err != nil {
			logger.Warnw("error walking", "path", relative, "err", err)
			errors += 1
			return err
		}

		if ignorePattern != "" && ignore.MatchString(relative) {
			skipped += 1
			return nil
		}

		if match.MatchString(relative) {
			logger.Debugw("discovered path", "path", relative)
			files <- relative
			return nil
		}

		skipped += 1
		return nil
	})
}
