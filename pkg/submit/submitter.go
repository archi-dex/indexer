package submit

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
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

func submitToEndpoint(logger util.Logger, ws *websocket.Conn, entity parse.Entity) error {
	return ws.WriteJSON(entity)
}

func Submitter(logger util.Logger, entities <-chan parse.Entity, outputFile, outputEndpoint string, dryrun bool) error {
	if !dryrun && outputFile == "" && outputEndpoint == "" {
		logger.Debug("skipping as no outputs were provided")
		return nil
	}

	var file *os.File
	var ws *websocket.Conn
	var err error

	if outputFile != "" {
		file, err = os.OpenFile(outputFile, fileFlags, fileMode)
		if err != nil {
			logger.Fatalw("failed to open output file", "err", err)
		}

		logger.Debugw("opened file", "filename", file.Name())
		defer func() {
			if err := file.Close(); err != nil {
				logger.Fatalw("failed to close output file", "err", err)
			}
		}()
	}

	if outputEndpoint != "" {
		var resp *http.Response
		ws, resp, err = websocket.DefaultDialer.Dial(outputEndpoint, nil)
		if err != nil {
			logger.Fatalw("failed to connect to ws", "resp", resp, "err", err)
		}

		logger.Debugw("connected to ws", "address", ws.RemoteAddr().String(), "resp", resp)
		defer func() {
			if err := ws.Close(); err != nil {
				logger.Fatalw("failed to close ws", "err", err)
			}
		}()
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
				if err := submitToEndpoint(logger, ws, entity); err != nil {
					logger.Warnw("error posting to endpoint", "endpoint", outputEndpoint, "err", err)
					errors += 1
				}
			}
		}
	}

	return nil
}
