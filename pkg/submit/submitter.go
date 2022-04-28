package submit

import (
	"github.com/hans-m-song/archidex/indexer/pkg/parse"
	"github.com/hans-m-song/archidex/indexer/pkg/util"
)

type destination interface {
	targetName() string
	open() error
	close() error
	submit(entity parse.Entity) error
}

func Submitter(logger util.Logger, entities <-chan parse.Entity, outputFile, outputEndpoint string, dryrun bool) error {
	if !dryrun && outputFile == "" && outputEndpoint == "" {
		logger.Debug("skipping as no outputs were provided")
		return nil
	}

	destinations := []destination{
		&fileDestination{target: outputFile, logger: logger},
		&wsDestination{target: outputEndpoint, logger: logger},
	}

	for _, dest := range destinations {
		if err := dest.open(); err != nil {
			logger.Warnw("failed to open destination", "target", dest.targetName(), "err", err)
			return err
		}

		defer func(dest destination) {
			if err := dest.close(); err != nil {
				logger.Warnw("failed to close destination", "target", dest.targetName(), "err", err)
			}
		}(dest)
	}

	errors := 0
	defer func() { logger.Infow("submission completed", "errors", errors) }()

	for entity := range entities {
		logger.Debugw("submitting entity", "entity", entity)

		if !dryrun {
			for _, dest := range destinations {
				if err := dest.submit(entity); err != nil {
					logger.Warnw("error submitting entity", "entity", entity, "err", err)
					errors += 1
				}
			}
		}
	}

	return nil
}
