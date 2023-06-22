package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	wire "github.com/jeroenrinzema/psql-wire"
	_ "github.com/lib/pq"
	"github.com/lib/pq/oid"
	"log"
)

type Server struct {
	port    int
	address string
	db      *sql.DB
	dbURL   string
}

type QueryFuncs struct {
	queryHandler        wire.PreparedStatementFn
	queryColumnsHandler QueryColumnsHandler
}
type QueryColumnsHandler func(ctx context.Context, query string) wire.Columns
type QueryHandler func(ctx context.Context, query string) QueryFuncs

func NewServer(port int, server string) (*Server, error) {
	db, err := sql.Open("postgres", server)
	if err != nil {
		log.Fatal(err)
	}

	return &Server{
		port:    port,
		address: fmt.Sprintf("127.0.0.1:%d", port),
		db:      db,
		dbURL:   server,
	}, nil
}

func (s *Server) Run() error {
	wire.ListenAndServe(s.address, s.handler)
	return nil
}

func (s *Server) selectQueryHandler(ctx context.Context, query string) QueryFuncs {
	statementHandler := func(ctx context.Context, writer wire.DataWriter, parameters []string) error {
		rows, err := s.db.Query(query)
		defer rows.Close()
		if err != nil {
			return err
		}

		columns, _ := rows.Columns()
		count := len(columns)
		values := make([]any, count)
		valuePtrs := make([]interface{}, count)
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		for rows.Next() {
			_ = rows.Scan(valuePtrs...)
			for i, val := range values {
				if val == nil {
					values[i] = "NULL"
				}
			}
			spew.Dump(values)
			writer.Row(values)
			writer.Row([]any{1, "Marry", "password", "email"})
		}

		return writer.Complete("SELECT 2")
	}

	columnHandler := func(ctx context.Context, query string) wire.Columns {
		rows, err := s.db.Query(query)
		defer rows.Close()
		if err != nil {
			panic(err)
		}

		cols, _ := rows.Columns()
		cols2 := make(wire.Columns, 0, len(cols))
		fmt.Println("got columns")
		for _, col := range cols {
			cols2 = append(cols2, wire.Column{
				Name:   col,
				Oid:    oid.T_text,
				Width:  1,
				Format: wire.TextFormat,
			})
			fmt.Println(col)
		}
		fmt.Println()
		return cols2
	}

	return QueryFuncs{statementHandler, columnHandler}
}

func (s *Server) defaultHandler(ctx context.Context, query string) QueryFuncs {
	f := func(ctx context.Context, writer wire.DataWriter, parameters []string) error {
		_, err := s.db.Exec(query)
		if err != nil {
			return err
		}
		return writer.Complete("OK")
	}

	columnsFunc := func(ctx context.Context, query string) wire.Columns {
		return wire.Columns{}
	}
	return QueryFuncs{f, columnsFunc}
}

func (s *Server) handler(ctx context.Context, query string) (wire.PreparedStatementFn, []oid.Oid, wire.Columns, error) {
	handlers := map[QueryType]QueryFuncs{
		CREATE: s.defaultHandler(ctx, query),
		INSERT: s.defaultHandler(ctx, query),
		SELECT: s.selectQueryHandler(ctx, query),
		UPDATE: s.defaultHandler(ctx, query),
		DELETE: s.defaultHandler(ctx, query),
		DROP:   s.defaultHandler(ctx, query),
	}
	queryFuncs := handlers[getQueryType(query)]

	return queryFuncs.queryHandler, wire.ParseParameters(query), queryFuncs.queryColumnsHandler(ctx, query), nil
}
