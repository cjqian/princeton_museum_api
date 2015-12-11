package sqlParser

import (
	_ "./mysql"
	"./sqlx"
	"strconv"
	//	"encoding/json"
	"fmt"
)

var (
	globalDB sqlx.DB

	tableIDMap   map[string]string
	tableXrefMap map[string]string
	colTypeMap   map[string]string
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
	tableIDMap = GetTableIDMap()
	tableXrefMap = GetTableXRefMap()
	colTypeMap = GetColTypeMap()
	return *db
}

/*********************************************************************************
 * HELPER FUNCTIONS
 ********************************************************************************/
//returns array of table name strings from queried database
func GetTableNames() []string {
	var tableNames []string

	tableRawBytes := make([]byte, 1)
	tableInterface := make([]interface{}, 1)

	tableInterface[0] = &tableRawBytes

	rows, err := globalDB.Query("show tables")
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
	tableRawBytes := make([]byte, 1)
	tableInterface := make([]interface{}, 1)

	tableInterface[0] = &tableRawBytes

	rows, err := globalDB.Query("SELECT distinct column_name from information_schema.columns where table_name='" + tableName + "' ORDER BY column_name asc")
	check(err)

	colMap := make([]string, 0)

	for rows.Next() {
		err = rows.Scan(tableInterface...)
		check(err)

		colMap = append(colMap, string(tableRawBytes))
	}

	return colMap
}

func GetColumnValues(tableName string, columnName string) []string {
	fmt.Println(tableName + ", " + columnName)
	tableRawBytes := make([]byte, 1)
	tableInterface := make([]interface{}, 1)

	tableInterface[0] = &tableRawBytes

	rows, err := globalDB.Query("select distinct " + columnName + " from " + tableName)
	check(err)

	colMap := make([]string, 0)

	for rows.Next() {
		err = rows.Scan(tableInterface...)
		check(err)

		colMap = append(colMap, string(tableRawBytes))
	}

	return colMap
}

//returns one-to-one structure, currently specifically for apititleobjxrefs table
func GetOneToOne(objectID int, tTable string, channel chan interface{}) {
	queryStr := "select * from " + tTable + " where ObjectID=" + strconv.Itoa(objectID)
	rows, err := globalDB.Queryx(queryStr)
	if err != nil {
		panic(err)
	}

	rowArray := make([]map[string]interface{}, 0)

	for rows.Next() {
		results := make(map[string]interface{}, 0)
		err = rows.MapScan(results)
		if err != nil {
			panic(err)
		}

		for k, v := range results {
			if b, ok := v.([]byte); ok {
				results[k], err = StringToType(b, colTypeMap[k])
				if err != nil {
					panic(err)
				}
			}
		}

		rowArray = append(rowArray, results)
	}

	channel <- rowArray
}

//returns one-to-many structure given objectID, original table, target table, and
//connecting table
func GetOneToMany(objectID int, oTable string, tTable string, channel chan interface{}) {
	cTable := tableXrefMap[tTable]

	oIDname := oTable + "." + tableIDMap[oTable]
	tIDname := tTable + "." + tableIDMap[tTable]
	cIDname := cTable + "." + tableIDMap[cTable]

	cIDtrans := cTable + "." + tableIDMap[tTable]

	queryStr := "select " + tTable + ".* from " + cTable + " INNER JOIN " + oTable + " ON " + cIDname + " = " + oIDname + " INNER JOIN " + tTable + " ON " + cIDtrans + " = " + tIDname + " where " + oIDname + "=" + strconv.Itoa(objectID)

	fmt.Println(queryStr)
	rows, err := globalDB.Queryx(queryStr)
	if err != nil {
		panic(err)
	}

	rowArray := make([]map[string]interface{}, 0)

	//for each row
	for rows.Next() {
		//map the column name to its value
		results := make(map[string]interface{}, 0)
		err = rows.MapScan(results)
		if err != nil {
			panic(err)
		}

		for k, v := range results {
			if b, ok := v.([]byte); ok {
				results[k], err = StringToType(b, colTypeMap[k])
				if err != nil {
					panic(err)
				}
			}
		}

		rowArray = append(rowArray, results)
	}

	channel <- rowArray
}

func AppendSpecial(tableName string, objectIDval int, results map[string]interface{}) {
	if tableName == "apiobjects" {
		titlesChan := make(chan interface{})
		constituentsChan := make(chan interface{})
		mediaChan := make(chan interface{})
		dimensionsChan := make(chan interface{})

		go GetOneToOne(objectIDval, "apititleobjxrefs", titlesChan)
		go GetOneToMany(objectIDval, "apiobjects", "apiconstituents", constituentsChan)
		go GetOneToMany(objectIDval, "apiobjects", "apimedia", mediaChan)
		go GetOneToMany(objectIDval, "apiobjects", "apidimelements", dimensionsChan)

		results["Titles"] = <-titlesChan
		results["Constituents"] = <-constituentsChan
		results["Media"] = <-mediaChan
		results["Dimensions"] = <-dimensionsChan
	}
}

/*********************************************************************************
 * GET FUNCTIONALITY
 ********************************************************************************/
func GetNumResults(tableName string, tableParameters []string) int {
	var numResults int

	queryStr := "select count(*) from " + tableName

	//where
	if len(tableParameters) > 0 {
		queryStr += " where "

		for _, v := range tableParameters {
			queryStr += v + " and "
		}

		queryStr = queryStr[:len(queryStr)-4]
	}

	rows, err := globalDB.Queryx(queryStr)
	if err != nil {
		return -1
	}

	for rows.Next() {
		err := rows.Scan(&numResults)
		if err != nil {
			return -1
		}
	}

	return numResults
}

func Get(tableName string, tableParameters []string, specialParameters map[string]int) ([]map[string]interface{}, error) {
	//pagination
	size := specialParameters["size"]
	page := specialParameters["page"]

	regStr := ""
	joinStr := ""

	cols := GetColumnNames(tableName)
	for _, col := range cols {
		regStr += tableName + "." + col + ","
	}

	regStr = regStr[:len(regStr)-1]

	if joinStr != "" {
		joinStr = ", " + joinStr[:len(joinStr)-1]
	}

	queryStr := "select " + regStr + joinStr + " from " + tableName + " "

	/* MAIN QUERY */
	//where
	if len(tableParameters) > 0 {
		queryStr += " where "

		for _, v := range tableParameters {
			queryStr += v + " and "
		}

		queryStr = queryStr[:len(queryStr)-4]
	}

	//limitation
	startNum := (page - 1.0) * size
	queryStr += " LIMIT " + strconv.Itoa(startNum) + ", " + strconv.Itoa(size)

	//do the query
	rows, err := globalDB.Queryx(queryStr)
	if err != nil {
		return nil, err
	}

	//map into an array of type map[colName]value
	rowArray := make([]map[string]interface{}, 0)

	objectIDstring := "ObjectID"
	//for each row
	for rows.Next() {
		//map the column name to its value
		results := make(map[string]interface{}, 0)
		err = rows.MapScan(results)
		if err != nil {
			return nil, err
		}

		for k, v := range results {
			if b, ok := v.([]byte); ok {
				results[k], err = StringToType(b, colTypeMap[k])
				if err != nil {
					return nil, err
				}
			}

		}

		objectIDval := results[objectIDstring].(int)
		AppendSpecial(tableName, objectIDval, results)

		rowArray = append(rowArray, results)
	}

	return rowArray, nil
}
