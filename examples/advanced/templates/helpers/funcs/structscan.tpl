func StructScan(r *pgx.Row, x *{{ Table.GoStructName }}) error {
    return r.Scan(
        {% for c in Table.SelectableColumns %}
            &(x.{{ c.GoFieldName }}),
        {% endfor %}
    )
}
