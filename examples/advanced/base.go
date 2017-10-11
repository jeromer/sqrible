package user

import (
	"fmt"
	"hash/fnv"
	"io"

	"github.com/jackc/pgx"
)

type Queryer interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (pgx.CommandTag, error)
}

type Preparer interface {
	Prepare(name, sql string) (*pgx.PreparedStatement, error)
}

func PrepareQueryRow(db Queryer, name, sql string, args ...interface{}) *pgx.Row {
	if p, ok := db.(Preparer); ok {
		// QueryRow doesn't return an error, the error is encoded in the pgx.Row.
		// Since that is private, Ignore the error from Prepare and run the query
		// without the prepared statement. It should fail with the same error.
		if _, err := p.Prepare(name, sql); err == nil {
			sql = name
		}
	}

	return db.QueryRow(sql, args...)
}

func PreparedName(baseName, sql string) string {
	h := fnv.New32a()
	if _, err := io.WriteString(h, sql); err != nil {
		// hash.Hash.Write never returns an error so this can't happen
		panic("failed writing to hash")
	}

	return fmt.Sprintf("%s%d", baseName, h.Sum32())
}
