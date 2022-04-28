package submit

import (
	"encoding/json"
	"os"

	"github.com/hans-m-song/archidex/indexer/pkg/parse"
	"github.com/hans-m-song/archidex/indexer/pkg/util"
)

var (
	_ destination = (*fileDestination)(nil)
)

const (
	fileFlags = os.O_APPEND | os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	fileMode  = 0644
)

type fileDestination struct {
	logger util.Logger
	target string
	dest   *os.File
}

func (d *fileDestination) targetName() string {
	return d.target
}

func (d *fileDestination) open() error {
	if d.target == "" {
		return nil
	}

	var err error
	if d.dest, err = os.OpenFile(d.target, fileFlags, fileMode); err != nil {
		return err
	}

	return nil
}

func (d *fileDestination) close() error {
	if d.dest == nil {
		return nil
	}

	if err := d.dest.Close(); err != nil {
		return err
	}

	return nil
}

func (d *fileDestination) submit(entity parse.Entity) error {
	if d.dest == nil {
		return nil
	}

	var str []byte
	var err error
	if str, err = json.Marshal(entity); err != nil {
		return err
	}

	if _, err := d.dest.Write(str); err != nil {
		return err
	}

	if _, err := d.dest.WriteString("\n"); err != nil {
		return err
	}

	return nil
}
