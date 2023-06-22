package server

import (
	"database/sql"
	"fmt"
	wire "github.com/jeroenrinzema/psql-wire"
	"github.com/lib/pq/oid"
	"strings"
)

// https://zontroy.com/postgresql-to-go-type-mapping/

//PostgreSQL Data Type	Go Data Type	Go Nullable Data Type	Go Primitive Type
//bigint	int64		int64
//bit	bool		bool
//boolean	bool		bool
//bytea	[]byte
//string		string
//character varying	string		string
//date
//}
//double precision	float32		float32
//integer	int		int
//money	float64		float64
//numeric	float64		float64
//real	float32		float32
//serial	int		int
//smallint	int16		int16
//smallserial	int16		int16
//text	string		string

// TODO: flesh out data type mapping
var dataTypes = map[string]oid.Oid{
	"character":         oid.T_text,
	"string":            oid.T_text,
	"text":              oid.T_text,
	"character varying": oid.T_text,
	"varchar":           oid.T_text,

	"boolean": oid.T_bool,
	"bit":     oid.T_bool,
	"bytea":   oid.T_bytea,

	"int4": oid.T_int4,

	"timestamp": oid.T_timestamp,
}

func wireColumn(name string, colType sql.ColumnType) wire.Column {
	psqlDataType := colType.DatabaseTypeName()
	psqlDataType = strings.ToLower(psqlDataType)
	oidType := dataTypes[psqlDataType]
	if oidType == 0 {
		oidType = oid.T_text
		err := fmt.Errorf("unknown data type %s", psqlDataType)
		fmt.Println(err.Error())
	}

	return wire.Column{
		Name:   name,
		Oid:    oidType,
		Width:  1,
		Format: wire.TextFormat,
	}
}
