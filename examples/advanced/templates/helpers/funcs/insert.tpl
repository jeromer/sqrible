func insert(db Queryer, x *{{ Table.GoStructName }}) (*{{ Table.GoStructName }}, error) {
    args := pgx.QueryArgs(make([]interface{}, 0, {{ Table.Columns | length }}))

    var columns, values []string

    {% for c in Table.InsertableColumns %}
        if x.{{ c.GoFieldName }}.Status != pgtype.Undefined {
            columns = append(columns, `{{ c.PGColumnName }}`)
            values = append(values, args.Append(&(x.{{ c.GoFieldName }})))
        }
    {% endfor %}

    sql := `INSERT INTO "{{ Table.Name }}" (` + strings.Join(columns, ", ") + `)
			VALUES(` + strings.Join(values, ",") + `)
			RETURNING {{ Table.SelectableColumns.PGNames | join:", " }};`

    stmtName := PreparedName("SqribleInsert{{ Table.GoStructName }}", sql)

    x2 := New()
    err := StructScan(
        PrepareQueryRow(db, stmtName, sql, args...), x2,
    )

    if err != nil {
        return nil, err
    }

    return x2, err
}

