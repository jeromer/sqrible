func StructScanRows(rows *pgx.Rows) ([]*{{ Table.GoStructName }}, error) {
    defer rows.Close()

    var buff[]*{{ Table.GoStructName }}

    var err error

    for rows.Next() {
        x := New()
        err = rows.Scan(
            {% for c in Table.SelectableColumns %}
                &(x.{{ c.GoFieldName }}),
            {% endfor %}
        )

        if err != nil {
            return []*{{ Table.GoStructName }}{}, err
        }

        buff = append(buff, x)
    }

    return buff, err
}
