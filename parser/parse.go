// Noah Snelson
// February 25, 2021
// sdb/parser/parse.go
//
// Contains main parsing function. The parser operates as a standard recursive
// descent parser -> https://en.wikipedia.org/wiki/Recursive_descent_parser
// Each statement type (`CREATE`, `USE`, `SELECT`) represent a top-level
// production, with each having its own parsing function in the respective
// `parser` package file.

package parser

import (
	"errors"
	"sdb/db"
	"sdb/statements"
	"sdb/utils"
	"strings"
)

// Main parsing function, prompt input/stdin is fed in as a parameter and a
// db.Executable interface is returned to be executed in sdb/main.go.
func Parse(input string) (db.Executable, error) {
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)

	if utils.IsComment(input) {
		return statements.Comment{}, nil
	}

	dropDB, err := ParseDropDBStatement(input)

	if err != nil {
		return nil, err
	} else if dropDB != nil {
		return dropDB, nil
	}

	dropTable, err := ParseDropTableStatement(input)

	if err != nil {
		return nil, err
	} else if dropTable != nil {
		return dropTable, nil
	}

	useDB, err := ParseUseDBStatement(input)

	if err != nil {
		return nil, err
	} else if useDB != nil {
		return useDB, nil
	}

	selectStatement, err := ParseSelectStatement(input)

	if err != nil {
		return nil, err
	} else if selectStatement != nil {
		return selectStatement, nil
	}

	alterStatement, err := ParseAlterStatement(input)

	if err != nil {
		return nil, err
	} else if alterStatement != nil {
		return alterStatement, nil
	}

	createTable, err := ParseCreateTableStatement(input)

	if err != nil {
		return nil, err
	} else if createTable != nil {
		return createTable, nil
	}

	createDB, err := ParseCreateDBStatement(input)

	if err != nil {
		return nil, err
	} else if createDB != nil {
		return createDB, nil
	}

	insert, err := ParseInsertStatement(input)

	if err != nil {
		return nil, err
	} else if insert != nil {
		return insert, nil
	}

	update, err := ParseUpdateStatement(input)

	if err != nil {
		return nil, err
	} else if update != nil {
		return update, nil
	}

	delete, err := ParseDeleteStatement(input)

	if err != nil {
		return nil, err
	} else if delete != nil {
		return delete, nil
	}

	transaction, err := ParseBeginTransaction(input)

	if err != nil {
		return nil, err
	} else if transaction != nil {
		return transaction, nil
	}

	commit, err := ParseCommit(input)

	if err != nil {
		return nil, err
	} else if commit != nil {
		return commit, nil
	}

	return nil, errors.New("!Syntax error.")
}
