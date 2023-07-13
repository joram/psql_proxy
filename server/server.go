package server

import (
	"context"
	"database/sql"
	"fmt"
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
	config  *Conf
}

type QueryFuncs struct {
	queryHandler        wire.PreparedStatementFn
	queryColumnsHandler QueryColumnsHandler
}
type QueryColumnsHandler func(ctx context.Context, query string) (wire.Columns, error)
type QueryHandler func(ctx context.Context, query string) QueryFuncs

func NewServer() (*Server, error) {
	config := GetConfig()

	db, err := sql.Open("postgres", config.ServerURI)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("proxying to %s\n", config.ServerURI)

	return &Server{
		config:  config,
		port:    config.Port,
		address: fmt.Sprintf("0.0.0.0:%d", config.Port),
		db:      db,
		dbURL:   config.ServerURI,
	}, nil
}

func (s *Server) Run() error {
	fmt.Printf("Listening to %s\n", s.address)
	return wire.ListenAndServe(s.address, s.handler)
}

func (s *Server) selectQueryHandler(_ context.Context, query string) QueryFuncs {

	statementHandler := func(ctx context.Context, writer wire.DataWriter, parameters []string) error {
		rows, err := s.db.Query(query)
		defer func() {
			err = rows.Close()
		}()
		if err != nil {
			return err
		}

		columns, err := rows.Columns()
		if err != nil {
			return err
		}

		columnTypes, err := rows.ColumnTypes()
		if err != nil {
			return err
		}

		count := len(columns)
		values := make([]any, count)
		valuePointers := make([]interface{}, count)
		for i := range values {
			valuePointers[i] = &values[i]
		}

		numRows := 0
		for rows.Next() {
			_ = rows.Scan(valuePointers...)
			for i, val := range values {
				values[i] = s.config.Anonymize(
					columns[i],
					columnTypes[i],
					val,
				)
			}
			err = writer.Row(values)
			if err != nil {
				return err
			}
			numRows += 1
		}

		return writer.Complete(fmt.Sprintf("SELECT %d", numRows))
	}

	columnHandler := func(ctx context.Context, query string) (wire.Columns, error) {
		rows, err := s.db.Query(query)
		defer func() {
			err = rows.Close()
		}()
		if err != nil {
			return nil, err
		}

		cols, _ := rows.Columns()
		colTypes, _ := rows.ColumnTypes()
		cols2 := make(wire.Columns, 0, len(cols))
		for i, name := range cols {
			cols2 = append(cols2, wireColumn(name, *colTypes[i]))
		}
		return cols2, nil
	}

	return QueryFuncs{statementHandler, columnHandler}
}

func (s *Server) defaultHandler(_ context.Context, query string) QueryFuncs {
	f := func(ctx context.Context, writer wire.DataWriter, parameters []string) error {
		_, err := s.db.Exec(query)
		if err != nil {
			fmt.Printf(err.Error())
			return err
		}
		return writer.Complete("OK")
	}

	columnsFunc := func(ctx context.Context, query string) (wire.Columns, error) {
		return wire.Columns{}, nil
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
	columns, err := queryFuncs.queryColumnsHandler(ctx, query)
	if err != nil {
		fmt.Printf("error! %s\n", err.Error())
		return nil, nil, nil, err
	}

	return queryFuncs.queryHandler, wire.ParseParameters(query), columns, nil
}
