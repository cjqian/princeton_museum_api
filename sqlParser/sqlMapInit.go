package sqlParser

import (
	/*"fmt"*/
	"strings"
)

//returns a map of each column name in table to its appropriate GoLang tpye (name string)
func GetColTypeMap() map[string]string {
	colMap := make(map[string]string, 0)

	cols, err := globalDB.Queryx("SELECT DISTINCT COLUMN_NAME, COLUMN_TYPE FROM information_schema.columns")
	check(err)

	for cols.Next() {
		var colName string
		var colType string

		err = cols.Scan(&colName, &colType)
		//split because SQL type returns are sometimes ex. int(11)
		colMap[colName] = strings.Split(colType, "(")[0]
	}

	return colMap
}

func GetTableMap() map[string][]string {
	var tableNames []string
	var tableMap = make(map[string][]string)

	tableRawBytes := make([]byte, 1)
	tableInterface := make([]interface{}, 1)

	tableInterface[0] = &tableRawBytes

	rows, err := globalDB.Query("SELECT TABLE_NAME FROM information_schema.tables where table_type='base table' or table_type='view'")
	check(err)

	for rows.Next() {
		err := rows.Scan(tableInterface...)
		check(err)

		tableNames = append(tableNames, string(tableRawBytes))
	}

	for _, table := range tableNames {
		rows, err = globalDB.Query("SELECT column_name from information_schema.columns where table_name='" + table + "' ORDER BY column_name asc")
		check(err)

		colMap := make([]string, 0)

		for rows.Next() {
			err = rows.Scan(tableInterface...)
			check(err)

			colMap = append(colMap, string(tableRawBytes))
		}

		tableMap[table] = colMap
	}
	return tableMap
}
