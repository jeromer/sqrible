package sqrible

import (
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

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
			map[string]TableColumnConfig, len(tableConfig.ConfigDetails),
		)

		for colName, details := range tableConfig.ConfigDetails {
			newTableConfig.TableCols[colName] = TableColumnConfig{
				IsIgnored:    details.isIgnored(),
				IsSelectable: details.isSelectable(),
				IsInsertable: details.isInsertable(),
				IsUpdateable: details.isUpdateable(),
			}
		}

		c.Tables[tableName] = *newTableConfig
	}

	return c
}

type TableColumnConfigDetails struct {
	Access string `yaml:"access"`
}

type TableColumnConfig struct {
	IsIgnored    bool
	IsSelectable bool
	IsInsertable bool
	IsUpdateable bool
}

type TableConfig struct {
	Template      string                              `yaml:"template"`
	ConfigDetails map[string]TableColumnConfigDetails `yaml:"tablecols"`
	TableCols     map[string]TableColumnConfig        `yaml:"-"`
	GoStruct      string
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

func (d TableColumnConfigDetails) isIgnored() bool {
	return d.Access == "-"
}

func (d TableColumnConfigDetails) isSelectable() bool {
	return d.accessContains("s")
}

func (d TableColumnConfigDetails) isInsertable() bool {
	return d.accessContains("i")
}

func (d TableColumnConfigDetails) isUpdateable() bool {
	return d.accessContains("u")
}

func (d TableColumnConfigDetails) accessContains(s string) bool {
	return strings.Contains(d.Access, s)
}
