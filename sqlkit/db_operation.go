package sqlkit

import (
	"database/sql"
)

type TxCmdStmtArgs struct {
	Statement string
	Args      []interface{}
}

type DBOperation struct {
	DB *sql.DB
}

func (dbo DBOperation) TxCommand(txCmdStmtArgs []*TxCmdStmtArgs) ([]*sql.Result, []error) {
	err := dbo.DB.Ping()
	if err != nil {
		return nil, []error{err}
	}

	tx, err := dbo.DB.Begin()
	if err != nil {
		return nil, []error{err}
	}
	var txPreparedStmtErrs []error
	var txExecErrs []error
	var results []*sql.Result
	for i := range txCmdStmtArgs {
		preparedStmt, err := tx.Prepare(txCmdStmtArgs[i].Statement)
		if err != nil {
			txPreparedStmtErrs = append(txPreparedStmtErrs, err)
		}
		if len(txPreparedStmtErrs) > 0 {
			break
		}

		result, err := preparedStmt.Exec(txCmdStmtArgs[i].Args...)
		if err != nil {
			txExecErrs = append(txExecErrs, err)
		}
		defer preparedStmt.Close()
		if len(txExecErrs) > 0 {
			break
		}
		results = append(results, &result)
	}

	if len(txPreparedStmtErrs) > 0 || len(txExecErrs) > 0 {
		return nil, append(txPreparedStmtErrs, txExecErrs...)
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
