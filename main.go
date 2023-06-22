package main

import "github.com/joram/psql_proxy/server"

// lexing/parsing strings: https://notes.eatonphil.com/database-basics.html
// client connection handling: https://github.com/rueian/pgbroker/blob/master/proxy/server.go#L38
// psql wire protocol in go: https://pkg.go.dev/github.com/jeroenrinzema/psql-wire#section-readme

func main() {
	s, err := server.NewServer(8080, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	err = s.Run()
	if err != nil {
		panic(err)
	}

	for {
		// do nothing
	}
}
