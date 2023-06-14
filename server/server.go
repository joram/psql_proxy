package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	Port int
	db   *sql.DB

	wg       sync.WaitGroup
	listener net.Listener
}

func (s *Server) handleNewClientConnection(conn net.Conn) {
	s.wg.Add(1)
	defer s.wg.Done()
	defer conn.Close()

	c, f := context.WithCancel(context.Background())
	handler := &ClientHandler{
		ClientConn: conn,
		Context:    c,
		Cancel:     f,
		AuthPhase:  PhaseStartup,
	}
	handler.handle()
}

func (s *Server) Listen() error {
	log.Println("listening on port", s.Port)

	// build listener
	portStr := fmt.Sprintf(":%d", s.Port)
	ln, err := net.Listen("tcp", portStr)
	if err != nil {
		return err
	}
	s.listener = ln
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go s.handleNewClientConnection(conn)
	}
	return nil
}
