package pgeasy

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/bendersilver/pgcopy"
	"github.com/bendersilver/sqleasy"
	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgtype"
)

func createTable(msg *pglogrepl.RelationMessage) (*sqleasy.Stmt, error) {
	var dbName = fmt.Sprintf("%s_%s", msg.Namespace, msg.RelationName)

	err := cfg.db.Exec("DROP TABLE IF EXISTS " + dbName)
	if err != nil {
		return nil, err
	}

	var cols, pk, mark []string
	for _, col := range msg.Columns {
		mark = append(mark, "?")
		switch col.DataType {
		case pgtype.BoolOID:
			cols = append(cols, col.Name+" BOOLEAN")
		case pgtype.Int2OID, pgtype.Int4OID, pgtype.Int8OID, pgtype.TimestampOID, pgtype.TimestamptzOID, pgtype.DateOID:
			cols = append(cols, col.Name+" INTEGER")
		case pgtype.NumericOID, pgtype.Float4OID, pgtype.Float8OID:
			cols = append(cols, col.Name+" REAL")
		case pgtype.TextOID, pgtype.VarcharOID, pgtype.NameOID:
			cols = append(cols, col.Name+" TEXT")
		default:
			cols = append(cols, col.Name+" BLOB")
		}
		if col.Flags == 1 {
			pk = append(pk, col.Name)
		}
	}

	err = cfg.db.Exec(fmt.Sprintf(`CREATE TABLE %s (\n%s\n,PRIMARY KEY (%s)\n) WITHOUT ROWID;`,
		dbName, strings.Join(cols, ",\n"), strings.Join(pk, ","),
	))
	if err != nil {
		return nil, err
	}
	return cfg.db.Prepare(
		fmt.Sprintf(
			"INSERT INTO %s VALUES (%s)", dbName, strings.Join(mark, ","),
		),
	)
}

func copyTable(t *TableRules) error {
	c, err := pgcopy.New(cfg.PgURL, t.Sheme, t.Table)
	if err != nil {
		return err
	}
	defer c.Close()

	insert, err := createTable(c.RelationMessage())
	if err != nil {
		return err
	}

	var dbName = fmt.Sprintf("%s.%s", t.Sheme, t.Table)

	for _, v := range []string{"DROP", "ADD"} {
		err = c.Exec(fmt.Sprintf("ALTER PUBLICATION %s %s TABLE %s;", cfg.Name, v, dbName))
	}
	if err != nil {
		return err
	}

	if t.InitSQL == "" {
		t.InitSQL = fmt.Sprintf("SELECT * FROM %s.%s;", t.Sheme, t.Table)
	}
	err = c.Read(t.InitSQL, func(vals []driver.Value) error {
		return insert.Exec(vals...)
	})

	return nil
}

func initCopy() (err error) {
	for _, tb := range cfg.TableRules {
		err := copyTable(&tb)
		if err != nil {
			return err
		}
	}
	return nil
}
