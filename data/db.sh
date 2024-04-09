rm -f sqlite.db
touch sqlite.db
sqlite3 sqlite.db < schema.sql
