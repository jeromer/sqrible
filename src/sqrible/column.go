package sqrible

type Column struct {
	PGColumnName      string
	PGDataType        string
	PGUDTName         string
	PGOrdinalPosition int32

	IsPK bool

	GoFieldName string
	PgxType     string

	Config *TableColumnConfig
}

func (c *Column) isIgnored() bool {
	return (c.IsConfigured() && c.Config.IsIgnored)
}

func (c *Column) IsConfigured() bool {
	return c.Config != nil
}

func (c *Column) IsSelectable() bool {
	if !c.IsConfigured() {
		return true
	}

	return c.Config.IsSelectable
}

func (c *Column) IsInsertable() bool {
	if !c.IsConfigured() {
		return true
	}

	return c.Config.IsInsertable
}

func (c *Column) IsUpdateable() bool {
	if !c.IsConfigured() {
		return true
	}

	return c.Config.IsUpdateable
}
