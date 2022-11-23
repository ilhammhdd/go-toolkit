package sqlkit

import (
	"database/sql"
	"fmt"
)

type TxCmdStmtArgs struct {
	Statement string
	Args      []interface{}
}

type DBOperation struct {
	DB *sql.DB
}

func (dbo DBOperation) TxCommand(txCmdStmtArgs []*TxCmdStmtArgs) ([]*sql.Result, error) {
	err := dbo.DB.Ping()
	if err != nil {
		return nil, err
	}

	tx, err := dbo.DB.Begin()
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, fmt.Errorf("error begin tx: %s, error rollback tx: %s", err.Error(), rollbackErr.Error())
		}
		return nil, err
	}
	var results []*sql.Result
	for i := range txCmdStmtArgs {
		preparedStmt, err := tx.Prepare(txCmdStmtArgs[i].Statement)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return nil, fmt.Errorf("error prepare tx statement: %s, error rollback tx: %s", err.Error(), rollbackErr.Error())
			}
			return nil, err
		}

		result, err := preparedStmt.Exec(txCmdStmtArgs[i].Args...)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return nil, fmt.Errorf("error exec tx prepared statement: %s, error rollback tx: %s", err.Error(), rollbackErr.Error())
			}
			return nil, err
		}
		defer preparedStmt.Close()
		results = append(results, &result)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (dbo DBOperation) Command(stmt string, args ...interface{}) (sql.Result, error) {
	err := dbo.DB.Ping()
	if err != nil {
		return nil, err
	}

	outStmt, err := dbo.DB.Prepare(stmt)
	if err != nil {
		return nil, err
	}
	defer outStmt.Close()

	result, err := outStmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (dbo DBOperation) Query(stmt string, args ...interface{}) (*sql.Rows, error) {
	err := dbo.DB.Ping()
	if err != nil {
		return nil, err
	}

	outStmt, err := dbo.DB.Prepare(stmt)
	if err != nil {
		return nil, err
	}
	defer outStmt.Close()

	rows, err := outStmt.Query(args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (dbo DBOperation) QueryRow(stmt string, args ...interface{}) (*sql.Row, error) {
	err := dbo.DB.Ping()
	if err != nil {
		return nil, err
	}

	outStmt, err := dbo.DB.Prepare(stmt)
	if err != nil {
		return nil, err
	}
	defer outStmt.Close()

	resultRow := outStmt.QueryRow(args...)

	return resultRow, nil
}

func (dbo DBOperation) QueryRowsToMap(stmt string, args ...interface{}) (*[]*map[string]interface{}, error) {
	err := dbo.DB.Ping()
	if err != nil {
		return nil, err
	}

	outStmt, err := dbo.DB.Prepare(stmt)
	if err != nil {
		return nil, err
	}
	defer outStmt.Close()

	rows, err := outStmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))

	var resultRows []*map[string]interface{}

	for rows.Next() {
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		err = rows.Scan(columnPointers...)
		if err != nil {
			return nil, err
		}

		columnTypes, _ := rows.ColumnTypes()
		resultRow := make(map[string]interface{})

		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			columnTypeName := (*columnTypes[i]).DatabaseTypeName()
			if columnTypeName == "TINYINT" {
				if (*val).(int64) == 1 {
					resultRow[colName] = true
				} else if (*val).(int64) == 0 {
					resultRow[colName] = false
				}
			} else {
				resultRow[colName] = *val
			}
		}

		resultRows = append(resultRows, &resultRow)
	}

	return &resultRows, nil
}
