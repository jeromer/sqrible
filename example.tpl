Example template
----------------

Table:
- name: {{ Table.Name }}
- go struct name: {{ Table.GoStructName }}
- template used: {{ Table.Template }}

{% macro dump_columns(cols) %}
{% for c in cols %}{{ dump_column(c) }}{% endfor %}
{% endmacro %}

{% macro dump_column(c) %}
{{ c.PGColumnName }} {{ c.PGDataType }} {{ c.PGOrdinalPosition }} {{ c.IsPK }} {{ c.GoFieldName }} {{ c.PgxType }}
{% endmacro %}

SELECTable columns
------------------
{{ dump_columns(Table.SelectableColumns) }}

INSERTable columns
------------------
{{ dump_columns(Table.InsertableColumns) }}

UPDATEable columns
------------------
{{ dump_columns(Table.UpdateableColumns) }}

Primary keys
------------
{% for c in Table.Pks %}{{ dump_column(c) }}{% endfor %}

all columns
-----------
{% for c in Table.Columns %}
{{ c.PGColumnName }} {{ c.PGDataType }} {{ c.PGOrdinalPosition }} {{ c.IsPK }} {{ c.GoFieldName }} {{ c.PgxType }} {{ c.IsConfigured }} {{ c.IsSelectable }} {{ c.IsInsertable }} {{ c.IsUpdateable}}
{% endfor %}
