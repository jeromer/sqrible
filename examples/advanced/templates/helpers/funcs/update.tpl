func update(db Queryer, x *{{ Table.GoStructName }}) (*{{ Table.GoStructName }}, error) {
	sets := make([]string, 0, {{ Table.UpdateableColumns | length}})
	args := pgx.QueryArgs(make([]interface{}, 0, {{ Table.UpdateableColumns | length + Table.Pks | length}}))

	{% for c in Table.UpdateableColumns %}
	if x.{{ c.GoFieldName }}.Status != pgtype.Undefined {
		sets = append(sets, `{{ c.PGColumnName }}`+"="+args.Append(&(x.{{ c.GoFieldName }})))
	}
	{% endfor %}

	if len(sets) == 0 {
		return x, nil
	}

	sql := `UPDATE "{{ Table.Name }}" SET ` + strings.Join(sets, ", ")
	sql += ` WHERE `
	sql += {% for c in Table.Pks %}
			 `{{ c.PGColumnName }}=` +args.Append(&(x.{{ c.GoFieldName }})) {% if forloop.Last == false %}+` AND `+{% endif %}
		  {% endfor %}
	sql += ` RETURNING {{ Table.SelectableColumns.PGNames | join:", " }};`

	stmtName := PreparedName("sqribleUpdate{{ Table.GoStructName }}", sql)

	x2 := New()
    err := StructScan(
        PrepareQueryRow(db, stmtName, sql, args...), x2,
    )

    if err != nil {
        return nil, err
    }

	return x2, err
}
