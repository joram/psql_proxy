package main

import "github.com/joram/psql_proxy/server"

// lexing/parsing strings: https://notes.eatonphil.com/database-basics.html
// client connection handling: https://github.com/rueian/pgbroker/blob/master/proxy/server.go#L38

func main() {
	s := &server.Server{
		Port: 8080,
	}
	err := s.Listen()
	if err != nil {
		panic(err)
	}
}
