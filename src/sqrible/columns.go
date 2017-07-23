package sqrible

type Columns []*Column

func (c Columns) PGNames() []string {
	return c.buff(
		func(col *Column) string { return col.PGColumnName },
	)
}

func (c Columns) GoNames() []string {
	return c.buff(
		func(col *Column) string { return col.GoFieldName },
	)
}

func (c Columns) buff(f func(*Column) string) []string {
	buff := make([]string, len(c))

	for i, col := range c {
		buff[i] = f(col)
	}

	return buff

}
