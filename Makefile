run:
	go run ./*.go

connect_proxy:
	psql postgresql://postgres:postgres@localhost:8080/postgres?sslmode=disable

connect_db:
	psql postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable

test_queries:
	psql postgresql://postgres:postgres@localhost:8080/postgres?sslmode=disable -f ./sql_examples/create_table.sql
	psql postgresql://postgres:postgres@localhost:8080/postgres?sslmode=disable -f ./sql_examples/insert_data.sql
	psql postgresql://postgres:postgres@localhost:8080/postgres?sslmode=disable -f ./sql_examples/select_data.sql