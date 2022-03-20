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

func submitToFile(file *os.File, entity parse.Entity) {
	str, _ := json.Marshal(entity)
	file.Write(str)
	file.WriteString("\n")
}

func submitToEndpoint(endpoint string) {}

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

	for entity := range entities {
		logger.Debugw("submitting entity", "entity", entity)

		if !dryrun {
			if outputFile != "" {
				submitToFile(file, entity)
			}

			if outputEndpoint != "" {
				submitToEndpoint(outputEndpoint)
			}
		}
	}

	return nil
}
