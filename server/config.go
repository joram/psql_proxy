package server

import (
	"database/sql"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var functions = map[string]AnonymizingFunction{
	"email":             email,
	"password.stars":    passwordStars,
	"name.first.dwarf":  nameFirstDwarf,
	"name.middle.dwarf": nameMiddleDwarf,
	"name.last.dwarf":   nameLastDwarf,
}

type Column struct {
	Name         string `yaml:"name"`
	FunctionName string `yaml:"function"`

	function *AnonymizingFunction
}

type Conf struct {
	Columns   []*Column `yaml:"columns"`
	Port      int       `yaml:"port"`
	ServerURI string    `yaml:"server_uri"`

	columnFuncs map[string]*AnonymizingFunction
}

func (c *Conf) getAnonymizingFunc(name string) (AnonymizingFunction, error) {
	f := functions[name]
	if f != nil {
		return f, nil
	}
	return nil, errors.New(fmt.Sprintf("Config Error: '%s' is not a valid anonymizing function", name))

}

func (c *Conf) Anonymize(columnName string, columnType *sql.ColumnType, value interface{}) interface{} {
	f := c.columnFuncs[columnName]
	if f == nil {
		return value
	}
	return (*f)(value)
}

func GetConfig() *Conf {

	// get file content
	yamlFile, err := ioutil.ReadFile("/config.yaml")
	if err != nil {
		panic(errors.New(fmt.Sprintf("yamlFile.Get err   #%v ", err)))
	}

	// unmarshal to struct
	c := &Conf{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		panic(errors.New(fmt.Sprintf("unmarshal:   #%v ", err)))
	}

	// find anonymizing funcs
	c.columnFuncs = make(map[string]*AnonymizingFunction)
	for _, col := range c.Columns {
		f, err := c.getAnonymizingFunc(col.FunctionName)
		if err != nil {
			panic(fmt.Sprintf("Config Error: %v", err))
		}
		col.function = &f
		c.columnFuncs[col.Name] = &f
	}

	return c
}
