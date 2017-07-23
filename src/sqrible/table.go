package sqrible

type Table struct {
	Name         string
	Columns      Columns
	GoStructName string
	Template     string
}

func NewTable(name string, cols Columns, cfg *TableConfig) *Table {
	return &Table{
		Name:         name,
		Columns:      cols,
		GoStructName: cfg.GoStruct,
		Template:     cfg.Template,
	}
}

func (t *Table) Pks() Columns {
	buff := Columns{}
	for _, c := range t.Columns {
		if c.IsPK {
			buff = append(buff, c)
		}
	}

	return buff
}

func (t *Table) SelectableColumns() Columns {
	return t.filterColumns(
		func(c *Column) bool {
			return c.IsSelectable()
		},
	)
}

func (t *Table) InsertableColumns() Columns {
	return t.filterColumns(
		func(c *Column) bool {
			return c.IsInsertable()
		},
	)
}

func (t *Table) UpdateableColumns() Columns {
	return t.filterColumns(
		func(c *Column) bool {
			return c.IsUpdateable()
		},
	)
}

func (t *Table) filterColumns(f func(*Column) bool) Columns {
	buff := Columns{}

	for _, c := range t.Columns {
		if f(c) {
			buff = append(buff, c)
		}
	}

	return buff
}
