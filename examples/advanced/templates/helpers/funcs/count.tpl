func Count(db Queryer) int {
    var n int

    sql := "SELECT COUNT(*) FROM {{ Table.Name }}"

    stmtName := PreparedName("SqribleCount{{ Table.GoStructName }}", sql)
    err := PrepareQueryRow(db, stmtName, sql).Scan(&n)

    if err != nil{
        return 0
    }

    return n
}
