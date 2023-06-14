run:
	go run ./*.go

connect:
	psql postgresql://john:password@localhost:8080/mydb?sslmode=disable