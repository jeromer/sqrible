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

	return c
}

type TableConfig struct {
	Template      string                              `yaml:"template"`
	ConfigDetails map[string]TableColumnConfigDetails `yaml:"tablecols"`
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

func (c Config) columnConfig(tn string, cn string) *TableColumnConfigDetails {
	tc, hasCfg := c.Tables[tn]
	if !hasCfg {
		return nil
	}

	tcd, found := tc.ConfigDetails[cn]
	if !found {
		return nil
	}

	return &tcd
}

type TableColumnConfigDetails struct {
	Access string `yaml:"access"`
}

func (d TableColumnConfigDetails) IsIgnored() bool {
	return d.Access == "-"
}

func (d TableColumnConfigDetails) IsSelectable() bool {
	return d.accessContains("s")
}

func (d TableColumnConfigDetails) IsInsertable() bool {
	return d.accessContains("i")
}

func (d TableColumnConfigDetails) IsUpdateable() bool {
	return d.accessContains("u")
}

func (d TableColumnConfigDetails) accessContains(s string) bool {
	return strings.Contains(d.Access, s)
}
