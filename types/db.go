package types

import (
	"fmt"
)

var ConstWidthTypes = []string{"float", "int"}
var VariableWidthTypes = []string{"char", "varchar"}

type DBState struct {
	CurrentDB string
}

func NewState() DBState {
	return DBState{
		CurrentDB: "",
	}
}

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

// TODO add type checking for when data is inserted/selected
type Type interface {
	ToString() string
}

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
