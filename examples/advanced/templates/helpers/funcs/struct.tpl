type {{ Table.GoStructName }} struct {
    {% for c in Table.SelectableColumns %}
        {{ c.GoFieldName }} {{ c.PgxType }} // {{ Table.Name }}.{{ c.PGColumnName }}({{ c.PGDataType}})
    {% endfor %}
}

