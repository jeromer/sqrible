package sqrible

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/jackc/pgx"
)

func ProcessTable(c *pgx.Conn, n string, cfg Config) *Table {
	if !tableExists(c, n) {
		Quit(fmt.Errorf("Table %s not found", n))
	}

	cols, err := tableColumns(c, n, cfg)
	if err != nil {
		Quit(err)
	}

	tcfg := cfg.tableConfig(n)

	return NewTable(n, cols, tcfg)
}

func Quit(e error) {
	fmt.Fprintln(os.Stderr, red("ERROR"), e)
	os.Exit(1)
}

func ApplyTemplate(t *Table, templateDir string, templateName string) ([]byte, error) {
	buff := []byte{}
	s := pongo2.NewSet("sqrible")
	err := s.SetBaseDirectory(templateDir)
	if err != nil {
		Quit(err)
		return buff, err
	}

	tpl, err := s.FromFile(templateName)
	if err != nil {
		Quit(err)
		return buff, err
	}

	return tpl.ExecuteBytes(pongo2.Context{
		"Table": t,
	})
}

func tableExists(c *pgx.Conn, n string) bool {
	count := new(int)

	err := c.QueryRow(
		`SELECT COUNT(*)
		 FROM pg_catalog.pg_class c
			  LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		 WHERE c.relkind IN ('r','v','m','S','f','')
			  AND n.nspname <> 'pg_catalog'
			  AND n.nspname <> 'information_schema'
			  AND n.nspname !~ '^pg_toast'
			  AND pg_catalog.pg_table_is_visible(c.oid)
			  AND c.relname = $1;`,
		n,
	).Scan(count)

	if err != nil {
		Quit(err)
		return false
	}

	return *count >= 1
}

func tableColumns(c *pgx.Conn, name string, cfg Config) ([]*Column, error) {
	var columns []*Column

	pks, err := tablePKs(c, name)
	if err != nil {
		return columns, err
	}

	rows, err := c.Query(
		`SELECT column_name,
				data_type,
				udt_name,
				ordinal_position,
				is_nullable, 
				is_updatable
		FROM information_schema.columns
		WHERE table_name=$1
		ORDER BY ordinal_position ASC`,
		name,
	)

	if err != nil {
		return columns, err
	}

	defer rows.Close()
	for rows.Next() {
		c := new(Column)

		err = rows.Scan(
			&c.PGColumnName,
			&c.PGDataType,
			&c.PGUDTName,
			&c.PGOrdinalPosition,
			&c.PGIsNullable,
			&c.PGIsUpdatable,
		)

		if err != nil {
			return []*Column{}, err
		}

		c.GoFieldName = asGoFieldName(c.PGColumnName)
		c.PgxType = asPgxType(c.PGDataType, c.PGUDTName, c.PGIsNullable)

		c.Config = cfg.columnConfig(name, c.PGColumnName)
		c.IsPK = colIsPk(c.PGColumnName, pks)

		if c.isIgnored() {
			continue
		}

		columns = append(columns, c)
	}

	if rows.Err() != nil {
		return []*Column{}, err
	}

	return columns, nil
}

func colIsPk(pgCol string, pks []string) bool {
	for _, pk := range pks {
		if pgCol == pk {
			return true
		}
	}

	return false
}

func tablePKs(c *pgx.Conn, tableName string) ([]string, error) {
	rows, err := c.Query(
		`SELECT a.attname
			FROM   pg_index i
			JOIN   pg_attribute a ON a.attrelid = i.indrelid
								 AND a.attnum = ANY(i.indkey)
			WHERE  i.indrelid = $1::regclass
			AND    i.indisprimary`,
		tableName,
	)

	if err != nil {
		return []string{}, err
	}

	defer rows.Close()

	pks := []string{}
	for rows.Next() {
		pk := new(string)
		err = rows.Scan(pk)
		if err != nil {
			return []string{}, err
		}

		pks = append(pks, *pk)
	}

	return pks, nil
}

func asGoFieldName(n string) string {
	parts := strings.Split(n, "_")
	buf := &bytes.Buffer{}

	for _, s := range parts {
		if isAcronym(s) {
			s = strings.ToUpper(s)
		} else {
			s = strings.Title(s)
		}

		buf.WriteString(s)
	}

	return buf.String()
}

func asPgxType(n string, udt string, nullable string) string {
	// These are the field transpositions IF
	// the field is marked nullable

	nulls := map[string]string{
		"int16":             "dat.NullInt64",
		"int8":              "dat.NullInt64",
		"bigint":            "dat.NullInt64",
		"bigserieal":        "dat.NullInt64",
		"int4":              "dat.NullInt64",
		"int":               "dat.NullInt64",
		"serial":            "dat.NullInt64",
		"int2":              "dat.NullInt64",
		"smallint":          "dat.NullInt64",
		"smallserial":       "dat.NullInt64",
		"double precision":  "dat.NullFloat64",
		"real":              "dat.NullFloat64",
		"bool":              "bool",
		"boolean":           "bool",
		"text":              "dat.NullString",
		"varchar":           "dat.NullString",
		"citext":            "dat.NullString",
		"character varying": "dat.NullString",
		"timestamp":         "dat.NullTime",
		"inet":              "inet",
		"uuid":              "uuid.UUID",
		"bytea":             "pgtype.ByteaArray",
		"jsonb":             "dat.JSON",
		"geography":         "pgtype.Point",
		"tsvector":          "dat.NullString",
	}

	full := map[string]string{


		"int16":             "int16",
		"int8":              "int8",
		"bigint":            "int8",
		"bigserieal":        "int8",
		"int4":              "int",
		"int":               "int",
		"serial":            "int",
		"int2":              "int",
		"smallint":          "int",
		"smallserial":       "int",
		"double precision":  "float64",
		"real":              "float64",
		"bool":              "bool",
		"boolean":           "bool",
		"text":              "string",
		"varchar":           "string",
		"citext":            "string",
		"character varying": "string",
		"timestamp":         "dat.NullTime",
		"inet":              "inet",
		"uuid":              "uuid.UUID",
		"bytea":             "[][]byte",
		"jsonb":             "dat.JSON",
		"geography":         "pgtype.Point",
		"tsvector":          "string",
	}

	if nullable == "YES" {

		t, found := nulls[udt]
		if found {
			return t
		}
		if !found {
			t, found := nulls[n]
			if found {
				return t
			}
		}
	} else {

		t, found := full[udt]
		if found {
			return t
		}

		if !found {
			t, found := full[n]
			if found {
				return t
			}
		}
	}

	if strings.ToLower(n) == "array" {
		return asPgxType(strings.ToLower(udt[1:]), udt, nullable) + "Array"
	}

	Quit(fmt.Errorf("Postgres type %s not found in pgx mapping", n))
	return ""
}

func isAcronym(s string) bool {
	acronyms := []string{
		"id", "ip", "url", "uid", "os",
	}

	for _, a := range acronyms {
		if a == s {
			return true
		}
	}

	return false
}

func red(s string) string {
	return "\033[1;31m" + s + "\033[0m"
}
