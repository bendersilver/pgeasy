package pgeasy

import (
	"fmt"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/bendersilver/sqleasy"
)

// TableRules -
type TableRules struct {
	Sheme   string `toml:"sheme"`
	Table   string `toml:"table"`
	InitSQL string `toml:"init_sql"`
}

// CleanRules -
type CleanRules struct {
	SQL     string `toml:"sql"`
	Timeout int    `toml:"timeout"`
}

// Config -
type Config struct {
	db     *sqleasy.Conn
	Name   string `toml:"slot_name"`
	PgURL  string `toml:"pg_url"`
	Server struct {
		Network string `toml:"network"`
		Addr    string `toml:"addr"`
	} `toml:"server"`
	TableRules []TableRules `toml:"table_rules"`
	CleanRules []CleanRules `toml:"clean_rules"`
}

// readConf -
func readConf() error {
	_, err := toml.DecodeFile(path.Join(os.Getenv("CONF_PATH"), "pgeasy.conf"), &cfg)
	if err != nil {
		return err
	}
	if cfg.PgURL == "" {
		return fmt.Errorf("pg_url parameters not set")
	}
	if cfg.Name == "" {
		cfg.Name = "pgeasy_slot"
	}

	if cfg.Server.Network == "" {
		cfg.Server.Network = "unix"
	}

	if cfg.Server.Addr == "" {
		switch cfg.Server.Network {
		case "unix":
			cfg.Server.Addr = "/tmp/pgeasy.sock"
		case "tcp":
			cfg.Server.Addr = "localhost:9021"
		default:
			return fmt.Errorf("wrong config server.network, allowed values unix|tcp")
		}
	}

	return writeConf()
}

// writeConf -
func writeConf() error {
	f, err := os.OpenFile(path.Join(os.Getenv("CONF_PATH"), "pe.cfg"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString("# config auto generated; DO NOT EDIT\n\n")

	enc := toml.NewEncoder(f)
	enc.Indent = "\t"
	return enc.Encode(&cfg)
}

// appendTableConf -
func appendTableConf() {}
