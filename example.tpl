Example template
----------------

Table:
- name: {{ Table.Name }}
- go struct name: {{ Table.GoStructName }}
- template used: {{ Table.Template }}

{% macro dump_columns(cols) %}
{% for c in cols %}{{ dump_column(c) }}{% endfor %}
{% endmacro %}

{% macro dump_column(c) %}PGColumnName      : {{ c.PGColumnName }}
PGDataType        : {{ c.PGDataType }}
PGOrdinalPosition : {{ c.PGOrdinalPosition }}
IsPK              : {{ c.IsPK }}
GoFieldName       : {{ c.GoFieldName }}
PgxType           : {{ c.PgxType }}
JSON              : {{ c.JSON}}

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
{% for c in Table.Columns%}{{ dump_column(c) }}{% endfor %}
