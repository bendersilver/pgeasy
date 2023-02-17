package pgeasy

import (
	"github.com/bendersilver/pglr"
	"github.com/bendersilver/sqleasy"
)

// EasyPG -
type EasyPG struct {
	db *sqleasy.Conn
	pl *pglr.Conn
}

var cfg Config

// InitConf -
func InitConf() error {
	return readConf()

}

// Start -
func Start() (err error) {
	var pe EasyPG
	pe.db, err = sqleasy.New()
	if err != nil {
		return err
	}

	pe.pl, err = pglr.NewConn(&pglr.Options{
		PgURL:     cfg.PgURL,
		Temporary: true,
	})
	if err != nil {
		return err
	}
	return nil
}

// Close -
func Close() {

}
