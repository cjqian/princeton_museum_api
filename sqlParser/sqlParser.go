package sqlParser

import (
	_ "./mysql"
	"./sqlx"
	"encoding/json"
	"fmt"
)

var (
	globalDB   sqlx.DB
	colTypeMap map[string]string
	tableMap   map[string][]string
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*********************************************************************************
 * DB INITIALIZE: Connects given DB creds, creates ColMap FOR SESSION
 ********************************************************************************/
func InitializeDatabase(username string, password string, environment string) sqlx.DB {
	db, err := sqlx.Connect("mysql", username+":"+password+"@tcp(localhost:3306)/"+environment)
	check(err)

	globalDB = *db

	//set global colTypeMap
	tableMap = GetTableMap()
	colTypeMap = GetColTypeMap()
	return *db
}

/*********************************************************************************
 * HELPER FUNCTIONS
 ********************************************************************************/
//if is table, returns 1. else (for example, is view), returns 0.
func IsTable(serverTableName string) bool {
	//check if there is view. else, assume is table
	query := "select exists(select * from information_schema.tables where table_name='" + serverTableName + "' and table_name not in (select table_name from information_schema.views))"
	rows, err := globalDB.Query(query)
	check(err)

	//set up scan interface
	rawBytes := make([]byte, 1)
	scanInterface := make([]interface{}, 1)
	scanInterface[0] = &rawBytes

	//this should only return one row, but Scan panics if not called with Next
	for rows.Next() {
		err := rows.Scan(scanInterface...)
		check(err)
		//if exists as view, delete from view
		if string(rawBytes) == "1" {
			return true
		} else {
			return false
		}
	}

	return false
}

//returns array of table name strings from queried database
func GetTableNames() []string {
	var tableNames []string

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

	return tableNames
}

//returns array of column names from table in database
func GetColumnNames(tableName string) []string {
	colNames := make([]string, 0)
	colNames = append(colNames, tableMap[tableName]...)

	return colNames
}

/*********************************************************************************
 * GET FUNCTIONALITY
 ********************************************************************************/
func Get(tableName string, tableParameters []string) ([]map[string]interface{}, error) {
	regStr := ""
	joinStr := ""
	onStr := ""

	cols := GetColumnNames(tableName)
	for _, col := range cols {
		regStr += tableName + "." + col + ","
	}

	regStr = regStr[:len(regStr)-1]

	if joinStr != "" {
		joinStr = ", " + joinStr[:len(joinStr)-1]
	}

	queryStr := "select " + regStr + joinStr + " from " + tableName + " "

	queryStr += onStr

	//where
	if len(tableParameters) > 0 {
		queryStr += " where "

		for _, v := range tableParameters {
			queryStr += v + " and "
		}

		queryStr = queryStr[:len(queryStr)-4]
	}

	fmt.Println(queryStr)
	//do the query
	rows, err := globalDB.Queryx(queryStr)
	if err != nil {
		return nil, err
	}

	//map into an array of type map[colName]value
	rowArray := make([]map[string]interface{}, 0)

	for rows.Next() {
		results := make(map[string]interface{}, 0)
		err = rows.MapScan(results)
		if err != nil {
			return nil, err
		}

		for k, v := range results {
			//converts the byte array to its correct type
			if b, ok := v.([]byte); ok {
				results[k], err = StringToType(b, colTypeMap[k])
				if err != nil {
					return nil, err
				}
			}
		}

		rowArray = append(rowArray, results)
	}

	return rowArray, nil
}
