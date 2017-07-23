package sqrible

import (
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type TableColumConfigFlag string

type TableColumnConfig struct {
	IsIgnored    bool
	IsSelectable bool
	IsInsertable bool
	IsUpdateable bool
}

type TableConfig struct {
	Template  string                          `yaml:"template"`
	ColFlags  map[string]TableColumConfigFlag `yaml:"tablecols"`
	TableCols map[string]TableColumnConfig    `yaml:"-"`
	GoStruct  string
}

type Config struct {
	Tables map[string]TableConfig
}

func (c Config) tableConfig(tn string) *TableConfig {
	cfg, found := c.Tables[tn]
	if !found {
		return nil
	}

	return &cfg
}

func (c Config) TableConfigurationProvided(tn string) bool {
	return c.tableConfig(tn) != nil
}

func (c Config) columnConfig(tn string, cn string) *TableColumnConfig {
	tc, hasCfg := c.Tables[tn]
	if !hasCfg {
		return nil
	}

	tcc, found := tc.TableCols[cn]
	if !found {
		return nil
	}

	return &tcc
}

func ParseConfig(f string) Config {
	c := Config{}

	data, err := ioutil.ReadFile(f)
	if err != nil {
		Quit(err)
	}

	err = yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		Quit(err)
	}

	for tableName, tableConfig := range c.Tables {
		newTableConfig := new(TableConfig)
		*newTableConfig = tableConfig
		newTableConfig.TableCols = make(
			map[string]TableColumnConfig, len(tableConfig.ColFlags),
		)

		for colName, flags := range tableConfig.ColFlags {
			newTableConfig.TableCols[colName] = TableColumnConfig{
				IsIgnored:    flags.isIgnored(),
				IsSelectable: flags.isSelectable(),
				IsInsertable: flags.isInsertable(),
				IsUpdateable: flags.isUpdateable(),
			}
		}

		c.Tables[tableName] = *newTableConfig
	}

	return c
}

func (f TableColumConfigFlag) isIgnored() bool {
	return f == "-"
}

func (f TableColumConfigFlag) isSelectable() bool {
	return f.contains("s")
}

func (f TableColumConfigFlag) isInsertable() bool {
	return f.contains("i")
}

func (f TableColumConfigFlag) isUpdateable() bool {
	return f.contains("u")
}

func (f TableColumConfigFlag) contains(s string) bool {
	return strings.Contains(string(f), s)
}
