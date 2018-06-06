func fetchOneWithSQL(db Queryer, sql string, args ...interface{}) (*{{ Table.GoStructName}}, error) {
    stmtName := PreparedName(
        "sqribleFetchOneWithSQL{{ Table.GoStructName }}", sql,
    )

    x := New()

    err := StructScan(
        PrepareQueryRow(db, stmtName, sql, args...), x,
    )

    if err == pgx.ErrNoRows {
        return nil, ErrNotFound
    } else if err != nil {
        return nil, err
    }

    return x, nil
}
