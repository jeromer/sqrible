func fetchOne(db Queryer, fieldName string, fieldValue interface{}) (*{{ Table.GoStructName}}, error) {
    sql := `SELECT {{ Table.SelectableColumns.PGNames | join:", " }}
          FROM {{ Table.Name }}
          WHERE ` + fieldName + `=$1;`

    stmtName := PreparedName(
        "sqribleFetchOne{{ Table.GoStructName }}By" + fieldName, sql,
    )

    x := New()

    err := StructScan(
        PrepareQueryRow(db, stmtName, sql, fieldValue), x,
    )

    if err == pgx.ErrNoRows {
        return nil, ErrNotFound
    } else if err != nil {
        return nil, err
    }

    return x, nil
}
