package sqlParser

import (
	_ "./mysql"
	"./sqlx"
	"fmt"
	"strconv"
	//"encoding/json"
)

var tableNameToId = map[string]string{
	"apiobjects":      "ObjectID",
	"apiconstituents": "ConstituentID",
}

var (
	globalDB sqlx.DB

	colTypeMap map[string]string
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
	//fmt.Println(tableName + ", " + columnName)
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

/*********************************************************************************
 * APIOBJECTS HELPER FUNCTIONS
 ********************************************************************************/
func QueryRows(queryStr string) []map[string]interface{} {
	fmt.Println(queryStr)

	rows, err := globalDB.Queryx(queryStr)
	if err != nil {
		panic(err)
	}

	rowArray := make([]map[string]interface{}, 0)

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

	return rowArray
}

//returns data from the apiconobjxrefs data
func GetBibliography() []map[string]interface{} {
	queryStr := "select ReferenceID, ObjCitation, apibibobjxrefs.SysTimeStamp, apiobjects.ObjectID from apibibobjxrefs INNER JOIN queryView as apiobjects ON apibibobjxrefs.ObjectID = apiobjects.ObjectID"

	return QueryRows(queryStr)
}

func GetConstituentsTrunc() []map[string]interface{} {
	params := "apiconstituents.Active, apiconstituents.AlphaSort, apiconstituents.Approved, apiconstituents.BeginDate, apiconstituents.BeginDateISO, apiconstituents.ConstituentID, apiconobjxrefs.DisplayOrder, apiconobjxrefs.Displayed, apiconobjxrefs.Prefix, apiconobjxrefs.Remarks, apiconobjxrefs.Role, apiconobjxrefs.Suffix, apiconstituents.SysTimeStamp, apiobjects.ObjectID"

	queryStr := "select " + params + " from apiconobjxrefs " + "INNER JOIN queryView as apiobjects ON apiconobjxrefs.ObjectID = apiobjects.ObjectID " + "INNER JOIN apiconstituents ON apiconobjxrefs.ConstituentID = apiconstituents.ConstituentID ORDER BY apiconobjxrefs.DisplayOrder"

	return QueryRows(queryStr)
}

func GetDimElements(whereStr string, limStr string, rowCount int) []map[string]interface{} {
	queryStr := `select apidimelements.*, apidimobjxrefs.* from apidimobjxrefs 
		INNER JOIN apiobjects ON apidimobjxrefs.ObjectID = apiobjects.ObjectID 
		INNER JOIN apidimelements ON apidimobjxrefs.DimItemElemXrefID = apidimelements.DimItemElemXrefID `
	queryStr += whereStr + limStr

	return QueryRows(queryStr)
}

//returns data from the apiconobjxrefs data
func GetExhibitions(whereStr string, limStr string, rowCount int) []map[string]interface{} {
	queryStr := "select ExhibitionID, RunningCaption, apiexhobjxrefs.SysTimeStamp from apiexhobjxrefs INNER JOIN apiobjects " + whereStr + limStr

	return QueryRows(queryStr)
}

//returns data from the apiconobjxrefs data
func GetGeography(whereStr string, limStr string, rowCount int) []map[string]interface{} {
	queryStr := `select ObjGeographyID, apiobjgeography.ObjectID, GeoCode, PrimaryDisplay, Continent, SubContinent, Country, 
		Region, State, City, Country, SubRegion, Locale, Locus, River, Excavation, 
		Latitude, Longitude, GeoNames, apiobjgeography.SysTimeStamp from apiobjgeography INNER JOIN apiobjects `
	queryStr += whereStr + limStr

	return QueryRows(queryStr)
}

func GetMedia(whereStr string, limStr string, rowCount int) []map[string]interface{} {
	queryStr := `select apimedia.* from apimediaxrefs 
		INNER JOIN apiobjects ON apimediaxrefs.ID = apiobjects.ObjectID 
		INNER JOIN apimedia ON apimediaxrefs.MediaMasterID = apimedia.MediaMasterID `

	queryStr += whereStr + limStr

	return QueryRows(queryStr)
}

//returns data from the apiconobjxrefs data
func GetTerms(whereStr string, limStr string, rowCount int) []map[string]interface{} {
	queryStr := "select * from apitermsobjxrefs INNER JOIN apiobjects " + whereStr + limStr

	return QueryRows(queryStr)
}

func GetTitles(whereStr string, limStr string, rowCount int) []map[string]interface{} {
	queryStr := "select * from apititleobjxrefs INNER JOIN apiobjects " + whereStr + limStr

	return QueryRows(queryStr)
}

//returns data from the apiconobjxrefs data
func GetUri(constituentID int, channel chan interface{}) {
	queryStr := "select * from apiconuris where ConstituentID = " + strconv.Itoa(constituentID)
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

func AddSubObjects(curResult map[string]interface{}, tableName string, subResults []map[string]interface{}, subIdx int) int {
	subArray := make([]interface{}, 0)

	for subIdx < len(subResults) && curResult["ObjectID"] == subResults[subIdx]["ObjectID"] {
		fmt.Println("Cur object id: " + strconv.Itoa(curResult["ObjectID"].(int)))
		fmt.Println("Bib object id: " + strconv.Itoa(subResults[subIdx]["ObjectID"].(int)))

		delete(subResults[subIdx], "ObjectID")
		subArray = append(subArray, subResults[subIdx])
		subIdx++
	}

	curResult[tableName] = subArray

	return subIdx
}

/*********************************************************************************
 * GET TABLES
 ********************************************************************************/
//func QueryConstituents(tableName string, idVal int, results map[string]interface{}) {
//uriChan := make(chan interface{})
//go GetUri(idVal, uriChan)
//results["URIs"] = <-uriChan
//}

func QueryObjects(whereStr string, limStr string, rowCount int, results []map[string]interface{}) {
	bibliographyResults := GetBibliography()
	constituentResults := GetConstituentsTrunc()
	//dimensionResults := GetDimElements(whereStr, limStr, rowCount)
	//exhibitionResults := GetExhibitions(whereStr, limStr, rowCount)
	//geographyResults := GetGeography(whereStr, limStr, rowCount)
	//mediaResults := GetMedia(whereStr, limStr, rowCount)
	//termResults := GetTerms(whereStr, limStr, rowCount)
	//titleResults := GetTitles(whereStr, limStr, rowCount)

	bibIdx := 0
	constituentIdx := 0

	for i := 0; i < len(results); i++ {
		fmt.Println(i)
		bibIdx = AddSubObjects(results[i], "Bibliography", bibliographyResults, bibIdx)
		constituentIdx = AddSubObjects(results[i], "Constituents", constituentResults, constituentIdx)
		//results[i]["Constituents"] = constituentResults[i]
		//results[i]["Dimensions"] = dimensionResults[i]
		//results[i]["Exhibitions"] = exhibitionResults[i]
		//results[i]["Geography"] = geographyResults[i]
		//results[i]["Media"] = mediaResults[i]
		//results[i]["Terms"] = termResults[i]
		//results[i]["Titles"] = titleResults[i]
	}
}

/*********************************************************************************
 * GET FUNCTIONALITY
 ********************************************************************************/
func GetNumRows(tableName string, whereStr string, limStr string) int {
	var rowCount int
	queryStr := "SELECT count(*) FROM " + tableName + " " + whereStr + limStr

	err := globalDB.Get(&rowCount, queryStr)
	if err != nil {
		panic(err)
	}

	return rowCount
}

func GetWhereString(tableParameters []string) string {
	whereStr := ""
	//where
	if len(tableParameters) > 0 {
		whereStr += " where "

		for _, v := range tableParameters {
			whereStr += v + " and "
		}

		whereStr = whereStr[:len(whereStr)-4]
	}

	return whereStr
}

func GetLimString(specialParameters map[string]int) string {
	//pagination
	size := specialParameters["size"]
	page := specialParameters["page"]

	//limitation
	startNum := (page - 1.0) * size
	limStr := " LIMIT " + strconv.Itoa(startNum) + ", " + strconv.Itoa(size)

	return limStr
}

//makes a view with name "queryView", no return value
func MakeView(tableName string, whereStr string, limitStr string) {
	selectStatement := "select * from " + tableName + " " + whereStr + " " + limitStr
	query := "create view queryView as " + selectStatement
	fmt.Println(query)

	_, err := globalDB.Query("DROP VIEW IF EXISTS queryView")
	if err != nil {
		panic(err)
	}

	_, err = globalDB.Query(query)
	if err != nil {
		panic(err)
	}

	//now, there should be a view with the name "queryView"
}

func Get(tableName string, tableParameters []string, specialParameters map[string]int) ([]map[string]interface{}, error) {
	//pagination

	regStr := ""
	//joinStr := ""

	cols := GetColumnNames(tableName)
	for _, col := range cols {
		regStr += col + ","
	}

	if len(cols) > 0 {
		regStr = regStr[:len(regStr)-1]
	}

	//if joinStr != "" {
	//joinStr = ", " + joinStr[:len(joinStr)-1]
	//}

	whereStr := GetWhereString(tableParameters)
	limStr := GetLimString(specialParameters)
	//make the view
	MakeView(tableName, whereStr, limStr)

	//queryStr := "select " + regStr + joinStr + " from queryView"
	queryStr := "select " + regStr + " from queryView"

	//get number of rows
	rowCount := GetNumRows(tableName, whereStr, limStr)

	//map into an array of type map[colName]value
	rowArray := QueryRows(queryStr)

	//query the special tables
	if tableName == "apiobjects" {
		QueryObjects(whereStr, limStr, rowCount, rowArray)
	} else if tableName == "apiconstituents" {
		//QueryConstituents(whereStr, size)
	}

	//then, remove the view

	return rowArray, nil
}
