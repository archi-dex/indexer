package submit

import (
	"github.com/gorilla/websocket"
	"github.com/hans-m-song/archidex/indexer/pkg/parse"
	"github.com/hans-m-song/archidex/indexer/pkg/util"
)

var (
	_ destination = (*wsDestination)(nil)
)

type wsDestination struct {
	logger util.Logger
	target string
	dest   *websocket.Conn
}

func (d *wsDestination) targetName() string {
	return d.target
}

func (d *wsDestination) open() error {
	if d.target == "" {
		return nil
	}

	var err error
	if d.dest, _, err = websocket.DefaultDialer.Dial(d.target, nil); err != nil {
		return err
	}

	return nil
}

func (d *wsDestination) close() error {
	if d.dest == nil {
		return nil
	}

	if err := d.dest.Close(); err != nil {
		return err
	}

	return nil
}

func (d *wsDestination) submit(entity parse.Entity) error {
	if d.dest == nil {
		return nil
	}

	return d.dest.WriteJSON(entity)
}
