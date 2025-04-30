task docker

task up

.env:

LVL_DEPLOYMENT=debug

DB_PASS=123
DB_USER=user
DB_NAME=dbname
DB_PORT=5433

DB_PATH="postgres://user:123@localhost:5433/dbname?sslmode=disable&search_path=public"
