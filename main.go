package main

import (
	"errors"
	"flag"
	"fmt"
	"sqrible/src/sqrible"

	"github.com/jackc/pgx"
)

var configFile = flag.String("c", "", "/path/to/config/file.yaml")
var tableName = flag.String("t", "", "table name")
var templateDir = flag.String("d", "", "template dir")

func init() {
	flag.Parse()

	checkFlags()
}

func main() {
	cfg := sqrible.ParseConfig(*configFile)
	if !cfg.TableConfigurationProvided(*tableName) {
		sqrible.Quit(
			fmt.Errorf("Configuration for table %s not found in %s", *tableName, *configFile),
		)
	}

	conn := connectPG()
	defer conn.Close()

	t := sqrible.ProcessTable(conn, *tableName, cfg)

	buff, err := sqrible.ApplyTemplate(t, *templateDir, t.Template)
	if err != nil {
		sqrible.Quit(err)
	}

	fmt.Printf(string(buff))
}

func connectPG() *pgx.Conn {
	conn, err := pgx.Connect(pgConnConfig())
	if err != nil {
		sqrible.Quit(err)
	}

	return conn
}

func pgConnConfig() pgx.ConnConfig {
	cfg, err := pgx.ParseEnvLibpq()
	if err != nil {
		sqrible.Quit(err)
	}

	return cfg
}

func checkFlags() {
	if len(*configFile) <= 0 {
		sqrible.Quit(errors.New("-f flag not provided"))
	}

	if len(*tableName) <= 0 {
		sqrible.Quit(errors.New("-t flag not provided"))
	}

	if len(*templateDir) <= 0 {
		sqrible.Quit(errors.New("-d flag not provided"))
	}
}
