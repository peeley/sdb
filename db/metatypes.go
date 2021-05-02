// Noah Snelson
// February 25, 2021
// sdb/db/metatypes.go
//
// Contains types for core database functionality, including current database
// state and all the types stored by the database.

package db

import (
	"fmt"
	"strings"
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

type Column struct {
	Name string
	Type Type
}

type Value struct {
	Value interface{}
	Type Type
}

func (v *Value) GetValue() interface{} {
	return v.Value
}

func (v *Value) GetType() Type {
	return v.Type
}

func (v *Value) TypeMatches(t *Type) bool {
	if strings.Contains(v.GetType().ToString(), "int") {
		return v.GetType().ToString() == (*t).ToString()
	} else if strings.Contains(v.GetType().ToString(), "float") {
		return v.GetType().ToString() == (*t).ToString()
	}
	// v is a varchar or char, in which case it's valid as long as it's <= the
	// column's required length
	candVarChar, ok := v.GetType().(VarChar)
	if ok { // inserted type is varchar
		colVarChar, ok := (*t).(Char)
		if ok { // cand is varchar, column is char
			return candVarChar.Size <= colVarChar.Size
		} else { // cand is varchar, col is varchar
			colChar := (*t).(VarChar)
			return candVarChar.Size <= colChar.Size
		}
	}
	candChar := v.GetType().(Char)
	colChar, ok := (*t).(Char)
	if ok { // cand is char, column is char
		return candChar.Size == colChar.Size
	} else { // cand is char, col is varchar
		colVarChar := (*t).(VarChar)
		return candChar.Size <= colVarChar.Size
	}
}

func (v *Value) ToString() string {
	if v.Type.ToString() == "float" {
		return fmt.Sprintf("%v", v.Value)
	} else if v.Type.ToString() == "int" {
		return fmt.Sprintf("%v", v.Value)
	}
	// otherwise, value is a string of some kind
	return fmt.Sprintf("'%v'", v.Value)
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
	Size int
}

func (char Char) ToString() string {
	return fmt.Sprintf("char(%v)", char.Size)
}

type VarChar struct {
	Size int
}

func (varchar VarChar) ToString() string {
	return fmt.Sprintf("varchar(%v)", varchar.Size)
}

type Null struct {}

func (null Null) ToString() string {
	return "NULL"
}
