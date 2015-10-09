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
 * DELETE FUNCTIONALITY
 ********************************************************************************/
func Delete(serverTableName string, parameters []string) (bool, error) {
	if !IsTable(serverTableName) {
		return DeleteFromView(serverTableName, parameters)
	} else {
		return false, DeleteFromTable(serverTableName, parameters)
	}
}

//deletes from a table
func DeleteFromTable(tableName string, parameters []string) error {
	return RunDeleteQuery(tableName, parameters)
}

//deletes from a view
func DeleteFromView(viewName string, parameters []string) (bool, error) {
	if len(parameters) == 0 {
		qStr := "drop view " + viewName
		_, err := globalDB.Query(qStr)
		return true, err
	} else {
		return false, RunDeleteQuery(viewName, parameters)
	}
}

//runs query of format "delete from tableName where parameterA=valueA and..."
func RunDeleteQuery(serverTableName string, parameters []string) error {
	//delete from tableName where x = a and y = b
	query := "delete from " + serverTableName

	if len(parameters) > 0 {
		query += " where "

		for _, v := range parameters {
			query += v + " and "
		}
		//removes last "and"
		query = query[:len(query)-4]
	}

	_, err := globalDB.Query(query)
	return err
}

/*********************************************************************************
 * GET FUNCTIONALITY
 ********************************************************************************/
func Get(tableName string) ([]map[string]interface{}, error) {
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

/*********************************************************************************
 * POST FUNCTIONALITY
 ********************************************************************************/
func Post(tableName string, jsonByte []byte) (string, error) {
	if IsTable(tableName) {
		err := PostRows(tableName, jsonByte)
		return tableName, err
	} else {
		return PostViews(jsonByte)
	}
}

//adds new row to table
func AddRow(newRow interface{}, tableName string) error {
	m := newRow.(map[string]interface{})
	//insert into table (colA, colB) values (valA, valB);
	query := "INSERT INTO " + tableName + " ("
	keyStr := ""
	valueStr := ""

	for k, v := range m {
		keyStr += k + ","
		typeStr, err := TypeToString(v)
		if err != nil {
			return err
		}

		valueStr += "'" + typeStr + "',"
	}

	keyStr = keyStr[:len(keyStr)-1]
	valueStr = valueStr[:len(valueStr)-1]

	query += keyStr + ") VALUES ( " + valueStr + " );"
	_, err := globalDB.Query(query)
	fmt.Println(query)
	return err
}

func AddRows(newRows []interface{}, tableName string) error {
	for _, row := range newRows {
		err := AddRow(row, tableName)
		if err != nil {
			return err
		}
	}

	return nil
}

//adds JSON from FILENAME to TABLE
//CURRENTLY ONLY ONE ROW
func PostRows(tableName string, jsonByte []byte) error {
	var f []interface{}

	err := json.Unmarshal(jsonByte, &f)
	if err != nil {
		return err
	}

	err2 := AddRows(f, tableName)
	if err2 != nil {
		return err
	}

	return nil
}

//view details are marshalled into this View struct
type View struct {
	Name  string
	Query string
}

//adds JSON from FILENAME to TABLE
func PostViews(jsonByte []byte) (string, error) {
	var views []View

	err := json.Unmarshal(jsonByte, &views)
	if err != nil {
		return "", err
	}

	var viewName string
	for _, view := range views {
		viewName = view.Name
		err = MakeView(view.Name, view.Query)
		if err != nil {
			return viewName, err
		}
	}

	return viewName, nil
}

func MakeView(viewName string, view string) error {
	qStr := "create view " + viewName + " as " + view
	_, err := globalDB.Query(qStr)
	tableMap = GetTableMap()
	return err
}

/*********************************************************************************
 * PUT FUNCTIONALITY
 ********************************************************************************/
func Put(tableName string, parameters []string, jsonByte []byte) error {
	//unmarshals the json into an interface
	var f []interface{}
	err := json.Unmarshal(jsonByte, &f)
	if err != nil {
		return err
	}
	//adds the interface row to table in database
	return UpdateRows(f, tableName, parameters)
}

func UpdateRows(newRows []interface{}, tableName string, parameters []string) error {
	for _, row := range newRows {
		err := UpdateRow(row, tableName, parameters)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateRow(newRow interface{}, tableName string, parameters []string) error {
	query := "update " + tableName

	updateParameters := newRow.(map[string]interface{})
	//new changes
	if len(updateParameters) > 0 {
		query += " set "

		for k, v := range updateParameters {
			typeStr, err := TypeToString(v)
			if err != nil {
				return err
			}
			query += k + "='" + typeStr + "', "
		}

		query = query[:len(query)-2]
	}

	//where
	if len(parameters) > 0 {
		query += " where "

		for _, v := range parameters {
			query += v + " and "
		}

		query = query[:len(query)-4]
	}

	_, err := globalDB.Query(query)
	return err
}
