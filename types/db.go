// Noah Snelson
// February 25, 2021
// sdb/types/db.go
//
// Contains types for core database functionality, including current database
// state and all the types stored by the database.

package types

import (
	"fmt"
)

// Types based on fixed with vs. dynamic width. Used in `parser`.
var ConstWidthTypes = []string{"float", "int"}
var VariableWidthTypes = []string{"char", "varchar"}

// The DBState as of right now just tracks the current database that the user is
// using, as specified in the `USE <db>;` query.
type DBState struct {
	CurrentDB string
}

// Create a new DBState with no current database. Will not be valid - user must
// execute a `USE` query before actually executing any queries other than
// `CREATE DATABASE`.
func NewState() DBState {
	return DBState{
		CurrentDB: "",
	}
}

// Converts an arbitrary string to a `Type` interface, with the appropriate type
// parameters.
func NewType(typename string, size int) Type {
	if typename == "float" {
		return Float{}
	}
	if typename == "int" {
		return Int{}
	}
	if typename == "char" {
		return Char{ size }
	}
	if typename == "varchar" {
		return VarChar{ size }
	}
	return nil
}

// TODO Will also be used for type checking when DB implements insert/select
// functionality.
type Type interface {
	ToString() string
}

// ----- TYPE STRUCTS --------------------------------------------------------
// Structs for the various types of values that the database implements. These
// do not store *values*, but rather the parameters for each type, as in the
// variable width `char` and `varchar` types. Will be more useful if enums or
// other more complex types are implemented.

type Float struct {}

func (float Float) ToString() string {
	return "float"
}

type Int struct {}

func (int Int) ToString() string {
	return "int"
}

type Char struct {
	size int
}

func (char Char) ToString() string {
	return fmt.Sprintf("char(%v)", char.size)
}

type VarChar struct {
	size int
}

func (varchar VarChar) ToString() string {
	return fmt.Sprintf("varchar(%v)", varchar.size)
}
