1. Import `./users.sql` in postgres
2. Build sqrible : `make -C ../../ build`
3. Run: `PGDATABASE=yourdb PGHOST=localhost ../../bin/sqrible -c ./sqrible.yml -d . -t users`
