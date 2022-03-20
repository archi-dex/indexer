package submit

import (
	"encoding/json"
	"os"

	"github.com/hans-m-song/archidex/indexer/pkg/parse"
	"github.com/hans-m-song/archidex/indexer/pkg/util"
)

const (
	fileFlags = os.O_APPEND | os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	fileMode  = 0644
)

func submitToFile(file *os.File, entity parse.Entity) error {
	str, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	file.Write(str)
	file.WriteString("\n")

	return nil
}

func submitToEndpoint(endpoint string) error {
	return nil
}

func Submitter(logger util.Logger, entities <-chan parse.Entity, outputFile, outputEndpoint string, dryrun bool) error {
	if outputFile == "" && outputEndpoint == "" {
		logger.Debug("skipping as no outputs were provided")
		return nil
	}

	var file *os.File
	var err error

	if outputFile != "" {
		file, err = os.OpenFile(outputFile, fileFlags, fileMode)
		defer func() {
			if err := file.Close(); err != nil {
				logger.Fatalw("failed to close output file", "err", err)
			}
		}()

		if err != nil {
			logger.Fatalw("failed to open output file", "err", err)
		}
	}

	errors := 0
	defer func() { logger.Infow("submission completed", "errors", errors) }()

	for entity := range entities {
		logger.Debugw("submitting entity", "entity", entity)

		if !dryrun {
			if outputFile != "" {
				if err := submitToFile(file, entity); err != nil {
					logger.Warnw("error writing to file", "file", outputFile, "err", err)
					errors += 1
				}
			}

			if outputEndpoint != "" {
				if err := submitToEndpoint(outputEndpoint); err != nil {
					logger.Warnw("error posting to endpoint", "endpoint", outputEndpoint, "err", err)
					errors += 1
				}
			}
		}
	}

	return nil
}
