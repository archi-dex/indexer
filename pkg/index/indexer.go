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

	return filepath.WalkDir(directory, func(filename string, d fs.DirEntry, err error) error {
		relative := strings.Replace(filename, cwd, "", 1)
		if err != nil {
			logger.Warnw("error walking", "path", relative, "err", err)
			return err
		}

		if ignorePattern != "" && ignore.MatchString(relative) {
			return nil
		}

		if match.MatchString(relative) {
			logger.Debugw("discovered path", "path", relative)
			files <- relative
		}

		return nil
	})
}
